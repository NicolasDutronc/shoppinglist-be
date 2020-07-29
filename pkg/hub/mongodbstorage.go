package hub

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type processorDocument struct {
	ID          primitive.ObjectID `bson:"_id"`
	ProcessorID string             `bson:"processor_id"`
}

type topicDocument struct {
	ID         primitive.ObjectID  `bson:"_id"`
	Name       string              `bson:"name"`
	Processors []processorDocument `bson:"processors"`
}

// MongoDBStorage is a mongodb-backed hub storage
type MongoDBStorage struct {
	TopicsCollection *mongo.Collection
	processors       map[string]Processor
}

// NewMongoDBStorage returns a MongoDBStorage based on a given collection
func NewMongoDBStorage(topicsCollection *mongo.Collection) Storage {
	return &MongoDBStorage{
		TopicsCollection: topicsCollection,
	}
}

// CreateTopic inserts a new topic document with the name of the topic and an empty array
func (ms *MongoDBStorage) CreateTopic(ctx context.Context, topic Topic) error {
	_, err := ms.TopicsCollection.InsertOne(ctx, topicDocument{
		Name:       string(topic),
		Processors: []processorDocument{},
	})

	return err
}

// DeleteTopic deletes a topic document based on the topic name
func (ms *MongoDBStorage) DeleteTopic(ctx context.Context, topic Topic) error {
	_, err := ms.TopicsCollection.DeleteOne(ctx, bson.M{
		"name": string(topic),
	})

	return err
}

// ListTopics returns all topics as an array
func (ms *MongoDBStorage) ListTopics(ctx context.Context) ([]Topic, error) {
	cursor, err := ms.TopicsCollection.Find(ctx, bson.M{}, &options.FindOptions{
		Projection: bson.M{
			"name": 1,
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var topicDocuments []topicDocument
	if err := cursor.Decode(&topicDocuments); err != nil {
		return nil, err
	}

	topics := make([]Topic, len(topicDocuments))

	for i, topicDoc := range topicDocuments {
		topics[i] = TopicFromString(topicDoc.Name)
	}

	return topics, nil
}

// TopicExists returns true if there is a document with the given topic as a name
func (ms *MongoDBStorage) TopicExists(ctx context.Context, topic Topic) (bool, error) {
	err := ms.TopicsCollection.FindOne(
		ctx,
		bson.M{
			"name": string(topic),
		},
	).Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// RegisterProcessor adds a processor in the processors map. An error is returned instead if the processor is already registered
func (ms *MongoDBStorage) RegisterProcessor(ctx context.Context, p Processor) error {
	if _, exists := ms.processors[p.GetID()]; exists {
		return fmt.Errorf("Processor %v is already registered", p.GetID())
	}

	ms.processors[p.GetID()] = p

	return nil
}

// UnregisterProcessor delete a processor from the processors map. An error is returned instead if the given processor is not registered
func (ms *MongoDBStorage) UnregisterProcessor(ctx context.Context, p Processor) error {
	if _, exists := ms.processors[p.GetID()]; !exists {
		return fmt.Errorf("Processor %v is not registered", p.GetID())
	}

	delete(ms.processors, p.GetID())

	return nil
}

// GetProcessor returns the processor that has the given id. An error is returned instead if no processor with the given id is found
func (ms *MongoDBStorage) GetProcessor(ctx context.Context, ID string) (Processor, error) {
	if _, exists := ms.processors[ID]; !exists {
		return nil, fmt.Errorf("Processor %v is not registered", ID)
	}

	return ms.processors[ID], nil
}

// ListProcessors returns all registered processors as a list
func (ms *MongoDBStorage) ListProcessors(ctx context.Context) ([]Processor, error) {
	processors := []Processor{}
	for _, p := range ms.processors {
		processors = append(processors, p)
	}

	return processors, nil
}

// Subscribe adds the processor to the processors array of the matching topic document
func (ms *MongoDBStorage) Subscribe(ctx context.Context, p Processor, t Topic) error {
	// check processor
	if _, err := ms.GetProcessor(ctx, p.GetID()); err != nil {
		return err
	}

	result, err := ms.TopicsCollection.UpdateOne(
		ctx,
		bson.M{
			"name": string(t),
		},
		bson.D{
			{"$addToSet", bson.D{{"processors", processorDocument{ProcessorID: p.GetID()}}}},
		},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("Topic %v does not exist", t)
	}

	return nil
}

// Unsubscribe removes the processor id from the processors array of the matching topic document
func (ms *MongoDBStorage) Unsubscribe(ctx context.Context, p Processor, t Topic) error {
	// check processor
	if _, err := ms.GetProcessor(ctx, p.GetID()); err != nil {
		return err
	}

	result, err := ms.TopicsCollection.UpdateOne(
		ctx,
		bson.M{
			"name": string(t),
		},
		bson.D{
			{"$pull", bson.D{
				{"processors", bson.M{
					"processor_id": p.GetID(),
				}},
			}},
		},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("Topic %v does not exist", t)
	}

	return nil
}

// GetSubscribers returns the processors that have subscribed to the given topic
func (ms *MongoDBStorage) GetSubscribers(ctx context.Context, t Topic) ([]Processor, error) {
	var topicDoc topicDocument
	if err := ms.TopicsCollection.FindOne(
		ctx,
		bson.M{
			"name": string(t),
		},
	).Decode(&topicDoc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("Topic %v does not exist", t)
		}
		return nil, err
	}

	processors := make([]Processor, len(topicDoc.Processors))
	for i, processorDoc := range topicDoc.Processors {
		processors[i] = ms.processors[processorDoc.ProcessorID]
	}

	return processors, nil
}

// HasSubscribed checks if the processor p is in the topic document processors array
func (ms *MongoDBStorage) HasSubscribed(ctx context.Context, p Processor, t Topic) (bool, error) {
	if err := ms.TopicsCollection.FindOne(
		ctx,
		bson.M{
			"name": string(t),
			"processors": bson.D{
				{"$in", bson.A{p.GetID()}},
			},
		},
	).Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
