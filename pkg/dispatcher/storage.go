package dispatcher

import (
	"context"
	"fmt"

	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
)

type DispatcherStorage interface {
	GetSubscribers(ctx context.Context, topic hub.Topic) (map[*HubServerinfo]bool, error)

	Subscribe(ctx context.Context, dispatcher *Dispatcher, topic hub.Topic) error

	Unsubscribe(ctx context.Context, dispatcher *Dispatcher, topic hub.Topic) error
}

type InMemoryStorage struct {
	state map[hub.Topic]map[*HubServerinfo]bool
}

func NewInMemoryStorage() DispatcherStorage {
	return &InMemoryStorage{
		state: make(map[hub.Topic]map[*HubServerinfo]bool),
	}
}

func (ms *InMemoryStorage) GetSubscribers(ctx context.Context, topic hub.Topic) (map[*HubServerinfo]bool, error) {
	if element, ok := ms.state[topic]; ok {
		return element, nil
	}
	return nil, fmt.Errorf("no such topic : %s", topic)
}

func (ms *InMemoryStorage) Subscribe(ctx context.Context, dispatcher *Dispatcher, topic hub.Topic) error {
	ms.state[topic][dispatcher.serverInfo] = true
	return nil
}

func (ms *InMemoryStorage) Unsubscribe(ctx context.Context, dispatcher *Dispatcher, topic hub.Topic) error {
	delete(ms.state[topic], dispatcher.serverInfo)
	return nil
}
