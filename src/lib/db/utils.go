package db

import (
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var onceDataEngine sync.Once
var once sync.Once
var databaseGetter func() (IDynamoDB, error)
var dynamoDB *DynamoDatabase

func GetDatabase() (IDynamoDB, error) {
	onceDataEngine.Do(func() {
		switch DB_ENGINE {
		case "DYNAMODB":
			databaseGetter = newDynamoDatabase
		default:
			databaseGetter = func() (IDynamoDB, error) {
				return nil, fmt.Errorf("Unknown DB_ENGINE: '%s'.", DB_ENGINE)
			}
		}
	})
	return databaseGetter()
}

func newDynamoDatabase() (IDynamoDB, error) {
	var err error = nil
	once.Do(func() {
		dynamoDB = new(DynamoDatabase)
		awsConf := &aws.Config{
			Region: aws.String(AWS_REGION),
		}
		session, errSession := session.NewSession(awsConf)
		if errSession != nil {
			log.Println("newDynamoDatabase error:", errSession)
			err = errSession
			return
		}
		dynamoDB.DB = dynamo.New(session, awsConf)
	})
	return dynamoDB, err
}
