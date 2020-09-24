package dynamodb_test

import (
	"log"
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
var TestTableAttribs = []dynamodb.TableAttributes{
	{
		Name:    "name",
		Type:    "S",
		KeyType: dynamodb.KeyTypePartition,
	},
}

// TestTableConf represents the test table configuration
var TestTableConf = dynamodb.TableConf{
	TableName:          TestTableNameValid,
	BillingMode:        dynamodb.BillingModePayPerRequest,
	ReadCapacityUnits:  0,
	WriteCapacityUnits: 0,
}

// CreateTableIfNotExists
func CreateTableIfNotExists(tableConf dynamodb.TableConf) error {

	// If the table doesn't exist then create it
	sess := internal.CreateAwsSession(true)
	exists, _ := dynamodb.TableExists(sess, tableConf.TableName)
	var err error
	if exists == false {
		err = dynamodb.CreateTable(sess, tableConf, TestTableAttribs)
	}
	return err
}

// DeleteTableIfExists
func DeleteTableIfExists(tableName string) error {

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
		log.Printf("AWS testing variable: %s not set or set to false", internal.TestAwsEnabled)
		os.Exit(0)
	}

	// Run the various tests then exit
	exitVal := m.Run()
	os.Exit(exitVal)
}
