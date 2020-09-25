package secretsmanager_test

import (
	"log"
	"os"
	"testing"

	"github.com/bradmccoydev/self-service-sdk/internal"
)

const (
	// The valid testing secret name
	TestSecretNameValid string = "testing"

	// An invalid testing secret name
	TestSecretNameInvalid string = "garbage"
)

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
