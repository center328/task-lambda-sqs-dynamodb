package db

import (
	"github.com/center328/task-lambda-sqs-dynamodb/src/api"
)

func RecordsMapper(records []api.Record) []RecordEntity {

	recordEntities := make([]RecordEntity, 0)
	for _, record := range records {
		newRecordEntity := RecordEntity{
			Id:				record.ID,
			RequestID: 		record.RequestID,
			RequestDate: 	record.RequestDate,
			Data: 			record.Data,
			ProcessDate: 	record.ProcessDate,
			ProcessStatus: 	record.ProcessStatus,
		}
		recordEntities = append(recordEntities, newRecordEntity)
	}
	return recordEntities
}
