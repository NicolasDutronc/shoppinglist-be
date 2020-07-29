package hub

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketProcessor is a processor that sends its messages through a websocket
type WebSocketProcessor struct {
	BaseProcessor

	close chan struct{}
	conn  *websocket.Conn

	readLimit  int64
	pongWait   time.Duration
	pingPeriod time.Duration
	writeWait  time.Duration
}

// Process encodes the message in JSON and writes it
func (p *WebSocketProcessor) Process(msg Message) error {
	p.conn.SetWriteDeadline(time.Now().Add(p.writeWait))
	w, err := p.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if err := EncodeJSON(w, msg); err != nil {
		return err
	}

	return w.Close()
}

// HandleClose sends the client a close message and close the connection
func (p *WebSocketProcessor) HandleClose() {
	p.conn.WriteMessage(websocket.CloseMessage, []byte{})
	p.conn.Close()
}

// GetDoneChannel returns the done channel
func (p *WebSocketProcessor) GetDoneChannel() <-chan struct{} {
	return p.close
}

// WebsocketHandler retrieves the processor id, sets up the websocket, initiates the processor and starts it
func WebsocketHandler(h Hub, writeWait time.Duration, maxMessageSize int64, pongWait time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		processorID := c.Request.URL.Query().Get("processor")

		upgrader := &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		processor := &WebSocketProcessor{
			BaseProcessor: BaseProcessor{
				ID:       processorID,
				messages: make(chan Message),
			},
			close:      make(chan struct{}),
			conn:       conn,
			pongWait:   pongWait,
			pingPeriod: (pongWait * 9) / 10,
			readLimit:  maxMessageSize,
			writeWait:  writeWait,
		}

		processor.conn.SetReadLimit(processor.readLimit)
		processor.conn.SetReadDeadline(time.Now().Add(processor.pongWait))
		processor.conn.SetPongHandler(func(string) error {
			processor.conn.SetReadDeadline(time.Now().Add(processor.pongWait))
			return nil
		})
		processor.conn.SetCloseHandler(func(code int, text string) error {
			processor.close <- struct{}{}
			return nil
		})

		if err := h.RegisterProcessor(c.Request.Context(), processor); err != nil {
			processor.close <- struct{}{}
		}

		go func() {
			if err := Run(processor); err != nil {
				processor.close <- struct{}{}
			}
		}()
	}
}
