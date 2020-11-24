package hub

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ChannelHub is a hub based on channels
type ChannelHub struct {
	state       Storage
	broadcast   chan *payload
	register    chan Processor
	unregister  chan Processor
	subscribe   chan *subscriptionMessage
	unsubscribe chan *unsubscriptionMessage
	topics      chan *topicOpMessage
	done        chan struct{}
}

// NewChannelHub creates a new channel hub and initiates a list of topics
func NewChannelHub(ctx context.Context, storage Storage, topics ...Topic) (Hub, error) {
	h := &ChannelHub{
		state:       storage,
		broadcast:   make(chan *payload),
		register:    make(chan Processor),
		unregister:  make(chan Processor),
		subscribe:   make(chan *subscriptionMessage),
		unsubscribe: make(chan *unsubscriptionMessage),
		topics:      make(chan *topicOpMessage),
		done:        make(chan struct{}),
	}

	for _, topic := range topics {
		if err := h.state.CreateTopic(ctx, topic); err != nil {
			return nil, err
		}
	}

	h.state.CreateTopic(ctx, TopicFromString("lists"))

	return h, nil
}

// GetProcessor returns the processor based on the id, returns an error if the processor is not registered
func (h *ChannelHub) GetProcessor(ctx context.Context, processorID string) (Processor, error) {
	return h.state.GetProcessor(ctx, processorID)
}

// RegisterProcessor registers the processor in the hub
func (h *ChannelHub) RegisterProcessor(ctx context.Context, p Processor) error {
	processor, _ := h.state.GetProcessor(ctx, p.GetID())
	if processor != nil {
		return fmt.Errorf("Processor %v is already registered", p.GetID())
	}

	h.register <- p

	return nil
}

// UnregisterProcessor unsubscribes the processor from all the topics and unregisters it
func (h *ChannelHub) UnregisterProcessor(ctx context.Context, p Processor) error {
	if _, err := h.state.GetProcessor(ctx, p.GetID()); err != nil {
		return err
	}

	h.unregister <- p

	return nil
}

// AddTopic creates a new topic
func (h *ChannelHub) AddTopic(ctx context.Context, topic Topic) error {
	// check if topic does not already exist
	exists, err := h.state.TopicExists(ctx, topic)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("Topic %v already exists", topic)
	}

	h.topics <- &topicOpMessage{
		topic: topic,
		op:    topicCreation,
	}

	return nil
}

// DeleteTopic removes a topic from the hub, send a delete message to its subscribers and unsubscribes them
func (h *ChannelHub) DeleteTopic(ctx context.Context, topic Topic) error {
	// check if topic exists
	exists, err := h.state.TopicExists(ctx, topic)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("%v does not exist", topic)
	}

	h.topics <- &topicOpMessage{
		topic: topic,
		op:    topicDeletion,
	}

	return nil
}

// Subscribe is a high level function that makes a processor subscribe to a topic
func (h *ChannelHub) Subscribe(ctx context.Context, p Processor, topic Topic) error {
	// check subscription
	subscribed, err := h.state.HasSubscribed(ctx, p, topic)
	if err != nil {
		return err
	}
	if subscribed {
		return fmt.Errorf("Processor %v has already subscribed to topic %v", p.GetID(), topic)
	}

	h.subscribe <- &subscriptionMessage{
		topic:     topic,
		processor: p,
	}

	return nil
}

// Unsubscribe is a high level function that unregisters a processor from a topic
func (h *ChannelHub) Unsubscribe(ctx context.Context, p Processor, topic Topic) error {
	// check subscription
	subscribed, err := h.state.HasSubscribed(ctx, p, topic)
	if err != nil {
		return err
	}

	if !subscribed {
		return fmt.Errorf("Processor %v has not subscribed to topic %v", p.GetID(), topic)
	}

	h.unsubscribe <- &unsubscriptionMessage{
		topic:     topic,
		processor: p,
	}

	return nil
}

// Publish sends a message to a topic
func (h *ChannelHub) Publish(ctx context.Context, msg Message) error {
	exists, err := h.state.TopicExists(ctx, msg.GetTopic())
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("The %v topic does not exist", msg.GetTopic())
	}
	log.Printf("Message published in topic %v", msg.GetTopic())
	h.broadcast <- &payload{
		ID:    msg.GetID(),
		Topic: msg.GetTopic(),
		Type:  msg.GetType(),
		Msg:   msg,
	}
	return nil
}

// Close safely stops the hub by unregistering all processors
func (h *ChannelHub) Close(ctx context.Context) error {
	processors, err := h.state.ListProcessors(ctx)
	if err != nil {
		return fmt.Errorf("Cannot close the hub : %w", err)
	}
	for _, processor := range processors {
		h.UnregisterProcessor(ctx, processor)
	}
	h.done <- struct{}{}

	return nil
}

// Run starts the hub. It listens to its register, unregister and message channels to act. Finally, if a message is received in the done channel, it stops
func (h *ChannelHub) Run(ctx context.Context, interrupt chan struct{}) error {
	defer func() {
		interrupt <- struct{}{}
	}()
	for {
		select {
		case topicOp := <-h.topics:
			switch topicOp.op {
			case topicCreation:
				if err := h.state.CreateTopic(ctx, topicOp.topic); err != nil {
					return err
				}
			case topicDeletion:
				processors, err := h.state.GetSubscribers(ctx, topicOp.topic)
				if err != nil {
					return err
				}

				msg := &DeleteTopicMessage{
					BaseMessage: BaseMessage{
						ID:    time.Now().Unix(),
						Topic: topicOp.topic,
					},
					Topic: topicOp.topic,
				}

				for _, processor := range processors {
					processor.GetMsgChannel() <- msg
					if err := h.state.Unsubscribe(ctx, processor, topicOp.topic); err != nil {
						return err
					}
				}

				if err := h.state.DeleteTopic(ctx, topicOp.topic); err != nil {
					return err
				}
			}

		case registration := <-h.register:
			if err := h.state.RegisterProcessor(ctx, registration); err != nil {
				return err
			}
			log.Printf("Registered processor %v", registration.GetID())

		case unregistration := <-h.unregister:
			topics, err := h.state.ListTopics(ctx)
			if err != nil {
				return err
			}
			for _, topic := range topics {
				subscribed, err := h.state.HasSubscribed(ctx, unregistration, topic)
				if err != nil {
					return err
				}
				if subscribed {
					if err := h.state.Unsubscribe(ctx, unregistration, topic); err != nil {
						return err
					}
				}
			}
			if err := h.state.UnregisterProcessor(ctx, unregistration); err != nil {
				return err
			}
			close(unregistration.GetMsgChannel())
			log.Printf("Unregistered processor %v", unregistration.GetID())

		case subscription := <-h.subscribe:
			if err := h.state.Subscribe(ctx, subscription.processor, subscription.topic); err != nil {
				return err
			}
			log.Printf("Processor %v successfully subscribed to topic %v", subscription.processor.GetID(), subscription.topic)

		case unsubscription := <-h.unsubscribe:
			if err := h.state.Unsubscribe(ctx, unsubscription.processor, unsubscription.topic); err != nil {
				return err
			}
			log.Printf("Processor %v unsubscribed from topic %v", unsubscription.processor.GetID(), unsubscription.topic)

		case message := <-h.broadcast:
			processors, err := h.state.GetSubscribers(ctx, message.Topic)
			if err != nil {
				return err
			}

			for _, processor := range processors {
				select {
				case processor.GetMsgChannel() <- message.Msg:
				default:
					h.unsubscribe <- &unsubscriptionMessage{
						topic:     message.Topic,
						processor: processor,
					}
				}
			}

		case <-h.done:
			return nil
		}
	}
}
