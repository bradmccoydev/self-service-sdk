package dynamodb_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/bradmccoydev/self-service-sdk/aws/dynamodb"
	"github.com/bradmccoydev/self-service-sdk/internal"
)

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
		{"No inputs", noCond, noFilter, noProj, false},
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

// Test ScanTable
func TestScanTable(t *testing.T) {

	// Setup response structures
	type ServiceItem struct {
		Service       string `json:"service"`
		Title         string `json:"title"`
		Description   string `json:"description"`
		Documentation string `json:"documentation"`
		Type          string `json:"type"`
	}
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
			err := dynamodb.ScanTable(sess, test.tableName, test.expr, &response)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}
