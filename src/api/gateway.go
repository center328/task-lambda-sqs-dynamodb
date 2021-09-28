package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/google/uuid"
)

func init() {

}

type Record struct {
	ID				string	`json:"id,omitempty"`
	ProcessStatus	bool	`json:"processStatus"`
	Data			string	`json:"data"`
	RequestID		string	`json:"requestID"`
	RequestDate		string	`json:"requestDate,omitempty"`	// YYYYMMDD
	ProcessDate		string	`json:"processDate,omitempty"`	// YYYYMMDD
}

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type RequestBody struct {
	EventsCount	int		`json:"EventsCount"`
}

type ResponseBody struct {
	Id      	string	`json:"Id,omitempty"`
	Description	string	`json:"Description,omitempty"`
}

func ValidateInputs(request Request) (RequestBody, error) {
	NewBody := RequestBody{}
	ErrorMessage := ""

	if len(request.Body) == 0 {
		ErrorMessage = "No inputs provided, please provide inputs in JSON format."
		return RequestBody{}, errors.New(ErrorMessage)
	}

	// De-serialize "request.Body" which is in JSON format into "NewDevice" in Go object.
	var err = json.Unmarshal([]byte(request.Body), &NewBody)

	if err != nil {
		ErrorMessage = "Wrong format: Inputs must be a valid JSON."
		return RequestBody{}, errors.New(ErrorMessage)
	}

	if NewBody.EventsCount < 1 {
		ErrorMessage = "Invalid Data: EventsCount, Must greater than 0"
		return RequestBody{}, errors.New(ErrorMessage)
	}

	// Everything looks fine, return created NewDevice in Go struct.
	return NewBody, nil
} // End of ValidateInputs function.

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request Request) (Response, error) {
	// First, we have to validate user input.
	NewBody, err := ValidateInputs(request)
	// if inputs are not suitable, return HTTP error code 400.
	if err != nil {
		return Response{
			Body:       "" + err.Error(),
			StatusCode: 400,
		}, nil
	}

	temp := Record{
		ID: "",
		ProcessStatus: false,
		Data: request.Body,
		RequestID: request.RequestContext.RequestID,
		RequestDate: time.Now().String(),
		ProcessDate: "",
	}

	for i := 0; i < NewBody.EventsCount; i++ {
		temp.ID = uuid.NewString()
		// Serialization/Encoding "NewBody" in "item" for using in DynamoDB functions.
		record, _ := dynamodbattribute.MarshalMap(temp)

		// Till now the user have provided a valid data input.
		// Let's add it to the DynamoDB table.
		_, err = TestAws.Put(record)

		// If internal database errors occurred, return HTTP error code 500.
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:       "Internal Server Error\nDatabase error.",
				StatusCode: 500,
			}, nil
		}

		// Serialization/Encoding "NewBody" to JSON.
		jsonResponse, _ := json.Marshal(NewBody)
		return events.APIGatewayProxyResponse{
			Body: string(jsonResponse),
			// Everything looks fine, return HTTP 201
			StatusCode: 201,
		}, nil
	}

	//body, err := json.Marshal(map[string]interface{}{
	//	"message": "Okay so your other function also executed successfully!",
	//})
	//if err != nil {
	//	return Response{StatusCode: 404}, err
	//}
	//json.HTMLEscape(&buf, body)
	//
	//resp := Response{
	//	StatusCode:      200,
	//	IsBase64Encoded: false,
	//	Body:            buf.String(),
	//	Headers: map[string]string{
	//		"Content-Type":           "application/json",
	//		"X-MyCompany-Func-Reply": "world-handler",
	//	},
	//}
}

func main() {
	lambda.Start(Handler)
}
