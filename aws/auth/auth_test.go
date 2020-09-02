package auth_test

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/bradmccoydev/self-service-sdk/aws/auth"
	"github.com/bradmccoydev/self-service-sdk/internal"
)

const (
	// Environment variable for AWS key
	envAwsKey string = "AWS_ACCESS_KEY_ID"
	// Environment variable for AWS Secret
	envAwsSecret string = "AWS_SECRET_ACCESS_KEY"
	// Environment variable for AWS region
	envAwsRegion string = "AWS_DEFAULT_REGION"
	// Environment variable for AWS session token
	envAwsSessionToken string = "AWS_SESSION_TOKEN"
	// Environment variable for ***PASSING IN*** a valid AWS key
	testValidAwsKey string = "TESTING_AWS_ACCESS_KEY_ID"
	// Environment variable for ***PASSING IN*** a valid AWS Secret
	testValidAwsSecret string = "TESTING_AWS_SECRET_ACCESS_KEY"
	// Environment variable for ***PASSING IN*** a valid AWS region
	testValidAwsRegion string = "TESTING_AWS_DEFAULT_REGION"
	// Environment variable for ***PASSING IN*** a valid AWS user id
	testValidAwsUserId string = "TESTING_AWS_USER_ID"
)

// Common structure for defining test data
type testdef struct {
	desc        string
	key         bool
	keyval      string
	secret      bool
	secval      string
	region      bool
	regval      string
	config      bool
	confval     string
	option      bool
	optnval     string
	expectErr   bool
	expectedVal string
}

// Common structure for managing credentials
type awscreds struct {
	key    string
	secret string
	region string
	userid string
}

// Global variable for AWS credentials
var AwsCreds awscreds

// createEnvVar - create an environment variable
func createEnvVar(t *testing.T, testcase string, varname string, varvalue string) {
	//t.Log(fmt.Printf("Test case: %s Creating environment variable: %s", testcase, varname))
	os.Setenv(varname, varvalue)
}

// deleteEnvVar - delete an environment variable
func deleteEnvVar(t *testing.T, testcase string, varname string) {
	//t.Log(fmt.Printf("Test case: %s Deleting environment variable: %s", testcase, varname))
	os.Unsetenv(varname)
}

// manageTestEnvVar handles creation/desctruction env vars
func manageTestEnvVar(t *testing.T, testdata testdef) func(t *testing.T) {

	// AWS Key
	if testdata.key {
		createEnvVar(t, testdata.desc, envAwsKey, testdata.keyval)
	}
	// AWS Secret
	if testdata.secret {
		createEnvVar(t, testdata.desc, envAwsSecret, testdata.secval)
	}
	// AWS Region
	if testdata.region {
		createEnvVar(t, testdata.desc, envAwsRegion, testdata.regval)
	}
	return func(t *testing.T) {

		// Delete environment variables
		deleteEnvVar(t, testdata.desc, envAwsKey)
		deleteEnvVar(t, testdata.desc, envAwsSecret)
		deleteEnvVar(t, testdata.desc, envAwsRegion)
	}
}

// TestMain routine for controlling setup/destroy for all auth package tests
func TestMain(m *testing.M) {

	// We need to grab valid AWS credentials from environment variables.
	// If any of these don't exist then we fail.
	key := os.Getenv(testValidAwsKey)
	if key == "" {
		log.Println("The environment variable TESTING_AWS_ACCESS_KEY_ID is not set!!!")
		os.Exit(1)
	}
	secret := os.Getenv(testValidAwsSecret)
	if secret == "" {
		log.Println("The environment variable TESTING_AWS_SECRET_ACCESS_KEY is not set!!!")
		os.Exit(1)
	}
	region := os.Getenv(testValidAwsRegion)
	if region == "" {
		log.Println("The environment variable TESTING_AWS_DEFAULT_REGION is not set!!!")
		os.Exit(1)
	}
	userid := os.Getenv(testValidAwsUserId)
	if userid == "" {
		log.Println("The environment variable TESTING_AWS_USER_ID is not set!!!")
		os.Exit(1)
	}

	// Set the global variable to make the values available for all tests
	AwsCreds.key = key
	AwsCreds.secret = secret
	AwsCreds.region = region
	AwsCreds.userid = userid

	// Run the various test functions
	exitVal := m.Run()

	// Exit
	os.Exit(exitVal)
}

// Test NewSession
func TestNewSession(t *testing.T) {

	// Setup test data
	var tests = []testdef{
		{"No values", false, "", false, "", false, "", false, "", false, "", false, ""},
		{"Just key", true, "fred", false, "", false, "", false, "", false, "", false, ""},
		{"Just secret", false, "", true, "fred", false, "", false, "", false, "", false, ""},
		{"Just region", false, "", false, "", true, "us-west-2", false, "", false, "", false, ""},
		{"Valid credentials", true, AwsCreds.key, true, AwsCreds.secret, true, AwsCreds.region, false, "", false, "", false, AwsCreds.userid},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageTestEnvVar(t, test)
			defer teardownTestCase(t)

			// Run the test
			var actual string
			sess, err := auth.NewSession()
			stsClient := sts.New(sess)
			req := sts.GetCallerIdentityInput{}
			id, _ := stsClient.GetCallerIdentity(&req)
			if id.UserId != nil {
				actual = *id.UserId
			} else {
				actual = ""
			}
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
				internal.Equals(t, test.expectedVal, actual)
			}
		})
	}
}
