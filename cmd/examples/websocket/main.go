package main

import (
	"context"
	"time"

	"github.com/NicolasDutronc/shoppinglist-be/internal/api"
	"github.com/NicolasDutronc/shoppinglist-be/internal/list"
	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
	"github.com/gin-gonic/gin"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type addItemMessage struct {
	hub.BaseMessage
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Done     bool   `json:"done"`
}

func (msg *addItemMessage) GetType() string {
	return "addItemMessageType"
}

type mockListItemAdder struct {
	h hub.Hub
}

func (m *mockListItemAdder) AddItem(ctx context.Context, listID string, name string, quantity string) (*list.Item, error) {
	if err := m.h.Publish(ctx, &addItemMessage{
		BaseMessage: hub.NewBaseMessage(time.Now().Unix(), hub.TopicFromString(listID)),
		Name:        name,
		Quantity:    quantity,
		Done:        false,
	}); err != nil {
		return nil, err
	}
	return &list.Item{
		Name:     name,
		Quantity: quantity,
		Done:     false,
	}, nil
}

func main() {
	ctx := context.Background()
	// Create an interruption channel
	quit := make(chan struct{}, 1)

	storage := hub.NewStorage()
	h, err := hub.NewChannelHub(ctx, storage, hub.TopicFromString("1"), hub.TopicFromString("2"))
	if err != nil {
		panic(err)
	}
	go h.Run(ctx, quit)
	defer h.Close(ctx)
	h.AddTopic(ctx, hub.TopicFromString("3"))

	r := gin.Default()

	r.StaticFile("/", "./cmd/examples/websocket/home.html")

	r.GET("/connect", hub.WebsocketHandler(h, writeWait, maxMessageSize, pongWait))
	r.POST("/subscribe", hub.SubscriptionHandler(h))
	r.POST("/unsubscribe", hub.UnsubscriptionHandler(h))
	r.POST("/lists/:id/send", api.AddItemHandler(&mockListItemAdder{h: h}))

	r.RunTLS(":8080", "server.crt", "server.key")
}
