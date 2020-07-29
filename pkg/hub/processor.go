package hub

import "log"

// Processor defines the processor interface
type Processor interface {
	GetMsgChannel() chan Message
	GetDoneChannel() <-chan struct{}
	Process(Message) error
	HandleClose()
	GetID() string
}

// BaseProcessor is a partial implementation of the processor API to implement the methods that are common to all processors
type BaseProcessor struct {
	ID       string
	messages chan Message
}

// GetID returns the id of the processor
func (p *BaseProcessor) GetID() string {
	return p.ID
}

// GetMsgChannel returns the messages channel
func (p *BaseProcessor) GetMsgChannel() chan Message {
	return p.messages
}

// Run starts the processor processing loop.
// It processes the messages from the message channel and quits if something is sent over the done channel
func Run(p Processor) error {
	log.Printf("Running processor n°%s", p.GetID())
	for {
		select {
		case msg, ok := <-p.GetMsgChannel():
			if !ok {
				return nil
			}
			if err := p.Process(msg); err != nil {
				return err
			}
		case <-p.GetDoneChannel():
			p.HandleClose()
			return nil
		}
	}
}
