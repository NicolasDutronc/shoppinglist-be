package list

import (
	"context"
	"time"

	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
)

// Service is a wrapper around all the single method interfaces defining the service
type Service interface {
	Repository
}

// ServiceImpl is the implementation of the Service interface
type ServiceImpl struct {
	repository Repository
	h          hub.Hub
}

// NewService returns a Shoppinglist service based on a shoplist repository, and a hub
func NewService(repo Repository, h hub.Hub) Service {
	return &ServiceImpl{
		repository: repo,
		h:          h,
	}
}

// FindListByID retrieves a list based on its id
func (s *ServiceImpl) FindListByID(ctx context.Context, listID string) (*Shoppinglist, error) {
	return s.repository.FindListByID(ctx, listID)
}

// FindAllLists retrieves all lists
func (s *ServiceImpl) FindAllLists(ctx context.Context) ([]*Shoppinglist, error) {
	return s.repository.FindAllLists(ctx)
}

// StoreList inserts a new empty list
func (s *ServiceImpl) StoreList(ctx context.Context, listName string) (*Shoppinglist, error) {
	list, err := s.repository.StoreList(ctx, listName)
	if err != nil {
		return nil, err
	}

	if err := s.h.AddTopic(ctx, hub.TopicFromString(list.ID.Hex())); err != nil {
		return nil, err
	}

	if err := s.h.Publish(ctx, &newListMessage{
		BaseMessage: hub.BaseMessage{
			ID:    time.Now().Unix(),
			Topic: hub.TopicFromString("lists"),
		},
		NewList: list,
	}); err != nil {
		return nil, err
	}

	return list, nil
}

// DeleteList removes a list
func (s *ServiceImpl) DeleteList(ctx context.Context, listID string) (int64, error) {
	n, err := s.repository.DeleteList(ctx, listID)
	if err != nil {
		return -1, err
	}

	if err := s.h.DeleteTopic(ctx, hub.TopicFromString(listID)); err != nil {
		return -1, err
	}

	if err := s.h.Publish(ctx, &deleteListMessage{
		BaseMessage: hub.BaseMessage{
			ID:    time.Now().Unix(),
			Topic: hub.TopicFromString("lists"),
		},
		ListID: listID,
	}); err != nil {
		return -1, err
	}

	return n, nil
}

// AddItem adds a new item to a list given by its id
func (s *ServiceImpl) AddItem(ctx context.Context, listID string, itemName string, itemQuantity string) (*Item, error) {
	item, err := s.repository.AddItem(ctx, listID, itemName, itemQuantity)
	if err != nil {
		return nil, err
	}

	if err := s.h.Publish(ctx, &addItemMessage{
		BaseMessage: hub.BaseMessage{
			ID:    time.Now().Unix(),
			Topic: hub.TopicFromString(listID),
		},
		Name:     itemName,
		Quantity: itemQuantity,
		Done:     false,
	}); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateItem updates an item based on its name and quantity in a list given its id
func (s *ServiceImpl) UpdateItem(ctx context.Context, listID string, itemName string, itemQuantity string, itemNewName string, itemNewQuantity string) (int64, error) {
	n, err := s.repository.UpdateItem(ctx, listID, itemName, itemQuantity, itemNewName, itemNewQuantity)
	if err != nil {
		return -1, err
	}

	if err := s.h.Publish(ctx, &updateItemMessage{
		BaseMessage: hub.BaseMessage{
			ID:    time.Now().Unix(),
			Topic: hub.TopicFromString(listID),
		},
		Name:        itemName,
		Quantity:    itemQuantity,
		NewName:     itemNewName,
		NewQuantity: itemNewQuantity,
	}); err != nil {
		return -1, err
	}

	return n, err
}

// ToggleItem changes the done boolean value of an item
func (s *ServiceImpl) ToggleItem(ctx context.Context, listID string, itemName string, itemQuantity string, itemDone bool) (int64, error) {
	n, err := s.repository.ToggleItem(ctx, listID, itemName, itemQuantity, itemDone)
	if err != nil {
		return -1, nil
	}

	if err := s.h.Publish(ctx, &toggleItemMessage{
		BaseMessage: hub.BaseMessage{
			ID:    time.Now().Unix(),
			Topic: hub.TopicFromString(listID),
		},
		Name:     itemName,
		Quantity: itemQuantity,
		Value:    itemDone,
	}); err != nil {
		return -1, err
	}

	return n, nil
}

// RemoveItem removes an item from a list
func (s *ServiceImpl) RemoveItem(ctx context.Context, listID string, itemName string, itemQuantity string) (int64, error) {
	n, err := s.repository.RemoveItem(ctx, listID, itemName, itemQuantity)
	if err != nil {
		return -1, err
	}

	if err := s.h.Publish(ctx, &deleteItemMessage{
		BaseMessage: hub.BaseMessage{
			ID:    time.Now().Unix(),
			Topic: hub.TopicFromString(listID),
		},
		Name:     itemName,
		Quantity: itemQuantity,
	}); err != nil {
		return -1, err
	}

	return n, nil
}

// RemoveAllItems removes all items from a list
func (s *ServiceImpl) RemoveAllItems(ctx context.Context, listID string) (int64, error) {
	n, err := s.repository.RemoveAllItems(ctx, listID)
	if err != nil {
		return -1, err
	}

	if err := s.h.Publish(ctx, &clearListMesssage{
		BaseMessage: hub.BaseMessage{
			ID:    time.Now().Unix(),
			Topic: hub.TopicFromString(listID),
		},
		ListID: listID,
	}); err != nil {
		return -1, err
	}

	return n, nil
}
