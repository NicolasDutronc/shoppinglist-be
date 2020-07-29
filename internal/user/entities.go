package user

import "github.com/NicolasDutronc/shoppinglist-be/internal/common"

// User represents a user
type User struct {
	common.BaseModel `bson:",inline"`
	Name             string `bson:"name" json:"name"`
	Password         string `bson:"password" json:"password"`
}
