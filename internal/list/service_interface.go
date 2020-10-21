package list

import "context"

// Service is the interface defining the list service api
type Service interface {
	FindListByID(ctx context.Context, listID string) (*Shoppinglist, error)

	FindAllLists(ctx context.Context) ([]*Shoppinglist, error)

	StoreList(ctx context.Context, listName string) (*Shoppinglist, error)

	DeleteList(ctx context.Context, listID string) (int64, error)

	AddItem(ctx context.Context, listID string, itemName string, itemQuantity string) (*Item, error)

	UpdateItem(ctx context.Context, listID string, itemName string, itemQuantity string, itemNewName string, itemNewQuantity string) (int64, error)

	ToggleItem(ctx context.Context, listID string, itemName string, itemQuantity string, itemDone bool) (int64, error)

	RemoveItem(ctx context.Context, listID string, itemName string, itemQuantity string) (int64, error)

	RemoveAllItems(ctx context.Context, listID string) (int64, error)
}
