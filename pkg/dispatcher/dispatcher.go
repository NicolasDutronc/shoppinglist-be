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
	Address string
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

	for subscriber := range subscribers {
		if subscriber.Address != d.serverInfo.Address {
			if err := d.notifyServer(ctx, msg, subscriber); err != nil {
				fmt.Printf("Could not notify server %s", subscriber.Address)
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

	response, err := d.client.Post(serverInfo.Address, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusInternalServerError:
		return fmt.Errorf("server %s had an error handling the message", serverInfo.Address)
	case http.StatusBadRequest:
		return fmt.Errorf("server %s says the request was malformed", serverInfo.Address)
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
