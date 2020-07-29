package common

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel is the base for all models containing common information like id, creation and update time
type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
