package internal

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/bradmccoydev/self-service-sdk/aws/auth"
)

const (
	// EnvAwsKey - the env var for the AWS key
	EnvAwsKey string = "AWS_ACCESS_KEY_ID"
	// EnvAwsSecret - the env var for the AWS Secret
	EnvAwsSecret string = "AWS_SECRET_ACCESS_KEY"
	// EnvAwsDefRegion - the env var for the AWS default region
	EnvAwsDefRegion string = "AWS_DEFAULT_REGION"
	// EnvAwsRegion - the env var for the AWS region
	EnvAwsRegion string = "AWS_REGION"
	// EnvAwsSessionToken - the env var for the AWS session token
	EnvAwsSessionToken string = "AWS_SESSION_TOKEN"
	// TestAwsEnabled - the env var that controls whether to run AWS related tests
	TestAwsEnabled string = "TESTING_AWS_ENABLED"
	// TestValidAwsKey - the env var for ***PASSING IN*** a valid AWS key
	TestValidAwsKey string = "TESTING_AWS_ACCESS_KEY_ID"
	// TestValidAwsSecret - the env var for ***PASSING IN*** a valid AWS Secret
	TestValidAwsSecret string = "TESTING_AWS_SECRET_ACCESS_KEY"
	// TestValidAwsRegion - the env var for ***PASSING IN*** a valid AWS region
	TestValidAwsRegion string = "TESTING_AWS_DEFAULT_REGION"
	// TestValidAwsUserID - the env var for ***PASSING IN*** a valid AWS user id
	TestValidAwsUserID string = "TESTING_AWS_USER_ID"
)

// AwsCreds - Structure for handling AWS credentials
type AwsCreds struct {
	Key    string
	Secret string
	Region string
	Userid string
}

// CreateAwsSession - create an AWS session
func CreateAwsSession(valid bool) *session.Session {

	// Setup
	var creds AwsCreds
	var err error = nil

	// Clear the environment variables
	os.Unsetenv(EnvAwsKey)
	os.Unsetenv(EnvAwsSecret)
	os.Unsetenv(EnvAwsDefRegion)
	os.Unsetenv(EnvAwsRegion)

	// Are we creating a valid session?
	if valid {

		// Get the credentials
		creds, err = LoadAwsCreds()
		if err != nil {
			log.Fatal(err)
		}

		// Set the environment variables
		os.Setenv(EnvAwsKey, creds.Key)
		os.Setenv(EnvAwsSecret, creds.Secret)
		os.Setenv(EnvAwsDefRegion, creds.Region)
		os.Setenv(EnvAwsRegion, creds.Region)
	}

	// Create the session
	sess, err := auth.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	// Return the session
	return sess
}

// LoadAwsCreds - attempts to load the AWS credentials (key, secret, region)
// from the TESTING_* environment variables
func LoadAwsCreds() (AwsCreds, error) {

	// We need to grab valid AWS credentials from environment variables.
	// If any of these don't exist then we fail.
	var values AwsCreds
	var err error = nil
	key := os.Getenv(TestValidAwsKey)
	if key == "" {
		err = errors.New("The environment variable TESTING_AWS_ACCESS_KEY_ID is not set")
		return values, err
	}
	secret := os.Getenv(TestValidAwsSecret)
	if secret == "" {
		err = errors.New("The environment variable TESTING_AWS_SECRET_ACCESS_KEY is not set")
		return values, err
	}
	region := os.Getenv(TestValidAwsRegion)
	if region == "" {
		err = errors.New("The environment variable TESTING_AWS_DEFAULT_REGION is not set")
		return values, err
	}
	userid := os.Getenv(TestValidAwsUserID)
	if userid == "" {
		err = errors.New("The environment variable TESTING_AWS_USER_ID is not set")
		return values, err
	}

	// Return the values
	values.Key = key
	values.Secret = secret
	values.Region = region
	values.Userid = userid
	return values, err
}

// PerformAwsTests - checks whether the TESTING_AWS_ENABLED
// environment variable is set to TRUE
func PerformAwsTests() bool {

	// Setup
	var doAwsTests bool = false
	env := os.Getenv(TestAwsEnabled)
	if strings.ToUpper(env) == "TRUE" {
		doAwsTests = true
	}

	// Return the values
	return doAwsTests
}
