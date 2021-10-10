package stream

import "github.com/aws/aws-sdk-go/service/sqs"

// SQS An interface for SQS operations
type SQS interface {
	Poll()
	Delete(msg *sqs.Message) error
	Enqueue(msgBatch []*sqs.SendMessageBatchRequestEntry) error
	RegisterPollHandler(pollHandler func(msg *sqs.Message))
	ChangeVisibilityTimeout(msg *sqs.Message, seconds int64) bool
}

