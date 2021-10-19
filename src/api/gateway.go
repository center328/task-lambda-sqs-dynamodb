package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	awsSqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/center328/task-lambda-sqs-dynamodb/src/lib/db"
	"github.com/google/uuid"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	sqs "github.com/center328/task-lambda-sqs-dynamodb/src/lib/stream"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type RequestBody struct {
	EventsCount	int		`json:"EventsCount"`
}

type ResponseBody struct {
	Id      	string	`json:"Id,omitempty"`
	Description	string	`json:"Description,omitempty"`
}

var logger *log.Logger
var logPrefix = "(task1-gateway) "

var queue *sqs.Config
var dynamoDB db.IDynamoDB

func init() {

	logger = log.New(os.Stdout, logPrefix, log.LstdFlags|log.Lshortfile)

	// Instantiate the queue with service connection
	queue1, err1 := sqs.NewSQS(sqs.Config{
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

	if err1 == nil {
		logger.Println(err1)
	} else {
		queue = queue1
	}

	db1, err2 := db.GetDatabase()

	if err2 == nil {
		logger.Println(err2)
	} else {
		dynamoDB = db1
	}

}

func ValidateInputs(request Request) (RequestBody, error) {
	NewBody := RequestBody{}
	ErrorMessage := ""

	if len(request.Body) == 0 {
		ErrorMessage = "No inputs provided, please provide inputs in JSON format."
		return RequestBody{}, errors.New(ErrorMessage)
	}

	// De-serialize "request.Body" which is in JSON format into "NewDevice" in Go object.
	var err = json.Unmarshal([]byte(request.Body), &NewBody)

	if err != nil {
		ErrorMessage = "Wrong format: Inputs must be a valid JSON."
		return RequestBody{}, errors.New(ErrorMessage)
	}

	if NewBody.EventsCount < 1 {
		ErrorMessage = "Invalid Data: EventsCount, Must greater than 0"
		return RequestBody{}, errors.New(ErrorMessage)
	}

	// Everything looks fine, return created NewDevice in Go struct.
	return NewBody, nil
} // End of ValidateInputs function.

func createRequests(size int, body string, requestID string) []db.Record {
	var records []db.Record

	var temp = db.Record{
		ID: "",
		ProcessStatus: false,
		Data: body,
		RequestID: requestID,
		RequestDate: time.Now().String(),
		ProcessDate: "",
	}

	for i := 0; i < size; i++ {
		temp.ID = uuid.NewString()
		records = append(records, temp)
	}

	return records
}

func createMessagesToEnqueue(msgs []db.Record)  []*awsSqs.SendMessageBatchRequestEntry {
	var msgBatch []*awsSqs.SendMessageBatchRequestEntry
	for i := 0; i < len(msgs); i++ {
		data, _ := json.Marshal(msgs[i])
		message := &awsSqs.SendMessageBatchRequestEntry{
			Id:                     aws.String(`uniqueID_` + strconv.Itoa(i)),
			MessageBody:            aws.String(string(data)),
			MessageDeduplicationId: aws.String(`dupID_` + strconv.Itoa(i)),
			MessageGroupId:         aws.String("task1Queue"),
		}
		msgBatch = append(msgBatch, message)
	}

	return msgBatch
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	// First, we have to validate user input.
	NewBody, err := ValidateInputs(request)
	// if inputs are not suitable, return HTTP error code 400.
	if err != nil {
		return Response{
			Body:       "" + err.Error(),
			StatusCode: 400,
		}, nil
	}

	reqs := createRequests(NewBody.EventsCount, request.Body, request.RequestContext.RequestID)
	errSQS := queue.Enqueue(createMessagesToEnqueue(reqs))

	if errSQS == nil {
		errDB := dynamoDB.RecordsCreate(reqs)
		if errDB != nil {
			logger.Println(errDB)
			return Response{StatusCode: 500}, errDB
		} else {
			var buf bytes.Buffer

			body, errMar := json.Marshal(reqs)
			if errMar != nil {
				return Response{StatusCode: 500}, errMar
			}
			json.HTMLEscape(&buf, body)

			resp := Response{
				StatusCode:      200,
				IsBase64Encoded: false,
				Body:            buf.String(),
				Headers: map[string]string{
					"Content-Type":           "application/json",
					"X-MyCompany-Func-Reply": "hello-handler",
				},
			}

			return resp, nil
		}
	} else {
		logger.Println("Queue Insert Failed...")
		return Response{StatusCode: 500}, errSQS
	}
}

func main() {
	lambda.Start(Handler)
}
