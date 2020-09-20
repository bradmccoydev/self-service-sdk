package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bradmccoydev/self-service-sdk/configutil"
	"github.com/bradmccoydev/self-service-sdk/logutil"
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

// Handler - This function contains the actual service logic
func Handler(request Request) (Response, error) {

	// Setup
	var resp Response

	// Create array of the environment variables
	// we need to get from the Lambda for configuration
	vars := []configutil.EnvVariable{
		{EnvVar: "Environment", ConfigKey: "env"},
		{EnvVar: "LogLevel", ConfigKey: "log.level"},
		{EnvVar: "Something", ConfigKey: "somekey"}}

	// Setup Config handler
	config, err := configutil.NewConfigFromEnv(nil, vars)
	if err != nil {
		return resp, err
	}

	// Setup log handler
	logLevel := config.GetString("log.level")
	logutil.New("", logLevel, false)

	return Response{
		Payload: fmt.Sprintf(request.Payload),
	}, nil
}

// Lambda entrypoint
func main() {
	lambda.Start(Handler)
}
