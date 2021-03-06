package dynamodb_test

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bradmccoydev/self-service-sdk/internal"
	"github.com/bradmccoydev/self-service-sdk/sdk/aws/dynamodb"
)

// Test CreateTable
func TestCreateTable(t *testing.T) {

	// Setup backend
	delerr := DeleteTableIfExists(TestTableNameValid)
	if delerr != nil {
		log.Fatal(delerr)
	}

	// Setup conf test data
	var emptyConf dynamodb.TableConf
	confNoMode := dynamodb.TableConf{TableName: TestTableNameValid, BillingMode: "", ReadCapacityUnits: 0, WriteCapacityUnits: 0}
	validConf := dynamodb.TableConf{TableName: TestTableNameValid, BillingMode: "PAY_PER_REQUEST", ReadCapacityUnits: 0, WriteCapacityUnits: 0}

	// Setup attribute test data
	var emptyInput []dynamodb.TableAttributes
	attribNoName := []dynamodb.TableAttributes{{Name: "", Type: "", KeyType: ""}}
	attribNoType := []dynamodb.TableAttributes{{Name: TestTableKeyFieldValid, Type: "", KeyType: ""}}
	attribKeyNoType := []dynamodb.TableAttributes{{Name: TestTableKeyFieldValid, Type: "S", KeyType: ""}}
	attribWithKey := []dynamodb.TableAttributes{{Name: TestTableKeyFieldValid, Type: "S", KeyType: "HASH"}}

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		conf      dynamodb.TableConf
		attribs   []dynamodb.TableAttributes
		expectErr bool
	}{
		{"No inputs", false, emptyConf, emptyInput, true},
		{"Just session", true, emptyConf, emptyInput, true},
		{"Session & table name", true, confNoMode, emptyInput, true},
		{"Session, table name & billing mode", true, validConf, emptyInput, true},
		{"Session, table name & attribute without a name", true, validConf, attribNoName, true},
		{"Session, table name & attribute without a type", true, validConf, attribNoType, true},
		{"Session, table name & key attribute without a type", true, validConf, attribKeyNoType, true},
		{"Session, table name & full key attribute", true, validConf, attribWithKey, false},
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
			err := dynamodb.CreateTable(sess, test.conf, test.attribs)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test DeleteTable
func TestDeleteTable(t *testing.T) {

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		expectErr bool
	}{
		{"No inputs", false, "", true},
		{"Just session", true, "", true},
		{"Session & invalid table name", true, TestTableNameInvalid, true},
		{"Session & valid table name", true, TestTableNameValid, false},
	}

	// Ensure table exists
	createerr := CreateTableIfNotExists(TestTableConf)
	if createerr != nil {
		log.Fatal(createerr)
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
			err := dynamodb.DeleteTable(sess, test.tableName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}
