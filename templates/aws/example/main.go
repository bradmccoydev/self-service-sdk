package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// Request - defines the input we expect to receive
type Request struct {
	RequestID string `json:"RequestID"`
	Payload   string `json:"Payload"`
}

// Response - defines the results we will send back
type Response struct {
	Payload string `json:"Payload"`
}

// Handler - This function contains the actual microservice logic
func Handler(request Request) (Response, error) {

	return Response{
		Payload: fmt.Sprintf(request.Payload),
	}, nil
}

// Lambda entrypoint
func main() {
	lambda.Start(Handler)
}
