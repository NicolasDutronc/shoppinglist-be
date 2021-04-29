package hub

import (
	"context"
	"fmt"
)

// Storage keeps track of processors and subscriptions
type Storage struct {
	Topics            map[Topic]map[Processor]bool
	Processors        map[string]Processor
	MessageHooks      map[Processor]bool
	SubscriptionHooks map[SubscriptionHook]bool
}

// NewStorage returns an empty Storage
func NewStorage() *Storage {
	return &Storage{
		Topics:            make(map[Topic]map[Processor]bool),
		Processors:        make(map[string]Processor),
		MessageHooks:      make(map[Processor]bool),
		SubscriptionHooks: make(map[SubscriptionHook]bool),
	}
}

// CreateTopic checks if the topic already exists and creates it. An error is returned instead if the topic already exists.
func (hs *Storage) CreateTopic(ctx context.Context, topic Topic) error {
	if _, exists := hs.Topics[topic]; exists {
		return fmt.Errorf("Topic %v already exists", topic)
	}

	hs.Topics[topic] = make(map[Processor]bool)

	return nil
}

// DeleteTopic checks if the topic exists and deletes it. An error is returned instead if the topic does not exist.
func (hs *Storage) DeleteTopic(ctx context.Context, topic Topic) error {
	if _, exists := hs.Topics[topic]; !exists {
		return fmt.Errorf("Topic %v does not exist", topic)
	}

	delete(hs.Topics, topic)

	return nil
}

// ListTopics returns all stored topics as a list
func (hs *Storage) ListTopics(ctx context.Context) ([]Topic, error) {
	topics := []Topic{}
	for topic := range hs.Topics {
		topics = append(topics, topic)
	}

	return topics, nil
}

// TopicExists returns true if the topic exists in the hub
func (hs *Storage) TopicExists(ctx context.Context, topic Topic) (bool, error) {
	_, exists := hs.Topics[topic]
	return exists, nil
}

// RegisterProcessor adds a processor in the processors map. An error is returned instead if the processor is already registered
func (hs *Storage) RegisterProcessor(ctx context.Context, p Processor) error {
	if _, exists := hs.Processors[p.GetID()]; exists {
		return fmt.Errorf("Processor %v is already registered", p.GetID())
	}

	hs.Processors[p.GetID()] = p

	return nil
}

// UnregisterProcessor delete a processor from the processors map. An error is returned instead if the given processor is not registered
func (hs *Storage) UnregisterProcessor(ctx context.Context, p Processor) error {
	if _, exists := hs.Processors[p.GetID()]; !exists {
		return fmt.Errorf("Processor %v is not registered", p.GetID())
	}

	delete(hs.Processors, p.GetID())

	return nil
}

// GetProcessor returns the processor that has the given id. An error is returned instead if no processor with the given id is found
func (hs *Storage) GetProcessor(ctx context.Context, ID string) (Processor, error) {
	if _, exists := hs.Processors[ID]; !exists {
		return nil, fmt.Errorf("Processor %v is not registered", ID)
	}

	return hs.Processors[ID], nil
}

// ListProcessors returns all registered processors as a list
func (hs *Storage) ListProcessors(ctx context.Context) ([]Processor, error) {
	processors := []Processor{}
	for _, p := range hs.Processors {
		processors = append(processors, p)
	}

	return processors, nil
}

// Subscribe adds the processor p to the subscribers of topic t. An error is returned instead if the processor or the topic does not exist
func (hs *Storage) Subscribe(ctx context.Context, p Processor, t Topic) error {
	// check processor
	if _, err := hs.GetProcessor(ctx, p.GetID()); err != nil {
		return err
	}

	// check topic
	if _, exists := hs.Topics[t]; !exists {
		return fmt.Errorf("Topic %v does not exist", t)
	}

	hs.Topics[t][p] = true

	return nil
}

// Unsubscribe deletes processor p from the subscribers of t. An error is returned instead if the processor or the topic does not exist
func (hs *Storage) Unsubscribe(ctx context.Context, p Processor, t Topic) error {
	// check processor
	if _, err := hs.GetProcessor(ctx, p.GetID()); err != nil {
		return err
	}

	// check topic
	if _, exists := hs.Topics[t]; !exists {
		return fmt.Errorf("Topic %v does not exist", t)
	}

	delete(hs.Topics[t], p)

	return nil
}

// GetSubscribers returns the subscribers of topic t. An error is returned instead if the topic does not exist
func (hs *Storage) GetSubscribers(ctx context.Context, t Topic) ([]Processor, error) {
	// check topic
	if _, exists := hs.Topics[t]; !exists {
		return nil, fmt.Errorf("Topic %v does not exist", t)
	}

	processors := []Processor{}
	for p := range hs.Topics[t] {
		processors = append(processors, p)
	}

	return processors, nil
}

// HasSubscribed returns true if the processor p has subscribed to the topic t. An error is returned instead if the processor or the topic does not exist
func (hs *Storage) HasSubscribed(ctx context.Context, p Processor, t Topic) (bool, error) {
	// check processor
	if _, err := hs.GetProcessor(ctx, p.GetID()); err != nil {
		return false, err
	}

	// check topic
	if _, exists := hs.Topics[t]; !exists {
		return false, fmt.Errorf("Topic %v does not exist", t)
	}

	// check subscription
	return hs.Topics[t][p], nil
}

func (hs *Storage) RegisterMessageHook(ctx context.Context, hook Processor) error {
	hs.MessageHooks[hook] = true

	return nil
}

func (hs *Storage) UnregisterMessageHook(ctx context.Context, hook Processor) error {
	delete(hs.MessageHooks, hook)

	return nil
}

func (hs *Storage) GetMessageHooks(ctx context.Context) ([]Processor, error) {
	hooks := []Processor{}
	for hook := range hs.MessageHooks {
		hooks = append(hooks, hook)
	}

	return hooks, nil
}
