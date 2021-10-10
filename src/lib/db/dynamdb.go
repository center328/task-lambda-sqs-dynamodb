package db

import (
	"log"

	"github.com/guregu/dynamo"
	"task-lambda-sqs-dynamodb/src/api"
)

type DynamoDatabase struct {
	IDynamoDB
	DB *dynamo.DB
}

func (db *DynamoDatabase) RecordsReadAll() ([]RecordEntity, error) {
	results := make([]RecordEntity, 0)
	unique := make(map[string]RecordEntity)
	table := db.DB.Table(TABLE_RECORDS)
	itr := table.Scan().Iter()
	var recordEntity RecordEntity
	for itr.Next(&recordEntity) {
		unique[recordEntity.Id] = recordEntity
	}
	if itr.Err() != nil {
		log.Print("RecordsReadAll iterator error:", itr.Err())
		return results, itr.Err()
	}
	for _, recordEntity := range unique {
		results = append(results, recordEntity)
	}
	return results, nil
}

func (db *DynamoDatabase) RecordsReadById(id string) (RecordEntity, error) {
	var lastError error
	table := db.DB.Table(TABLE_RECORDS)

	var result RecordEntity
	err := table.Get("@id", id).One(&result)
	if err != nil {
		log.Println("RecordGetById error:", err)
		lastError = err
	}
	return result, lastError
}

func (db *DynamoDatabase) RecordsCreate(records []api.Record) error {
	recordEntities := RecordsMapper(records)
	var lastError error
	table := db.DB.Table(TABLE_RECORDS)
	for _, recordEntity := range recordEntities {
		err := table.Put(recordEntity).Run()
		if err != nil {
			log.Println("RecordsCreate error:", err)
			lastError = err
		}
	}
	return lastError
}

func (db *DynamoDatabase) RecordCreate(record RecordEntity) error {
	var lastError error
	table := db.DB.Table(TABLE_RECORDS)
	err := table.Put(record).Run()
	if err != nil {
		log.Println("RecordsCreate error:", err)
		lastError = err
	}
	return lastError
}

func (db *DynamoDatabase) RecordUpdate(record RecordEntity) error {
	var lastError error
	table := db.DB.Table(TABLE_RECORDS)
	err := table.Put(record).If("$ = ?", "@id", record.Id).Run()
	if err != nil {
		log.Println("RecordUpdate error:", err)
		lastError = err
	}
	return lastError
}

func (db *DynamoDatabase) RecordsDelete(id string) error {
	var lastError error
	table := db.DB.Table(TABLE_RECORDS)
	err := table.Delete("@id", id).If("$ = ?", "@id", id).Run()
	if err != nil {
		log.Println("RecordDelete error:", err)
		lastError = err
	}
	return lastError
}
