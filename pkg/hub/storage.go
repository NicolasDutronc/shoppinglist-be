package hub

import (
	"context"
	"fmt"
)

// Storage is an interface defining how a hub can store its state, i.e. its topics and processors
type Storage interface {
	CreateTopic(ctx context.Context, topic Topic) error
	DeleteTopic(ctx context.Context, topic Topic) error
	ListTopics(ctx context.Context) ([]Topic, error)
	TopicExists(ctx context.Context, topic Topic) (bool, error)

	RegisterProcessor(ctx context.Context, p Processor) error
	UnregisterProcessor(ctx context.Context, p Processor) error
	GetProcessor(ctx context.Context, ID string) (Processor, error)
	ListProcessors(ctx context.Context) ([]Processor, error)

	Subscribe(ctx context.Context, p Processor, t Topic) error
	Unsubscribe(ctx context.Context, p Processor, t Topic) error
	GetSubscribers(ctx context.Context, t Topic) ([]Processor, error)
	HasSubscribed(ctx context.Context, p Processor, t Topic) (bool, error)
}

// InMemoryHubStorage is an in-memory hub storage
// Topics and processors are kept in memory
type InMemoryHubStorage struct {
	topics     map[Topic]map[Processor]bool
	processors map[string]Processor
}

// NewInMemoryHubStorage returns an empty InMemoryHubStorage
func NewInMemoryHubStorage() Storage {
	return &InMemoryHubStorage{
		topics:     make(map[Topic]map[Processor]bool),
		processors: make(map[string]Processor),
	}
}

// CreateTopic checks if the topic already exists and creates it. An error is returned instead if the topic already exists.
func (hs *InMemoryHubStorage) CreateTopic(ctx context.Context, topic Topic) error {
	if _, exists := hs.topics[topic]; exists {
		return fmt.Errorf("Topic %v already exists", topic)
	}

	hs.topics[topic] = make(map[Processor]bool)

	return nil
}

// DeleteTopic checks if the topic exists and deletes it. An error is returned instead if the topic does not exist.
func (hs *InMemoryHubStorage) DeleteTopic(ctx context.Context, topic Topic) error {
	if _, exists := hs.topics[topic]; !exists {
		return fmt.Errorf("Topic %v does not exist", topic)
	}

	delete(hs.topics, topic)

	return nil
}

// ListTopics returns all stored topics as a list
func (hs *InMemoryHubStorage) ListTopics(ctx context.Context) ([]Topic, error) {
	topics := []Topic{}
	for topic := range hs.topics {
		topics = append(topics, topic)
	}

	return topics, nil
}

// TopicExists returns true if the topic exists in the hub
func (hs *InMemoryHubStorage) TopicExists(ctx context.Context, topic Topic) (bool, error) {
	_, exists := hs.topics[topic]
	return exists, nil
}

// RegisterProcessor adds a processor in the processors map. An error is returned instead if the processor is already registered
func (hs *InMemoryHubStorage) RegisterProcessor(ctx context.Context, p Processor) error {
	if _, exists := hs.processors[p.GetID()]; exists {
		return fmt.Errorf("Processor %v is already registered", p.GetID())
	}

	hs.processors[p.GetID()] = p

	return nil
}

// UnregisterProcessor delete a processor from the processors map. An error is returned instead if the given processor is not registered
func (hs *InMemoryHubStorage) UnregisterProcessor(ctx context.Context, p Processor) error {
	if _, exists := hs.processors[p.GetID()]; !exists {
		return fmt.Errorf("Processor %v is not registered", p.GetID())
	}

	delete(hs.processors, p.GetID())

	return nil
}

// GetProcessor returns the processor that has the given id. An error is returned instead if no processor with the given id is found
func (hs *InMemoryHubStorage) GetProcessor(ctx context.Context, ID string) (Processor, error) {
	if _, exists := hs.processors[ID]; !exists {
		return nil, fmt.Errorf("Processor %v is not registered", ID)
	}

	return hs.processors[ID], nil
}

// ListProcessors returns all registered processors as a list
func (hs *InMemoryHubStorage) ListProcessors(ctx context.Context) ([]Processor, error) {
	processors := []Processor{}
	for _, p := range hs.processors {
		processors = append(processors, p)
	}

	return processors, nil
}

// Subscribe adds the processor p to the subscribers of topic t. An error is returned instead if the processor or the topic does not exist
func (hs *InMemoryHubStorage) Subscribe(ctx context.Context, p Processor, t Topic) error {
	// check processor
	if _, err := hs.GetProcessor(ctx, p.GetID()); err != nil {
		return err
	}

	// check topic
	if _, exists := hs.topics[t]; !exists {
		return fmt.Errorf("Topic %v does not exist", t)
	}

	hs.topics[t][p] = true

	return nil
}

// Unsubscribe deletes processor p from the subscribers of t. An error is returned instead if the processor or the topic does not exist
func (hs *InMemoryHubStorage) Unsubscribe(ctx context.Context, p Processor, t Topic) error {
	// check processor
	if _, err := hs.GetProcessor(ctx, p.GetID()); err != nil {
		return err
	}

	// check topic
	if _, exists := hs.topics[t]; !exists {
		return fmt.Errorf("Topic %v does not exist", t)
	}

	delete(hs.topics[t], p)

	return nil
}

// GetSubscribers returns the subscribers of topic t. An error is returned instead if the topic does not exist
func (hs *InMemoryHubStorage) GetSubscribers(ctx context.Context, t Topic) ([]Processor, error) {
	// check topic
	if _, exists := hs.topics[t]; !exists {
		return nil, fmt.Errorf("Topic %v does not exist", t)
	}

	processors := []Processor{}
	for p := range hs.topics[t] {
		processors = append(processors, p)
	}

	return processors, nil
}

// HasSubscribed returns true if the processor p has subscribed to the topic t. An error is returned instead if the processor or the topic does not exist
func (hs *InMemoryHubStorage) HasSubscribed(ctx context.Context, p Processor, t Topic) (bool, error) {
	// check processor
	if _, err := hs.GetProcessor(ctx, p.GetID()); err != nil {
		return false, err
	}

	// check topic
	if _, exists := hs.topics[t]; !exists {
		return false, fmt.Errorf("Topic %v does not exist", t)
	}

	// check subscription
	return hs.topics[t][p], nil
}
