package hub

import "context"

type TopicCreationHook interface {
	OnTopicCreation(ctx context.Context, topic Topic) error
}

type TopicDeletionHook interface {
	OnTopicDeletion(ctx context.Context, topic Topic) error
}

type SubscriptionHook interface {
	OnSubscription(ctx context.Context, processor Processor, topic Topic) error
}
