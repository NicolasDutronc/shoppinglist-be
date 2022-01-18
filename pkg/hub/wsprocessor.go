package hub

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	hub Hub
}

// Process encodes the message in JSON and writes it
func (p *WebSocketProcessor) Process(msg Message) error {
	p.conn.SetWriteDeadline(time.Now().Add(p.writeWait))
	w, err := p.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	encodingErr := EncodeJSON(w, msg)
	closingErr := w.Close()
	if encodingErr != nil {
		return err
	}

	return closingErr
}

// HandleClose sends the client a close message and close the connection
func (p *WebSocketProcessor) HandleClose() {
	log.Println("Closing websocket connection")
	p.conn.WriteMessage(websocket.CloseMessage, []byte{})
	p.conn.Close()
	p.hub.UnregisterProcessor(context.TODO(), p)
}

// GetDoneChannel returns the done channel
func (p *WebSocketProcessor) GetDoneChannel() <-chan struct{} {
	return p.close
}

func (p *WebSocketProcessor) handleMessages() error {
	type subMessage struct {
		MsgType string `json:"type"`
		Topic   string `json:"topic"`
	}

	defer p.conn.Close()
	var (
		newline = []byte{'\n'}
		space   = []byte{' '}
	)

	p.conn.SetReadLimit(p.readLimit)
	//p.conn.SetReadDeadline(time.Now().Add(p.pongWait))
	// p.conn.SetPongHandler(func(string) error {
	// 	p.conn.SetReadDeadline(time.Now().Add(p.pongWait))
	// 	return nil
	// })

	for {
		messageType, messageBytes, err := p.conn.ReadMessage()
		log.Printf("message received : %v, of type %v", string(messageBytes), messageType)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected error: %v", err)
				return err
			}
			log.Printf("error: %v", err)
			break
		}
		messageBytes = bytes.TrimSpace(bytes.Replace(messageBytes, newline, space, -1))

		var message subMessage
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("error while trying to unmarshal message: %v", err)
			return err
		}

		switch message.MsgType {
		case "SUB":
			p.hub.Subscribe(context.TODO(), p, Topic(message.Topic))
		case "UNSUB":
			p.hub.Unsubscribe(context.TODO(), p, Topic(message.Topic))
		}
	}

	log.Println("exiting message handler")
	p.close <- struct{}{}

	return nil
}

// WebsocketHandler retrieves the processor id, sets up the websocket, initiates the processor and starts it
func WebsocketHandler(h Hub, writeWait time.Duration, maxMessageSize int64, pongWait time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		processorID := uuid.NewString()

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
			hub:        h,
		}

		processor.conn.SetCloseHandler(func(code int, text string) error {
			log.Println("Close msg received from client")
			processor.close <- struct{}{}
			return nil
		})

		if err := h.RegisterProcessor(c.Request.Context(), processor); err != nil {
			log.Println("Error registering", err)
			processor.close <- struct{}{}
		}

		go func() {
			if err := Run(processor); err != nil {
				log.Println("Error processing message from server", err)
				processor.close <- struct{}{}
			}
		}()

		go func() {
			if err := processor.handleMessages(); err != nil {
				log.Println("Error reading message from client", err)
				processor.close <- struct{}{}
			}
		}()
	}
}
