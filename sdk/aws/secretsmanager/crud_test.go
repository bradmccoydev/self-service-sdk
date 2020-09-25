package secretsmanager_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bradmccoydev/self-service-sdk/internal"
	"github.com/bradmccoydev/self-service-sdk/sdk/aws/secretsmanager"
)

// Test CreateSecretKeyValue
func TestCreateSecretKeyValue(t *testing.T) {

	// Setup key/values
	type kv map[string]string
	var emptyKv kv
	noValues := make(kv)
	validKv := make(kv)
	validKv["fred"] = "nerk"

	// Setup test data
	tests := []struct {
		desc       string
		validSess  bool
		secretName string
		secretVal  kv
		expectErr  bool
	}{
		{"No values", false, "", emptyKv, true},
		{"With session but no secret name", true, "", emptyKv, true},
		{"With session and secret name", true, TestSecretNameValid, emptyKv, true},
		{"With session, secret name & empty kv", true, TestSecretNameValid, noValues, true},
		{"With session, secret name & secret value", true, TestSecretNameValid, validKv, false},
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
			err := secretsmanager.CreateSecretKeyValue(sess, test.secretName, test.desc, test.secretVal)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test CreateSecretString
func TestCreateSecretString(t *testing.T) {

	// Setup test data
	tests := []struct {
		desc       string
		validSess  bool
		secretName string
		secretVal  string
		expectErr  bool
	}{
		{"No values", false, "", "", true},
		{"With session but no secret name", true, "", "", true},
		{"With session and secret name", true, TestSecretNameValid, "", true},
		{"With session, secret name & secret value", true, TestSecretNameValid, "somesupersecretstring", false},
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
			err := secretsmanager.CreateSecretString(sess, test.secretName, test.desc, test.secretVal)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test DeleteSecret
func TestDeleteSecret(t *testing.T) {

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
			err := secretsmanager.DeleteSecret(sess, test.secretName)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
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
