package hub

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

// SubscribeJSONHandler retrieves the processor id of the client, initiates the processor, registers it and starts it
func SubscribeJSONHandler(h Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		processorID := c.GetHeader("Processor-ID")
		if processorID == "" {
			c.AbortWithError(http.StatusBadRequest, errors.New("Processor-ID header was not set"))
			return
		}

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

		if err := Run(processor); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}
