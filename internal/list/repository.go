package list

import (
	"context"
)

// FinderByID is a single method interface for finding a list by id
type FinderByID interface {
	FindListByID(ctx context.Context, listID string) (*Shoppinglist, error)
}

// Finder is a single method interface for listing the lists
type Finder interface {
	FindAllLists(ctx context.Context) ([]*Shoppinglist, error)
}

// Creator is a single method interface for creating a list
type Creator interface {
	StoreList(ctx context.Context, listName string) (*Shoppinglist, error)
}

// Deleter is a single method interface for deleting a list
type Deleter interface {
	DeleteList(ctx context.Context, listID string) (int64, error)
}

// ItemAdder is a single method interface for adding an item to a list
type ItemAdder interface {
	AddItem(ctx context.Context, listID string, itemName string, itemQuantity string) (*Item, error)
}

// ItemUpdater is a single method interface for updating an item inside a list
type ItemUpdater interface {
	UpdateItem(ctx context.Context, listID string, itemName string, itemQuantity string, itemNewName string, itemNewQuantity string) (int64, error)
}

// ItemToggler is a single method interface for toggling an item inside a list
type ItemToggler interface {
	ToggleItem(ctx context.Context, listID string, itemName string, itemQuantity string, itemDone bool) (int64, error)
}

// ItemRemover is a single method interface for removing an item from a list
type ItemRemover interface {
	RemoveItem(ctx context.Context, listID string, itemName string, itemQuantity string) (int64, error)
}

// Clearer is a single method interface for clearing all items from a list
type Clearer interface {
	RemoveAllItems(ctx context.Context, listID string) (int64, error)
}

// Repository is a wrapper around all the single method interfaces defining the service
type Repository interface {
	FinderByID
	Finder
	Creator
	Deleter
	ItemAdder
	ItemUpdater
	ItemToggler
	ItemRemover
	Clearer
}
