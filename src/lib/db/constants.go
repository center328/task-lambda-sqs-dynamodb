package db

import (
	"log"
	"os"
)

var (
	AWS_REGION      	= "ap-southeast-2"

	DB_ENGINE           = "DYNAMODB"
	DB_NAME             = "UNUSED"
	TABLE_RECORDS       = ""

	API_PREFIX_TLSD 	= "/api/v1/tlsd"
	ID_PREFIX_RECORD      = API_PREFIX_TLSD + "/records/"
)

func init() {
	awsRegion := os.Getenv("AWS_REGION")
	if len(awsRegion) > 0 {
		AWS_REGION = awsRegion
	}

	stageName := os.Getenv("STAGE_NAME")
	log.Println("Environment STAGE_NAME: ", stageName)

	TABLE_RECORDS = "tlsd-records-v1-" + stageName

	envDB_ENGINE := os.Getenv("DB_ENGINE")
	if len(envDB_ENGINE) > 0 {
		DB_ENGINE = envDB_ENGINE
		log.Println("DB_ENGINE:", DB_ENGINE)
	}

	envDB_NAME := os.Getenv("DB_NAME")
	if len(envDB_NAME) > 0 {
		DB_NAME = envDB_NAME
		log.Println("DB_NAME:", DB_NAME)
	}
}

