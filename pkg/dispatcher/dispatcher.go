package hub

import (
	"context"
	"fmt"

	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
	"github.com/gin-gonic/gin"
)

type Dispatcher struct {
	storage    DispatcherStorage
	serverInfo *HubServerinfo
	hub        hub.Hub
}

type HubServerinfo struct {
	Addresss string
}

type DispatcherStorage interface {
	GetSubscribers(ctx context.Context, topic hub.Topic) ([]*HubServerinfo, error)

	RegisterItself(ctx context.Context, topic hub.Topic) error

	UnregisterItSelf(ctx context.Context, topic hub.Topic) error
}

func (d *Dispatcher) NotifyServers(ctx context.Context, msg hub.Message) error {
	subscribers, err := d.storage.GetSubscribers(ctx, msg.GetTopic())
	if err != nil {
		return err
	}

	for _, subscriber := range subscribers {
		if err := d.notifyServer(ctx, msg, subscriber); err != nil {
			fmt.Printf("Could not notify server %s", subscriber.Addresss)
		}
	}

	return nil
}

func (d *Dispatcher) notifyServer(ctx context.Context, msg hub.Message, serverInfo *HubServerinfo) error {
	return nil
}

func (d *Dispatcher) MessageHandler() gin.HandlerFunc {
	return nil
}

type DispatcherProcessor struct {
	hub.BaseProcessor
	dispatcher *Dispatcher
}

func (dp *DispatcherProcessor) Process(msg hub.Message) error {
	if err := dp.dispatcher.storage.RegisterItself(context.TODO(), msg.GetTopic()); err != nil {
		return err
	}
	return dp.dispatcher.NotifyServers(context.TODO(), msg)
}
