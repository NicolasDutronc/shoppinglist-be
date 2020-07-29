package hub

import (
	"encoding/json"
	"io"
)

// Message is an interface defining a message going through the hub
type Message interface {
	GetType() string
	GetID() int64
	GetTopic() Topic
}

// EncodeJSON provides a way to encode and send messages to json
func EncodeJSON(w io.Writer, msg Message) error {
	type payload struct {
		Type string  `json:"type"`
		ID   int64   `json:"id"`
		Msg  Message `json:"message"`
	}
	return json.NewEncoder(w).Encode(&payload{
		Type: msg.GetType(),
		ID:   msg.GetID(),
		Msg:  msg,
	})
}

// BaseMessage is an abstract struct for all messages that implements the GetID method of the message interface
type BaseMessage struct {
	ID    int64 `json:"-"`
	Topic Topic `json:"-"`
}

// GetID returns the id of the message
func (msg *BaseMessage) GetID() int64 {
	return msg.ID
}

// GetTopic returns the message topic
func (msg *BaseMessage) GetTopic() Topic {
	return msg.Topic
}

// DeleteTopicMessage is a special message that is sent to subscribers when a topic is deleted so they can safely react
type DeleteTopicMessage struct {
	BaseMessage
	Topic Topic `json:"topic"`
}

// GetType returns the type of the message
func (msg *DeleteTopicMessage) GetType() string {
	return "deleteTopicMessageType"
}
