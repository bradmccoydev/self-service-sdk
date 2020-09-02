package dynamodb_test

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bradmccoydev/self-service-sdk/aws/dynamodb"
	"github.com/bradmccoydev/self-service-sdk/internal"
)

// Global variable for AWS credentials
var AwsCreds internal.AwsCreds

// TestMain routine for controlling setup/destruction for all tests in this package
func TestMain(m *testing.M) {

	// Set the global variable to make the values available for all tests
	var err error = nil
	AwsCreds, err = internal.LoadAwsCreds()
	if err != nil {
		log.Fatal(err)
	}

	// Run the various tests then exit
	exitVal := m.Run()
	os.Exit(exitVal)
}

// Test GetTableList
func TestGetTableList(t *testing.T) {

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		expectErr bool
	}{
		{"No session", false, true},
		{"With session", true, false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			var sess *session.Session
			if test.validSess {
				sess = internal.CreateAwsSession(true)
			} else {
				sess = internal.CreateAwsSession(false)
			}
			_, err := dynamodb.GetTableList(sess)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test NewFilterExpression
func TestNewFilterExpression(t *testing.T) {

	// Setup test data
	tests := []struct {
		desc      string
		field     string
		operator  string
		value     string
		expectErr bool
	}{
		{"No values", "", "", "", true},
		{"Just a field", "fred", "", "", true},
		{"Just an operator", "", "fred", "", true},
		{"Just a value", "", "fred", "", true},
		{"Invalid operator", "fred", "fred", "fred", true},
		{"All valid - begins with", "fred", "BW", "fred", false},
		{"All valid - contains", "fred", "CO", "fred", false},
		{"All valid - equals", "fred", "EQ", "fred", false},
		{"All valid - greater than", "fred", "GT", "fred", false},
		{"All valid - greater than or equals", "fred", "GE", "fred", false},
		{"All valid - in", "fred", "IN", "fred", false},
		{"All valid - less than", "fred", "LT", "fred", false},
		{"All valid - less than or equals", "fred", "LE", "fred", false},
		{"All valid - not equals", "fred", "NE", "fred", false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			filters := []dynamodb.Filter{{test.field, test.operator, test.value}}
			_, err := dynamodb.NewFilterExpression(filters)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}
