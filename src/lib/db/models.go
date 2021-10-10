package db

type RecordEntity struct {
	Id				string       `dynamo:"@id" json:"@id"`
	Data			string       `dynamo:"data" json:"data"`
	RequestID		string		 `dynamo:"requestID" json:"requestID"`
	RequestDate		string		 `dynamo:"requestDate,omitempty" json:"requestDate,omitempty"`	// YYYYMMDD
	ProcessDate		string		 `dynamo:"processDate,omitempty" json:"processDate,omitempty"`	// YYYYMMDD
	ProcessStatus	bool         `dynamo:"processStatus" json:"processStatus"`
}
