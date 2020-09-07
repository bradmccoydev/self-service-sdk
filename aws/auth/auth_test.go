package auth_test

import (
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/bradmccoydev/self-service-sdk/aws/auth"
	"github.com/bradmccoydev/self-service-sdk/internal"
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

// Global variable for AWS credentials
var AwsCreds internal.AwsCreds

// manageTestEnvVar handles creation/desctruction env vars
func manageTestEnvVar(t *testing.T, testdata testdef) func(t *testing.T) {

	// AWS Key
	if testdata.key {
		os.Setenv(internal.EnvAwsKey, testdata.keyval)
	}
	// AWS Secret
	if testdata.secret {
		os.Setenv(internal.EnvAwsSecret, testdata.secval)
	}
	// AWS Region
	if testdata.region {
		os.Setenv(internal.EnvAwsRegion, testdata.regval)
	}
	return func(t *testing.T) {

		// Delete environment variables
		os.Unsetenv(internal.EnvAwsKey)
		os.Unsetenv(internal.EnvAwsSecret)
		os.Unsetenv(internal.EnvAwsRegion)
	}
}

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

// Test NewSession
func TestNewSession(t *testing.T) {

	// Setup test data
	var tests = []testdef{
		{"No values", false, "", false, "", false, "", false, "", false, "", false, ""},
		{"Just key", true, "fred", false, "", false, "", false, "", false, "", false, ""},
		{"Just secret", false, "", true, "fred", false, "", false, "", false, "", false, ""},
		{"Just region", false, "", false, "", true, "us-west-2", false, "", false, "", false, ""},
		{"Valid credentials", true, AwsCreds.Key, true, AwsCreds.Secret, true, AwsCreds.Region, false, "", false, "", false, AwsCreds.Userid},
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
