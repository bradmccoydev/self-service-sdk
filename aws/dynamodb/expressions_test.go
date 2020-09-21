package dynamodb_test

import (
	"fmt"
	"testing"

	"github.com/bradmccoydev/self-service-sdk/aws/dynamodb"
	"github.com/bradmccoydev/self-service-sdk/internal"
)

// Test NewExpression
func TestNewExpression(t *testing.T) {

	// Setup key condition test data
	var noCond []dynamodb.Condition
	emptyCond := []dynamodb.Condition{{}}
	validCond := []dynamodb.Condition{{Field: "service", Operator: "EQ", Value: "123"}}

	// Setup filter test data
	var noFilter []dynamodb.Condition
	emptyFilter := []dynamodb.Condition{{}}
	validFilter := []dynamodb.Condition{{Field: "service", Operator: "EQ", Value: "123"}}

	// Setup projection test data
	var noProj []dynamodb.Field
	emptyProj := []dynamodb.Field{{}}
	noNameProj := []dynamodb.Field{{Name: ""}}
	singleProj := []dynamodb.Field{{Name: "fred"}}
	multipleProj := []dynamodb.Field{{Name: "fred"}, {Name: "harry"}, {Name: "norm"}}

	// Setup test data
	tests := []struct {
		desc        string
		keys        []dynamodb.Condition
		filters     []dynamodb.Condition
		projections []dynamodb.Field
		expectErr   bool
	}{
		{"No inputs", noCond, noFilter, noProj, true},
		{"Empty inputs", emptyCond, emptyFilter, emptyProj, true},
		{"Valid key condition", validCond, noFilter, noProj, false},
		{"Valid key condition & empty filter", validCond, emptyFilter, emptyProj, true},
		{"Valid key condition & filter & invalid projection", validCond, validFilter, noNameProj, true},
		{"Valid key condition & filter & single projection", validCond, validFilter, singleProj, false},
		{"Valid key condition & filter & multi projection", validCond, validFilter, multipleProj, false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			expr, err := dynamodb.NewExpression(test.keys, test.filters, test.projections)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
				fmt.Println("Expression", expr)
			}
		})
	}
}
