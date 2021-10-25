package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/center328/task-lambda-sqs-dynamodb/src/config"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/center328/task-lambda-sqs-dynamodb/src/lib/db"
	sqs "github.com/center328/task-lambda-sqs-dynamodb/src/lib/stream"
	"time"
)

var logger *log.Logger
var logPrefix = "(task1-handler) "

var queue *sqs.Config
var dynamoDB db.IDynamoDB

func init() {
	env := config.Env()

	logger = log.New(os.Stdout, logPrefix, log.LstdFlags|log.Lshortfile)

	// Instantiate the queue with service connection
	queue1, err1 := sqs.NewSQS(sqs.Config{
		// aws config
		AWSRegion:  		env.AWSRegion,
		URL:               	env.SQSURL,
		BatchSize:         	env.SQSBatchSize,
		VisibilityTimeout: 	120,
		WaitSeconds:       	20,

		// misc config
		RunInterval: 		20,
		RunOnce:     		env.RunOnce,
		MaxHandlers: 		100,
		BusyTimeout: 		30,
	})

	if err1 != nil {
		logger.Println(err1)
	} else {
		queue = queue1
	}

	db1, err2 := db.GetDatabase()

	if err2 != nil {
		logger.Println(err2)
	} else {
		dynamoDB = db1
	}

}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(sqsEvent events.SQSEvent) {
	if len(sqsEvent.Records) == 0 {
		logger.Println("error: No SQS message passed to function")
		return
	}
	// simulate processing a request for 2 seconds
	for _, msg := range sqsEvent.Records {
		logger.Println("Wait 2 seconds for:", msg.MessageId)
		wait := time.Duration(2) * time.Second
		<-time.After(wait)

		logger.Println("Processing:", msg.MessageId, msg.Body)

		record := db.Record{}
		err := json.Unmarshal([]byte(msg.Body), &record)

		if err != nil {
			logger.Println("xerr while converting json to struct", err)
			//return Response{StatusCode: 500}, err
		}
		logger.Println("record    ", record)

		recordEn, err := dynamoDB.RecordsReadById(record.ID)

		if err != nil {
			logger.Println("error get in repo", err)
		}

		recordEn.ProcessDate = time.Now().String()
		recordEn.ProcessStatus = true
		logger.Println("record    ", recordEn)

		err1 := dynamoDB.RecordUpdate(recordEn)

		if err1 != nil {
			logger.Println("error update in repo", err1)
		} else {
			log.Println("Record Updated in repo")
			//
			//// Simulate processing time as 10 seconds
			//time.Sleep(10 * time.Second)
			log.Println("Finished:", msg.MessageId)

			// Send back to the queue
			queue.Delete(msg)
		}
	}

	// Poll from the SQS queue
	queue.Poll()
}

func main() {
	lambda.Start(Handler)
}
