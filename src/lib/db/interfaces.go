package db

type IDynamoDB interface {
	RecordsReadAll() ([]RecordEntity, error)
	RecordsReadById(id string) (RecordEntity, error)
	RecordsCreate(records []Record) error
	RecordCreate(record RecordEntity) error
	RecordUpdate(record RecordEntity) error
	RecordsDelete(id string) error
}

type Record struct {
	ID				string	`json:"id,omitempty"`
	ProcessStatus	bool	`json:"processStatus"`
	Data			string	`json:"data"`
	RequestID		string	`json:"requestID"`
	RequestDate		string	`json:"requestDate,omitempty"`	// YYYYMMDD
	ProcessDate		string	`json:"processDate,omitempty"`	// YYYYMMDD
}