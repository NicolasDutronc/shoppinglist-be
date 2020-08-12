package user

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDBRepository is a user repository based on mongodb
type MongoDBRepository struct {
	UserCollection *mongo.Collection
}

// NewMongoDBRepository initis a new mongodb user repository
func NewMongoDBRepository(collection *mongo.Collection) Repository {
	return &MongoDBRepository{
		UserCollection: collection,
	}
}

// FindByID returns the first user that matches the given id
func (r *MongoDBRepository) FindByID(ctx context.Context, userID string) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var user User
	if err := r.UserCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByName returns the first user that matches the given name
func (r *MongoDBRepository) FindByName(ctx context.Context, userName string) (*User, error) {
	var user User
	if err := r.UserCollection.FindOne(ctx, bson.M{"name": userName}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Store creates a new user and stores it
func (r *MongoDBRepository) Store(ctx context.Context, name string, password string, permissions ...*Permission) (*User, error) {
	user := NewUser(name, password, permissions...)

	if _, err := r.UserCollection.InsertOne(ctx, *user); err != nil {
		return nil, err
	}

	return user, nil

}

// UpdateName updates the name of the first user that matches the given id
func (r *MongoDBRepository) UpdateName(ctx context.Context, userID string, newName string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return -1, err
	}

	result, err := r.UserCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"name", newName}}},
			{"$set", bson.D{{"updated_at", time.Now()}}},
		},
	)

	if err != nil {
		return -1, err
	}

	return result.ModifiedCount, nil
}

// UpdatePassword updates the password
func (r *MongoDBRepository) UpdatePassword(ctx context.Context, userID string, currentPassword string, newPassword string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return -1, err
	}

	result, err := r.UserCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.D{
			{"$set", bson.D{{"name", newPassword}}},
			{"$set", bson.D{{"updated_at", time.Now()}}},
		},
	)

	if err != nil {
		return -1, err
	}

	return result.ModifiedCount, nil

}

// Delete deletes a user
func (r *MongoDBRepository) Delete(ctx context.Context, userID string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return -1, err
	}

	result, err := r.UserCollection.DeleteOne(
		ctx,
		bson.M{"_id": objectID},
	)

	if err != nil {
		return -1, err
	}

	return result.DeletedCount, nil
}
