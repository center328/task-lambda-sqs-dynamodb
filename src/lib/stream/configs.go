package stream

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"sync"
)

// Config Wrapper for Config methods
type Config struct {

	AWSKey    string
	AWSSecret string
	AWSRegion string

	// Poll from this SQS URL
	URL string

	// Maximum number of time to attempt AWS service connection
	MaxRetries int

	// Maximum number of messages to retrieve per batch
	BatchSize int64

	// The maximum poll time (0 <= 20)
	WaitSeconds int64

	// Once a message is received by a consumer, the maximum time in seconds till others can see this
	VisibilityTimeout int64

	// Poll only once and exit
	RunOnce bool

	// Poll every X seconds defined by this value
	RunInterval int

	// Maximum number of handlers to spawn for batch processing
	MaxHandlers int

	// BusyTimeout in seconds
	BusyTimeout int

	svc          *sqs.SQS
	handlerCount int
	mutex        *sync.Mutex
	pollHandler  func(msg *sqs.Message)
}
