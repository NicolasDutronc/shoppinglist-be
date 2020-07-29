package list

import (
	"context"
	"time"

	"github.com/NicolasDutronc/shoppinglist-be/internal/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDBRepository contains all the methods to interact with the shoplist collection
type MongoDBRepository struct {
	ShoppinglistsCollection *mongo.Collection
}

// NewMongoDBRepository is a constructor for MongoDBRepository
func NewMongoDBRepository(coll *mongo.Collection) Repository {
	return &MongoDBRepository{
		ShoppinglistsCollection: coll,
	}
}

// FindListByID retrieves a list based on its id
func (r *MongoDBRepository) FindListByID(ctx context.Context, id string) (*Shoppinglist, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var list Shoppinglist
	if err := r.ShoppinglistsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&list); err != nil {
		return nil, err
	}

	return &list, nil
}

// FindAllLists retrieves all lists
func (r *MongoDBRepository) FindAllLists(ctx context.Context) ([]*Shoppinglist, error) {
	var lists []*Shoppinglist
	cursor, err := r.ShoppinglistsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &lists); err != nil {
		return nil, err
	}

	return lists, err
}

// StoreList inserts a new empty list
func (r *MongoDBRepository) StoreList(ctx context.Context, name string) (*Shoppinglist, error) {
	list := Shoppinglist{
		BaseModel: common.BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:  name,
		Items: []*Item{},
	}

	_, err := r.ShoppinglistsCollection.InsertOne(ctx, list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

// DeleteList removes a list
func (r *MongoDBRepository) DeleteList(ctx context.Context, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	result, err := r.ShoppinglistsCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return -1, err
	}

	return result.DeletedCount, nil
}

// AddItem adds a new item to a list given by its id
func (r *MongoDBRepository) AddItem(ctx context.Context, id string, name string, quantity string) (*Item, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	newItem := Item{
		Name:     name,
		Quantity: quantity,
		Done:     false,
	}
	_, err = r.ShoppinglistsCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.D{
			{"$addToSet", bson.D{{"items", newItem}}},
			{"$set", bson.D{{"updated_at", time.Now()}}},
		},
	)
	if err != nil {
		return nil, err
	}

	return &newItem, nil
}

// UpdateItem updates an item based on its name and quantity in a list given its id
func (r *MongoDBRepository) UpdateItem(ctx context.Context, id string, name string, quantity string, newName string, newQuantity string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	result, err := r.ShoppinglistsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":            objectID,
			"items.name":     name,
			"items.quantity": quantity,
		},
		bson.D{
			{"$set", bson.D{{"items.$.name", newName}}},
			{"$set", bson.D{{"items.$.quantity", newQuantity}}},
			{"$set", bson.D{{"updated_at", time.Now()}}},
		},
	)

	if err != nil {
		return -1, err
	}

	return result.ModifiedCount, nil
}

// ToggleItem changes the done boolean value of an item
func (r *MongoDBRepository) ToggleItem(ctx context.Context, id string, name string, quantity string, done bool) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	result, err := r.ShoppinglistsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":            objectID,
			"items.name":     name,
			"items.quantity": quantity,
		},
		bson.D{
			{"$set", bson.D{{"items.$.done", done}}},
			{"$set", bson.D{{"updated_at", time.Now()}}},
		},
	)
	if err != nil {
		return -1, err
	}

	return result.ModifiedCount, nil
}

// RemoveItem removes an item from a list
func (r *MongoDBRepository) RemoveItem(ctx context.Context, id string, name string, quantity string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	result, err := r.ShoppinglistsCollection.UpdateOne(
		ctx,
		bson.M{
			"_id": objectID,
		},
		bson.D{
			{"$pull", bson.D{
				{"items", bson.D{
					{"name", name},
					{"quantity", quantity},
				}},
			}},
			{"$set", bson.D{{"updated_at", time.Now()}}},
		},
	)
	if err != nil {
		return -1, err
	}

	return result.ModifiedCount, nil
}

// RemoveAllItems removes all items from a list
func (r *MongoDBRepository) RemoveAllItems(ctx context.Context, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, err
	}

	result, err := r.ShoppinglistsCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.D{
			{"$pull", bson.D{{"items", bson.D{}}}},
			{"$set", bson.D{{"updated_at", time.Now()}}},
		},
	)
	if err != nil {
		return -1, err
	}

	return result.ModifiedCount, nil
}
