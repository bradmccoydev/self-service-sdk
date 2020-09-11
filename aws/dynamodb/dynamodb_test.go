package dynamodb_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/bradmccoydev/self-service-sdk/aws/dynamodb"
	"github.com/bradmccoydev/self-service-sdk/internal"
)

// ServiceItem represents an item from the service table
type ServiceItem struct {
	Service       string `json:"service"`
	Version       string `json:"version"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Documentation string `json:"documentation"`
	Type          string `json:"type"`
}

// Global variable for AWS credentials
var AwsCreds internal.AwsCreds

// TestMain routine for controlling setup/destruction for all tests in this package
func TestMain(m *testing.M) {

	// Do we need to do these tests?
	var doTests bool = internal.PerformAwsTests()
	if doTests == false {
		os.Exit(0)
	}

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

// Test GetTableArn
func TestGetTableArn(t *testing.T) {

	// Setup test data
	tests := []struct {
		desc      string
		tableName string
		expectErr bool
	}{
		{"Invalid table name", "FredSmith", true},
		{"Valid table name", "service", false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			sess := internal.CreateAwsSession(true)
			_, err := dynamodb.GetTableArn(sess, test.tableName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test GetTableItemCount
func TestGetTableItemCount(t *testing.T) {

	// Setup test data
	tests := []struct {
		desc      string
		tableName string
		expectErr bool
	}{
		{"Invalid table name", "FredSmith", true},
		{"Valid table name", "service", false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			sess := internal.CreateAwsSession(true)
			_, err := dynamodb.GetTableItemCount(sess, test.tableName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test GetTableDetails
func TestGetTableDetails(t *testing.T) {

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		expectErr bool
	}{
		{"No session", false, "", true},
		{"With session but no table name", true, "", true},
		{"With session and invalid table name", true, "FredSmith", true},
		{"With session and valid table name", true, "service", false},
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
			_, err := dynamodb.GetTableDetails(sess, test.tableName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
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
