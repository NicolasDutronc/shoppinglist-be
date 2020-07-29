package list

import "github.com/NicolasDutronc/shoppinglist-be/pkg/hub"

type newListMessage struct {
	hub.BaseMessage
	NewList *Shoppinglist `json:"new_list"`
}

func (msg *newListMessage) GetType() string {
	return "newListMessageType"
}

type deleteListMessage struct {
	hub.BaseMessage
	ListID string `json:"listID"`
}

func (msg *deleteListMessage) GetType() string {
	return "deleteListMessageType"
}

type addItemMessage struct {
	hub.BaseMessage
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Done     bool   `json:"done"`
}

func (msg *addItemMessage) GetType() string {
	return "addItemMessageType"
}

type updateItemMessage struct {
	hub.BaseMessage
	Name        string `json:"name"`
	Quantity    string `json:"quantity"`
	NewName     string `json:"new_name"`
	NewQuantity string `json:"new_quantity"`
}

func (msg *updateItemMessage) GetType() string {
	return "updateItemMessageType"
}

type toggleItemMessage struct {
	hub.BaseMessage
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Value    bool   `json:"value"`
}

func (msg *toggleItemMessage) GetType() string {
	return "toggleItemMessageType"
}

type deleteItemMessage struct {
	hub.BaseMessage
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
}

func (msg *deleteItemMessage) GetType() string {
	return "deleteItemMessageType"
}

type clearListMesssage struct {
	hub.BaseMessage
	ListID string `json:"listID"`
}

func (msg *clearListMesssage) GetType() string {
	return "clearListMessageType"
}
