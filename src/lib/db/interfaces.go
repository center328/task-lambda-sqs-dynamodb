package db

import "github.com/center328/task-lambda-sqs-dynamodb/src/api"

type IDynamoDB interface {
	RecordsReadAll() ([]RecordEntity, error)
	RecordsReadById(id string) (RecordEntity, error)
	RecordsCreate(records []api.Record) error
	RecordCreate(record RecordEntity) error
	RecordUpdate(record RecordEntity) error
	RecordsDelete(id string) error
}
