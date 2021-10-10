package api

type Record struct {
	ID				string	`json:"id,omitempty"`
	ProcessStatus	bool	`json:"processStatus"`
	Data			string	`json:"data"`
	RequestID		string	`json:"requestID"`
	RequestDate		string	`json:"requestDate,omitempty"`	// YYYYMMDD
	ProcessDate		string	`json:"processDate,omitempty"`	// YYYYMMDD
}

