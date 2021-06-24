package list

import (
	"context"
	"fmt"
	"time"

	"github.com/NicolasDutronc/shoppinglist-be/internal/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InMemoryRepository is an in-memory shoplist repository
type InMemoryRepository struct {
	lists map[string]*Shoppinglist
}

// NewInMemoryRepository is a constructor of InMemoryRepository
func NewInMemoryRepository() Repository {
	return &InMemoryRepository{
		lists: make(map[string]*Shoppinglist),
	}
}

// FindListByID retrieves a list based on its id
func (r *InMemoryRepository) FindListByID(ctx context.Context, listID string) (*Shoppinglist, error) {
	list, exists := r.lists[listID]
	if !exists {
		return nil, fmt.Errorf("there is no list with id %v", listID)
	}

	return list, nil
}

// FindAllLists retrieves all lists
func (r *InMemoryRepository) FindAllLists(ctx context.Context) ([]*Shoppinglist, error) {
	var lists []*Shoppinglist
	for _, list := range r.lists {
		lists = append(lists, list)
	}

	return lists, nil
}

// StoreList inserts a new empty list
func (r *InMemoryRepository) StoreList(ctx context.Context, listName string) (*Shoppinglist, error) {
	exists := false
	for _, list := range r.lists {
		if list.Name == listName {
			exists = true
			break
		}
	}

	if exists {
		return nil, fmt.Errorf("A list already exists with the name %v", listName)
	}

	newList := &Shoppinglist{
		BaseModel: common.BaseModel{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Items: []*Item{},
		Name:  listName,
	}

	r.lists[newList.ID.Hex()] = newList

	return newList, nil
}

// DeleteList removes a list
func (r *InMemoryRepository) DeleteList(ctx context.Context, listID string) (int64, error) {
	list, err := r.FindListByID(ctx, listID)
	if err != nil {
		return -1, err
	}

	delete(r.lists, list.ID.Hex())

	return 1, nil
}

// AddItem adds a new item to a list given by its id
func (r *InMemoryRepository) AddItem(ctx context.Context, listID string, itemName string, itemQuantity string) (*Item, error) {
	list, err := r.FindListByID(ctx, listID)
	if err != nil {
		return nil, err
	}

	newitem := &Item{
		Name:     itemName,
		Quantity: itemQuantity,
		Done:     false,
	}

	list.Items = append(list.Items, newitem)
	list.UpdatedAt = time.Now()

	return newitem, nil
}

// UpdateItem updates an item based on its name and quantity in a list given its id
func (r *InMemoryRepository) UpdateItem(ctx context.Context, listID string, itemName string, itemQuantity string, itemNewName string, itemNewQuantity string) (int64, error) {
	list, err := r.FindListByID(ctx, listID)
	if err != nil {
		return -1, err
	}

	found := false

	for _, item := range list.Items {
		if item.Name == itemName && item.Quantity == itemQuantity {
			item.Name = itemNewName
			item.Quantity = itemNewQuantity
			found = true
			break
		}
	}

	if !found {
		return -1, fmt.Errorf("Could not find any item matching the name %v and the quantity %v in the list %v", itemName, itemQuantity, listID)
	}

	list.UpdatedAt = time.Now()

	return 1, nil
}

// ToggleItem changes the done boolean value of an item
func (r *InMemoryRepository) ToggleItem(ctx context.Context, listID string, itemName string, itemQuantity string, itemDone bool) (int64, error) {
	list, err := r.FindListByID(ctx, listID)
	if err != nil {
		return -1, err
	}

	found := false

	for _, item := range list.Items {
		if item.Name == itemName && item.Quantity == itemQuantity {
			item.Done = itemDone
			found = true
			break
		}
	}

	if !found {
		return -1, fmt.Errorf("Could not find any item matching the name %v and the quantity %v in the list %v", itemName, itemQuantity, listID)
	}

	list.UpdatedAt = time.Now()

	return 1, nil
}

// RemoveItem removes an item from a list
func (r *InMemoryRepository) RemoveItem(ctx context.Context, listID string, itemName string, itemQuantity string) (int64, error) {
	list, err := r.FindListByID(ctx, listID)
	if err != nil {
		return -1, err
	}

	index := -1

	for i, item := range list.Items {
		if item.Name == itemName && item.Quantity == itemQuantity {
			index = i
			break
		}
	}

	if index == -1 {
		return -1, fmt.Errorf("Could not find any item matching the name %v and the quantity %v in the list %v", itemName, itemQuantity, listID)
	}

	list.Items = append(list.Items[:index], list.Items[index+1:]...)
	list.UpdatedAt = time.Now()

	return 1, nil
}

// RemoveAllItems removes all items from a list
func (r *InMemoryRepository) RemoveAllItems(ctx context.Context, listID string) (int64, error) {
	list, err := r.FindListByID(ctx, listID)
	if err != nil {
		return -1, err
	}

	n := len(list.Items)
	list.Items = []*Item{}

	list.UpdatedAt = time.Now()

	return int64(n), nil
}
