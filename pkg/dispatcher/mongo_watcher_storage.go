package dispatcher

import (
	"context"
	"fmt"
	"sync"

	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var initOnce sync.Once

type MongoWatcherStorage struct {
	collection *mongo.Collection
	localState *InMemoryStorage
	mu         sync.Mutex
}

func NewMongoWatcherStorage(collection *mongo.Collection) DispatcherStorage {
	return &MongoWatcherStorage{
		collection: collection,
		localState: &InMemoryStorage{
			state: make(map[hub.Topic]map[*HubServerinfo]bool),
		},
		mu: sync.Mutex{},
	}
}

func (m *MongoWatcherStorage) GetSubscribers(ctx context.Context, topic hub.Topic) (map[*HubServerinfo]bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	initOnce.Do(func() {
		// get data from mongo
		// update internal state
		// start watching
	})

	return m.localState.GetSubscribers(ctx, topic)
}

func (m *MongoWatcherStorage) Subscribe(ctx context.Context, dispatcher *Dispatcher, topic hub.Topic) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	initOnce.Do(func() {
		m.watch(ctx)
	})

	// update mongo

	return m.localState.Subscribe(ctx, dispatcher, topic)
}

func (m *MongoWatcherStorage) Unsubscribe(ctx context.Context, dispatcher *Dispatcher, topic hub.Topic) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	initOnce.Do(func() {
		m.watch(ctx)
	})

	// update mongo

	return m.localState.Unsubscribe(ctx, dispatcher, topic)
}

func (m *MongoWatcherStorage) watch(ctx context.Context) error {
	pipeline := mongo.Pipeline{}
	stream, err := m.collection.Watch(ctx, pipeline)
	if err != nil {
		return err
	}
	defer stream.Close(ctx)

	for stream.Next(ctx) {
		var data bson.M
		if err := stream.Decode(&data); err != nil {
			return err
		}

		fmt.Printf("%v\n", data)
	}

	return nil
}
