package dispatcher

import (
	"context"
	"strings"

	"github.com/NicolasDutronc/shoppinglist-be/pkg/hub"
)

type DispatcherProcessor struct {
	hub.BaseProcessor
	dispatcher  *Dispatcher
	doneChannel <-chan struct{}
}

func (dp *DispatcherProcessor) Process(msg hub.Message) error {
	if strings.HasPrefix(string(msg.GetTopic()), "_internal") {
		return nil
	}
	return dp.dispatcher.NotifyServers(context.TODO(), msg)
}

func (dp *DispatcherProcessor) GetDoneChannel() <-chan struct{} {
	return dp.doneChannel
}

func (dp *DispatcherProcessor) HandleClose() {

}

type SubscriptionProcessor struct {
	hub.BaseProcessor
	dispatcher    *Dispatcher
	doneChannel   <-chan struct{}
	subscriptions map[hub.Topic]int
}

func (sp *SubscriptionProcessor) Process(msg hub.Message) error {
	switch msg.GetTopic() {
	case hub.SUBSCRIPTION_TOPIC:
		sub := msg.(hub.SubscriptionMessage)
		if val, ok := sp.subscriptions[sub.Topic]; ok {
			sp.subscriptions[sub.Topic] = val + 1
		} else {
			sp.subscriptions[sub.Topic] = 1

			if err := sp.dispatcher.storage.Subscribe(context.TODO(), sp.dispatcher, sub.Topic); err != nil {
				return err
			}
		}
		return nil
	case hub.UNSUBSCRIPTION_TOPIC:
		sub := msg.(hub.UnsubscriptionMessage)
		if val, ok := sp.subscriptions[sub.Topic]; ok {
			sp.subscriptions[sub.Topic] = val - 1

			if val == 1 {
				if err := sp.dispatcher.storage.Unsubscribe(context.TODO(), sp.dispatcher, sub.Topic); err != nil {
					return err
				}
				delete(sp.subscriptions, sub.Topic)
			}
		}
		return nil
	}

	return nil
}

func (dp *SubscriptionProcessor) GetDoneChannel() <-chan struct{} {
	return dp.doneChannel
}

func (dp *SubscriptionProcessor) HandleClose() {

}
