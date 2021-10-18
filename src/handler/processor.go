package handler

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	awsSqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/center328/task-lambda-sqs-dynamodb/src/lib/db"
	sqs "github.com/center328/task-lambda-sqs-dynamodb/src/lib/stream"
	"time"
)

var logger *log.Logger
var logPrefix = "(task1-gateway) "

var queue *sqs.Config
var dynamoDB db.IDynamoDB

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler() {

	// Instantiate the queue with service connection
	queue, _ := sqs.NewSQS(sqs.Config{
		// aws config
		MaxRetries:			10,

		BatchSize:         	10,
		VisibilityTimeout: 	120,
		WaitSeconds:       	20,

		// misc config
		RunInterval: 		20,
		RunOnce:     		true,
		MaxHandlers: 		100,
		BusyTimeout: 		30,
	})

	// simulate processing a request for 2 seconds
	queue.RegisterPollHandler(func(msg *awsSqs.Message) {
		logger.Println("Wait 2 seconds for:", *msg.MessageId)
		wait := time.Duration(2) * time.Second
		<-time.After(wait)

		logger.Println("Processing:", *msg.MessageId, *msg.Body)

		record := db.Record{}
		err := json.Unmarshal([]byte(*msg.Body), &record)

		if err != nil {
			logger.Println("xerr while converting json to struct", err)
			//return Response{StatusCode: 500}, err
		}
		logger.Println("record", record)

		recordEn, err := dynamoDB.RecordsReadById(record.ID)

		if err != nil {
			logger.Println("error get in repo", err)
		}

		recordEn.ProcessDate = time.Now().String()
		recordEn.ProcessStatus = true

		err1 := dynamoDB.RecordUpdate(recordEn)

		if err1 != nil {
			logger.Println("error update in repo", err1)
		} else {

			// Simulate processing time as 10 seconds
			time.Sleep(10 * time.Second)
			log.Println("Finished:", *msg.MessageId)

			// Send back to the queue
			queue.Delete(msg)
		}
	})

	// Poll from the SQS queue
	queue.Poll()
}

func main() {
	lambda.Start(Handler)
}
