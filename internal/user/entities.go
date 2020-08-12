package user

import (
	"fmt"
	"time"

	"github.com/NicolasDutronc/shoppinglist-be/internal/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user
type User struct {
	common.BaseModel `bson:",inline"`
	Name             string        `bson:"name" json:"name"`
	Password         string        `bson:"password"`
	Permissions      []*Permission `bson:"permissions"`
}

// NewUser is a User constructor
func NewUser(name string, password string, permissions ...*Permission) *User {
	return &User{
		BaseModel: common.BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        name,
		Password:    password,
		Permissions: permissions,
	}
}

// Can checks if the user can do the action on the resource
func (u *User) Can(action string, resourceID string) error {
	for _, p := range u.Permissions {
		if p.Action == "*" && p.ResourceID == "*" {
			return nil
		}
		if p.Action == action && p.ResourceID == resourceID {
			return nil
		}
	}

	return fmt.Errorf("User %s cannot %s on resource %s", u.Name, action, resourceID)
}

// Permission is a triplet that represents the right to do some action on a resource
type Permission struct {
	ResourceID string `bson:"resource_id"`
	Action     string `bson:"action"`
}
