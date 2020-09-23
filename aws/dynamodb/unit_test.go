package dynamodb_test

import (
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
