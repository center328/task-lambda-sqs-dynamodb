package api

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	awsSqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/center328/task-lambda-sqs-dynamodb/src/lib/db"
	"github.com/google/uuid"
	"strconv"
	"time"
)

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
