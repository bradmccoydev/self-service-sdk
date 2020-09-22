package dynamodb_test

import (
	"log"
	"os"
	"testing"

	"github.com/bradmccoydev/self-service-sdk/internal"
)

// TestTableItem represents an item from the service table
type TestTableItem struct {
	Service       string `json:"service"`
	Version       string `json:"version"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Documentation string `json:"documentation"`
	Type          string `json:"type"`
}

// TestTableKeys represents the key fields for the test table
type TestTableKeys struct {
	Service string `json:"service"`
	Version string `json:"version"`
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
