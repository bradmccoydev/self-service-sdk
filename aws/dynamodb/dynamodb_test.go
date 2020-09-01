package dynamodb_test

import (
	"testing"

	"github.com/bradmccoydev/self-service-sdk/aws/dynamodb"
	"github.com/bradmccoydev/self-service-sdk/internal"
)

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
