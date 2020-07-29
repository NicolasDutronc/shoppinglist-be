package hub

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type payload struct {
	ID    int64
	Type  string
	Topic Topic
	Msg   Message
}

type subscriptionMessage struct {
	topic     Topic
	processor Processor
}

type unsubscriptionMessage struct {
	topic     Topic
	processor Processor
}

// Hub defines the methods used to (un)register processors, (un)subscribe processors, perform topic operations, publishing messages into the hub as well as the methods to run and close the hub
type Hub interface {
	GetProcessor(ctx context.Context, processorID string) (Processor, error)
	RegisterProcessor(ctx context.Context, p Processor) error
	UnregisterProcessor(ctx context.Context, p Processor) error
	AddTopic(ctx context.Context, topic Topic) error
	DeleteTopic(ctx context.Context, topic Topic) error
	Subscribe(ctx context.Context, p Processor, topic Topic) error
	Unsubscribe(ctx context.Context, p Processor, topic Topic) error
	Publish(ctx context.Context, msg Message) error
	Close(ctx context.Context) error
	Run(ctx context.Context, interrupt chan struct{}) error
}

// SubscriptionHandler is a gin http handler for topic subscription
func SubscriptionHandler(h Hub) gin.HandlerFunc {
	type request struct {
		Topic       string `json:"topic"`
		ProcessorID string `json:"processor"`
	}
	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		p, err := h.GetProcessor(c.Request.Context(), req.ProcessorID)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err := h.Subscribe(c.Request.Context(), p, TopicFromString(req.Topic)); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"subscription": gin.H{
				"processor": req.ProcessorID,
				"topic":     req.Topic,
			},
		})
	}
}

// UnsubscriptionHandler is a gin http handler for topic unsubscription
func UnsubscriptionHandler(h Hub) gin.HandlerFunc {
	type request struct {
		Topic       string `json:"topic"`
		ProcessorID string `json:"processor"`
	}

	return func(c *gin.Context) {
		var req request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		p, err := h.GetProcessor(c.Request.Context(), req.ProcessorID)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err := h.Unsubscribe(c.Request.Context(), p, TopicFromString(req.Topic)); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"unsubscription": gin.H{
				"processor": req.ProcessorID,
				"topic":     req.Topic,
			},
		})
	}
}
