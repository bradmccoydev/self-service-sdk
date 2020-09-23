package dynamodb_test

import (
	"errors"
	"os"
	"testing"

	"github.com/bradmccoydev/self-service-sdk/aws/dynamodb"
	"github.com/bradmccoydev/self-service-sdk/internal"
)

const (
	// The valid testing table name
	TestTableNameValid string = "testing"

	// An invalid testing table name
	TestTableNameInvalid string = "garbage"

	// The valid testing table key field
	TestTableKeyFieldValid string = "name"

	// An invalid testing table name
	TestTableKeyFieldInvalid string = "garbage"
)

// TestTableItem represents an item from the service table
type TestTableItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TestTableAttribs represents the test table attributes
var TestTableKeys = []dynamodb.TableAttributes{
	{
		Name:    "name",
		Type:    "S",
		IsKey:   true,
		KeyType: "HASH",
	},
}

// DeleteTableIfExists
func DeleteTableIfExists(tableName string) error {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}

	// If the table exists then delete it
	sess := internal.CreateAwsSession(true)
	exists, _ := dynamodb.TableExists(sess, TestTableNameValid)
	var err error
	if exists {
		err = dynamodb.DeleteTable(sess, tableName)
	}
	return err
}

// TestMain routine for controlling setup/destruction for all tests in this package
func TestMain(m *testing.M) {

	// Do we need to do these tests?
	var doTests bool = internal.PerformAwsTests()
	if doTests == false {
		os.Exit(0)
	}

	// Run the various tests then exit
	exitVal := m.Run()
	os.Exit(exitVal)
}
