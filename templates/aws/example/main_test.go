package main_test

import (
	"os"
	"testing"
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

// TestMain - This function provides an overall wrapper for testing
func TestMain(m *testing.M) {

	// Run the test cases
	exitVal := m.Run()

	// Exit with the appropriate return code
	os.Exit(exitVal)
}

// Test Handler
func TestHandler(t *testing.T) {

	// Setup request data for testing
	emptyRequest := Request{}
	emptyPayload := Request{RequestID: "12345"}

	// Setup test cases
	tests := []struct {
		desc      string
		request   Request
		expectErr bool
	}{
		{"Empty request", emptyRequest, false},
		{"Empty payload", emptyPayload, false},
	}

	// Iterate through the test cases
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			_, err := main.Handler(sess, test.tableName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}
