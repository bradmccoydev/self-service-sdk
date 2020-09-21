package dynamodb_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/bradmccoydev/self-service-sdk/aws/dynamodb"
	"github.com/bradmccoydev/self-service-sdk/internal"
)

// Test CreateItem
func TestCreateItem(t *testing.T) {

	// Setup input test data
	timeStamp := time.Now().Format(time.RFC3339)
	emptyInput := ServiceItem{}
	noKey := ServiceItem{Title: "Fred"}
	emptyKey := ServiceItem{Service: "", Title: "Fred"}
	validInput := ServiceItem{Service: "Fred", Version: timeStamp, Title: "Fred"}

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		input     ServiceItem
		expectErr bool
	}{
		{"No inputs", false, "", emptyInput, true},
		{"Just session", true, "", emptyInput, true},
		{"Session & invalid table name", true, "fred", emptyInput, true},
		{"Session & valid table name", true, "service", emptyInput, true},
		{"Session, valid table name & no key", true, "service", noKey, true},
		{"Session, valid table name & empty key", true, "service", emptyKey, true},
		{"Valid input", true, "service", validInput, false},
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
			err := dynamodb.CreateItem(sess, test.tableName, test.input)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test QueryItems
func TestQueryItems(t *testing.T) {

	// Setup response array
	var response *[]ServiceItem

	// Setup key condition test data
	var emptyKey []dynamodb.Condition
	invalidKeyField := []dynamodb.Condition{{Field: "fred", Operator: "EQ", Value: "123"}}
	invalidKeyVal := []dynamodb.Condition{{Field: "service", Operator: "EQ", Value: "123"}}
	validKey := []dynamodb.Condition{{Field: "service", Operator: "EQ", Value: "123"}}

	// Setup expression test data
	var emptyExpression expression.Expression
	emptyKeyExpr, _ := dynamodb.NewExpression(emptyKey, nil, nil)
	invalidKeyFieldExpr, _ := dynamodb.NewExpression(invalidKeyField, nil, nil)
	invalidKeyValueExpr, _ := dynamodb.NewExpression(invalidKeyVal, nil, nil)
	validKeyExpr, _ := dynamodb.NewExpression(validKey, nil, nil)

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		expr      expression.Expression
		expectErr bool
	}{
		{"No inputs", false, "", emptyExpression, true},
		{"Just session", true, "", emptyExpression, true},
		{"Session & invalid table name", true, "fred", emptyExpression, true},
		{"Session & valid table name", true, "service", emptyKeyExpr, true},
		{"Session, valid table name & empty key", true, "service", emptyExpression, true},
		{"Session, valid table name & invalid key field", true, "service", invalidKeyFieldExpr, true},
		{"Session, valid table name & invalid key value", true, "service", invalidKeyValueExpr, false},
		{"Session, valid table name & valid key", true, "service", validKeyExpr, false},
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
			err := dynamodb.QueryItems(sess, test.tableName, test.expr, &response)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test ScanItems
func TestScanItems(t *testing.T) {

	// Setup response array
	var response *[]ServiceItem

	// Setup filter test data
	invalidFilter := []dynamodb.Condition{{Field: "fred", Operator: "EQ", Value: "123"}}
	validFilter := []dynamodb.Condition{{Field: "service", Operator: "EQ", Value: "123"}}

	// Setup expression test data
	var emptyExpression expression.Expression
	invalidFilterExpr, _ := dynamodb.NewExpression(nil, invalidFilter, nil)
	validFilterExpr, _ := dynamodb.NewExpression(nil, validFilter, nil)

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		expr      expression.Expression
		expectErr bool
	}{
		{"No inputs", false, "", emptyExpression, true},
		{"Just session", true, "", emptyExpression, true},
		{"Session & invalid table name", true, "fred", emptyExpression, true},
		{"Session & valid table name", true, "service", emptyExpression, false},
		{"Session, valid table name & invalid filter", true, "service", invalidFilterExpr, false},
		{"Session, valid table name & valid filter", true, "service", validFilterExpr, false},
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
			err := dynamodb.ScanItems(sess, test.tableName, test.expr, &response)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}
