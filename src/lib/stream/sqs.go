package stream

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var logger *log.Logger
var logPrefix = "(task1-sqs) "

func init() {
	logger = log.New(os.Stdout, logPrefix, log.LstdFlags|log.Lshortfile)
}

// NewSQS Instantiate a SQS instance
func NewSQS(opts Config) (*Config, error) {
	// Validate parameters
	validateErr := validateOpts(opts)
	if validateErr != nil {
		logger.Println(validateErr)
		return nil, validateErr
	}

	creds := credentials.NewEnvCredentials()
	if _, err := creds.Get(); err != nil {
		logger.Println("AWS Credential error", err)
		return nil, errors.New("Invalid AWS credentials. Please make sure that `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` is present in the env")
	}

	// Create AWS Config
	awsConfig := aws.NewConfig().WithRegion(AWSRegion).WithMaxRetries(opts.MaxRetries).WithCredentials(creds)
	if awsConfig == nil {
		logger.Println("Invalid AWS Config")
		return nil, errors.New("Something is wrong with your AWS config parameters")
	}

	// Establish a session
	newSession := session.Must(session.NewSession(awsConfig))
	if newSession == nil {
		logger.Println("Unable to create session")
		return nil, errors.New("Unable to create session")
	}

	// Create a service connection
	svc := sqs.New(newSession)
	if svc == nil {
		logger.Println("Unable to connect to SQS")
		return nil, errors.New("Unable to create a service connection with AWS SQS")
	}

	logger.Println("Fetching queue attributes")
	if _, err := svc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		QueueUrl: &URL,
	}); err != nil {
		logger.Println("Unable to fetch queue attributes", err)
		return nil, errors.New("Unable to get queue attributes")
	}
	logger.Println("Connected to Queue")

	opts.svc = svc
	opts.mutex = &sync.Mutex{}
	return &opts, nil
}

// Poll for messages in the queue
func (s *Config) Poll() {
	if s.svc == nil {
		logger.Fatalln("No service connection")
	}

	wg := sync.WaitGroup{}
	batch := 0

	for {
		batch++
		childLogger := log.New(os.Stdout, fmt.Sprintf("%sbatch-%d ", logPrefix, batch), log.LstdFlags|log.Lshortfile)

		childLogger.Println("Start receiving messages")

		maxMsgs := s.BatchSize

		// Is running at capacity?
		if s.MaxHandlers > 0 {
			for s.handlerCount >= s.MaxHandlers {
				childLogger.Printf("Reached max handler count")
				childLogger.Printf("Going to wait state for %d seconds", s.BusyTimeout)
				<-time.After(time.Duration(s.BusyTimeout) * time.Second)
			}
			availableHandlers := int64(s.MaxHandlers - s.handlerCount)
			if availableHandlers < maxMsgs {
				maxMsgs = availableHandlers
			}
		}

		childLogger.Printf("Polling for a maximum of %d messages", maxMsgs)

		result, err := s.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            &URL,
			MaxNumberOfMessages: &maxMsgs,
			WaitTimeSeconds:     &s.WaitSeconds,
			VisibilityTimeout:   &s.VisibilityTimeout,
		})

		// Retrieve error?
		if err != nil {
			childLogger.Println("ReceiveMessageError:", err)
			break
		}

		// Message log
		if len(result.Messages) == 0 {
			childLogger.Println("Queue is empty")
		} else {
			childLogger.Printf("Fetched %d messages", len(result.Messages))
		}

		// Process messages
		for _, msg := range result.Messages {
			if s.pollHandler == nil {
				childLogger.Println("No Poll handler registered. Register a handler for custom handling")
			} else {
				s.handlerCount++
				wg.Add(1)

				go func(m *sqs.Message) {
					s.pollHandler(m)

					s.mutex.Lock()
					s.handlerCount--
					s.mutex.Unlock()

					wg.Done()
				}(&(*msg))
			}

			childLogger.Printf("Spawned handler for %s", *msg.MessageId)
		}

		if s.RunOnce == true {
			childLogger.Println(`Exiting since configured to run once`)
			break
		} else {
			childLogger.Printf("Waiting for %d seconds before polling for next batch", s.RunInterval)
			<-time.After(time.Duration(s.RunInterval) * time.Second)
		}

		childLogger.Println("Finished polling")
	}

	wg.Wait()
}

// Enqueue messages to SQS
func (s *Config) Enqueue(msgBatch []*sqs.SendMessageBatchRequestEntry) error {
	if s.svc == nil {
		logger.Fatal("No service connection")
	}

	logger.Printf(`Enqueuing %d messages`, len(msgBatch))

	result, err := s.svc.SendMessageBatch(&sqs.SendMessageBatchInput{
		QueueUrl: &URL,
		Entries:  msgBatch,
	})

	logger.Printf("Enqueue result: %d success, %d failed", len(result.Successful), len(result.Failed))
	return err
}

// Delete a SQS message from the queue
func (s *Config) Delete(msg *sqs.Message) error {

	logger.Printf("Delete message with ID %s", *msg.MessageId)
	_, err := s.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &URL,
		ReceiptHandle: msg.ReceiptHandle,
	})

	return err
}

// RegisterPollHandler : A method to register a custom Poll Handling method
func (s *Config) RegisterPollHandler(pollHandler func(msg *sqs.Message)) {
	s.pollHandler = pollHandler
}

// ChangeVisibilityTimeout : Method to change visibility timeout of a message.
func (s *Config) ChangeVisibilityTimeout(msg *sqs.Message, seconds int64) bool {
	retVal := false
	logger.Printf("change visibility timeout for message ID %s", *msg.MessageId)

	if s.svc == nil {
		logger.Fatal("SQS Connection failed")
		return retVal
	}

	strURL := &URL
	receiptHandle := *msg.ReceiptHandle

	changeMessageVisibilityInput := sqs.ChangeMessageVisibilityInput{}

	changeMessageVisibilityInput.SetQueueUrl(*strURL)
	changeMessageVisibilityInput.SetReceiptHandle(receiptHandle)
	changeMessageVisibilityInput.SetVisibilityTimeout(seconds)

	_, err := s.svc.ChangeMessageVisibility(&changeMessageVisibilityInput)

	if err == nil {
		logger.Printf("changed visibility timeout success for %s", *msg.MessageId)
		retVal = true
	} else {
		logger.Printf("change visibility timeout failed for %s", *msg.MessageId)
	}

	return retVal
}
