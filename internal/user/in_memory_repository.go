package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InMemoryRepository is a user repository based on mongodb
type InMemoryRepository struct {
	UserCollection []*User
}

// NewInMemoryRepository initis a new mongodb user repository
func NewInMemoryRepository() Repository {
	return &InMemoryRepository{
		UserCollection: []*User{},
	}
}

// FindByID returns the first user that matches the given id
func (r *InMemoryRepository) FindByID(ctx context.Context, userID string) (*User, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	for _, user := range r.UserCollection {
		if user.ID == objectID {
			return user, nil
		}
	}

	return nil, fmt.Errorf("User not found for id %s", userID)
}

// FindByName returns the first user that matches the given name
func (r *InMemoryRepository) FindByName(ctx context.Context, userName string) (*User, error) {
	for _, user := range r.UserCollection {
		if user.Name == userName {
			return user, nil
		}
	}

	return nil, fmt.Errorf("User not found with name %s", userName)
}

// Store creates a new user and stores it
func (r *InMemoryRepository) Store(ctx context.Context, name string, password string, permissions ...*Permission) (*User, error) {
	user := NewUser(name, password, permissions...)

	r.UserCollection = append(r.UserCollection, user)

	return user, nil
}

// UpdateName updates the name of the first user that matches the given id
func (r *InMemoryRepository) UpdateName(ctx context.Context, userID string, newName string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return -1, err
	}

	for idx, user := range r.UserCollection {
		if user.ID == objectID {
			r.UserCollection[idx].Name = newName
			return 1, nil
		}
	}

	return -1, fmt.Errorf("User not found for id %s", userID)
}

// UpdatePassword updates the password
func (r *InMemoryRepository) UpdatePassword(ctx context.Context, userID string, newPassword string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return -1, err
	}

	for idx, user := range r.UserCollection {
		if user.ID == objectID {
			r.UserCollection[idx].Password = newPassword
			return 1, nil
		}
	}

	return -1, fmt.Errorf("User not found for id %s", userID)

}

// Delete deletes a user
func (r *InMemoryRepository) Delete(ctx context.Context, userID string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return -1, err
	}

	for idx, user := range r.UserCollection {
		if user.ID == objectID {
			r.UserCollection = append(r.UserCollection[:idx], r.UserCollection[idx+1:]...)
			return 1, nil
		}
	}

	return -1, fmt.Errorf("User not found for id %s", userID)
}

// AddPermissions adds the permissions to the user
func (r *InMemoryRepository) AddPermissions(ctx context.Context, userID string, permissions ...*Permission) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return -1, err
	}

	for idx, user := range r.UserCollection {
		if user.ID == objectID {
			r.UserCollection[idx].Permissions = append(r.UserCollection[idx].Permissions, permissions...)
			return 1, nil
		}
	}

	return -1, fmt.Errorf("User not found for id %s", userID)
}

// RemovePermissions adds the permissions to the user
func (r *InMemoryRepository) RemovePermissions(ctx context.Context, userID string, permissions ...*Permission) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return -1, err
	}

	for idx, user := range r.UserCollection {
		if user.ID == objectID {
			userPermissions := make([]*Permission, len(user.Permissions))
			copy(userPermissions, user.Permissions)
			for i, userPermission := range userPermissions {
				for _, permissionToRemove := range permissions {
					if userPermission == permissionToRemove {
						r.UserCollection[idx].Permissions = append(r.UserCollection[idx].Permissions[:i], r.UserCollection[idx].Permissions[i+1:]...)
					}
				}
			}
			return 1, nil
		}
	}

	return -1, fmt.Errorf("User not found for id %s", userID)
}
