package db

import (
	"fmt"
	"github.com/center328/task-lambda-sqs-dynamodb/src/config"
	"sync"

	"github.com/guregu/dynamo"
)

var onceDataEngine sync.Once
var once sync.Once
var databaseGetter func() (IDynamoDB, error)
var dynamoDB *DynamoDatabase

func GetDatabase() (IDynamoDB, error) {
	onceDataEngine.Do(func() {
		switch config.Env().DbEngine {
			case "DYNAMODB":
				databaseGetter = newDynamoDatabase
			default:
				databaseGetter = func() (IDynamoDB, error) {
					return nil, fmt.Errorf("Unknown DB_ENGINE: '%s'.", config.Env().DbEngine)
				}
		}
	})
	return databaseGetter()
}

func newDynamoDatabase() (IDynamoDB, error) {
	var err error = nil
	once.Do(func() {
		dynamoDB = new(DynamoDatabase)
		dynamoDB.DB = dynamo.New(config.AwsSession, config.AwsConf)
	})
	return dynamoDB, err
}
