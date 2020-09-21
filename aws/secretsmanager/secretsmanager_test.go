package secretsmanager_test

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bradmccoydev/self-service-sdk/aws/secretsmanager"
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

// Test GetSecretKeyValue
func TestGetSecretKeyValue(t *testing.T) {

	// Setup test data
	tests := []struct {
		desc       string
		validSess  bool
		secretName string
		expectErr  bool
	}{
		{"No session", false, "", true},
		{"With session but no secret name", true, "", true},
		{"With session and invalid secret name", true, "Fred", false},
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
			_, err := secretsmanager.GetSecretKeyValue(sess, test.secretName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test GetSecretString
func TestGetSecretString(t *testing.T) {

	// Setup test data
	tests := []struct {
		desc       string
		validSess  bool
		secretName string
		expectErr  bool
	}{
		{"No session", false, "", true},
		{"With session but no secret name", true, "", true},
		{"With session and invalid secret name", true, "Fred", false},
		{"With session and valid secret name", true, "service", false},
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
			_, err := secretsmanager.GetSecretString(sess, test.secretName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}
