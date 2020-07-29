package hub

// Topic represents a channel to which you can publish and subscribe
type Topic string

func emptyTopic() Topic {
	return ""
}

// TopicFromString returns a Topic based on an id
func TopicFromString(id string) Topic {
	return Topic(id)
}

type topicOp int

type topicOpMessage struct {
	topic Topic
	op    topicOp
}

const (
	topicCreation topicOp = iota
	topicDeletion
)
