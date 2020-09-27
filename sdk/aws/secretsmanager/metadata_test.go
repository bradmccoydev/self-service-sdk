package secretsmanager_test

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bradmccoydev/self-service-sdk/internal"
	"github.com/bradmccoydev/self-service-sdk/sdk/aws/secretsmanager"
)

// Test DescribeSecret
func TestDescribeSecret(t *testing.T) {

	// Setup backend
	createerr := CreateTestSecretIfNotExists()
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup test data
	tests := []struct {
		desc       string
		validSess  bool
		secretName string
		expectErr  bool
	}{
		{"No session", false, "", true},
		{"With session but no secret name", true, "", true},
		{"With session and invalid secret name", true, TestSecretNameInvalid, true},
		{"With session and valid secret name", true, TestSecretNameValid, false},
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
			_, err := secretsmanager.DescribeSecret(sess, test.secretName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}
