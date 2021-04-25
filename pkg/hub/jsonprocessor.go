package hub

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// JSONProcessor sends its messages as JSON payloads
type JSONProcessor struct {
	BaseProcessor
	ginContext *gin.Context
}

// Process encodes the message as a JSON payload and send it to the client
func (p *JSONProcessor) Process(msg Message) error {
	if err := EncodeJSON(p.ginContext.Writer, msg); err != nil {
		log.Printf("Error occured encoding message of type %v : %v", msg.GetType(), err.Error())
		p.ginContext.AbortWithError(http.StatusInternalServerError, err)
		return err
	}
	p.ginContext.Status(http.StatusOK)
	p.ginContext.Writer.Flush()
	log.Println("Message sent")

	return nil
}

// HandleClose sends a http status ok to the client to indicate a normal close
func (p *JSONProcessor) HandleClose() {
	p.ginContext.AbortWithStatus(http.StatusOK)
}

// GetDoneChannel returns the client request done channel which is close when the client closes the request
func (p *JSONProcessor) GetDoneChannel() <-chan struct{} {
	return p.ginContext.Request.Context().Done()
}

type connectionMessage struct {
	BaseMessage
	ProcessorID string `json:"processor-id"`
}

func (msg *connectionMessage) GetType() string {
	return "connectionMessage"
}

// SubscribeJSONHandler retrieves the processor id of the client, initiates the processor, registers it and starts it
func SubscribeJSONHandler(h Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		processorID := uuid.NewString()

		processor := &JSONProcessor{
			BaseProcessor: BaseProcessor{
				ID:       processorID,
				messages: make(chan Message),
			},
			ginContext: c,
		}

		if err := h.RegisterProcessor(c.Request.Context(), processor); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer func() {
			if err := h.UnregisterProcessor(c.Request.Context(), processor); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}()

		processor.Process(&connectionMessage{
			BaseMessage: BaseMessage{
				ID:    time.Now().Unix(),
				Topic: TopicFromString("internal"),
			},
			ProcessorID: processorID,
		})

		if err := Run(processor); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
