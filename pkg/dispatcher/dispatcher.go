package dispatcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
	"github.com/gin-gonic/gin"
)

type Dispatcher struct {
	storage           DispatcherStorage
	serverInfo        *HubServerinfo
	hub               hub.Hub
	foreignerMessages map[int64]bool
	client            *http.Client
}

type HubServerinfo struct {
	Addresss string
}

type DispatcherStorage interface {
	GetSubscribers(ctx context.Context, topic hub.Topic) ([]*HubServerinfo, error)

	Subscribe(ctx context.Context, dispatcher *Dispatcher, topic hub.Topic) error

	Unsubscribe(ctx context.Context, dispatcher *Dispatcher, topic hub.Topic) error
}

func (d *Dispatcher) notifyServers(ctx context.Context, msg hub.Message) error {
	if d.foreignerMessages[msg.GetID()] {
		delete(d.foreignerMessages, msg.GetID())
		return nil
	}
	subscribers, err := d.storage.GetSubscribers(ctx, msg.GetTopic())
	if err != nil {
		return err
	}

	for _, subscriber := range subscribers {
		if subscriber.Addresss != d.serverInfo.Addresss {
			if err := d.notifyServer(ctx, msg, subscriber); err != nil {
				fmt.Printf("Could not notify server %s", subscriber.Addresss)
			}
		}
	}

	return nil
}

func (d *Dispatcher) notifyServer(ctx context.Context, msg hub.Message, serverInfo *HubServerinfo) error {
	jsonValue, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	response, err := d.client.Post(serverInfo.Addresss, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusInternalServerError:
		return fmt.Errorf("server %s had an error handling the message", serverInfo.Addresss)
	case http.StatusBadRequest:
		return fmt.Errorf("server %s says the request was malformed", serverInfo.Addresss)
	}

	return nil
}

func (d *Dispatcher) MessageHandler() gin.HandlerFunc {
	type request struct {
		Msg hub.Message `json:"message"`
	}
	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		d.foreignerMessages[req.Msg.GetID()] = true
		if err := d.hub.Publish(c.Request.Context(), req.Msg); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusOK)
	}
}
