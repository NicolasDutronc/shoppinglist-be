package list

import "github.com/NicolasDutronc/shoppinglist-be/internal/common"

// Item is the item model containing a name and a quantity
type Item struct {
	Name     string `bson:"name" json:"name"`
	Quantity string `bson:"quantity" json:"quantity"`
	Done     bool   `bson:"done" json:"done"`
}

// Shoppinglist is a struct defining a shoplist in the collection
type Shoppinglist struct {
	common.BaseModel `bson:",inline"`
	Name             string  `bson:"name" json:"name"`
	Items            []*Item `bson:"items" json:"items"`
}
