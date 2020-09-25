package dynamodb_test

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bradmccoydev/self-service-sdk/internal"
	"github.com/bradmccoydev/self-service-sdk/sdk/aws/dynamodb"
)

// Test DescribeTable
func TestDescribeTable(t *testing.T) {

	// Setup backend
	createerr := CreateTableIfNotExists(TestTableConf)
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		expectErr bool
	}{
		{"No session", false, "", true},
		{"With session but no table name", true, "", true},
		{"With session and invalid table name", true, TestTableNameInvalid, true},
		{"With session and valid table name", true, TestTableNameValid, false},
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
			_, err := dynamodb.DescribeTable(sess, test.tableName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test GetTableArn
func TestGetTableArn(t *testing.T) {

	// Setup backend
	createerr := CreateTableIfNotExists(TestTableConf)
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup test data
	tests := []struct {
		desc      string
		tableName string
		expectErr bool
	}{
		{"Invalid table name", TestTableNameInvalid, true},
		{"Valid table name", TestTableNameValid, false},
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

	// Setup backend
	createerr := CreateTableIfNotExists(TestTableConf)
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup test data
	tests := []struct {
		desc      string
		tableName string
		expectErr bool
	}{
		{"Invalid table name", TestTableNameInvalid, true},
		{"Valid table name", TestTableNameValid, false},
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

// Test TableExists
func TestTableExists(t *testing.T) {

	// Setup backend
	createerr := CreateTableIfNotExists(TestTableConf)
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup test data
	tests := []struct {
		desc        string
		validSess   bool
		tableName   string
		expectFound bool
		expectErr   bool
	}{
		{"No session", false, "", false, true},
		{"With session but no table name", true, "", false, true},
		{"With session and invalid table name", true, TestTableNameInvalid, false, false},
		{"With session and valid table name", true, TestTableNameValid, true, false},
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
			actual, err := dynamodb.TableExists(sess, test.tableName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
				internal.Equals(t, test.expectFound, actual)
			}
		})
	}
}
