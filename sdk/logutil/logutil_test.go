package logutil_test

import (
	"testing"

	"github.com/bradmccoydev/self-service-sdk/logutil"
)

// TestCase structure for managing test data
type TestCase struct {
	desc       string
	msg        string
	timeFormat string
	logLevel   string
}

// TestData contains the individual test cases to be performed
var TestData = [...]TestCase{
	{"No setup values", "No setup values", "", ""},
	{"Invalid time format", "Invalid time format", "fred", ""},
	{"Invalid time format & invalid log level", "Invalid time format & invalid log level", "fred", "fred"},
	{"Valid setup with no message", "", "UNIX", "TRACE"},
	{"Valid setup with message", "Valid setup with message", "UNIX", "TRACE"},
}

// manageLogConf handles creation of zerolog configuration
// for the test case NewConfigFromEnv
func manageLogConf(t *testing.T, timeFormat string, logLevel string) func(t *testing.T) {
	t.Logf("Setup zerolog configuration for test: %s", string(t.Name()))
	logutil.New(timeFormat, logLevel, false)
	return func(t *testing.T) {
	}
}

// Test LogDebug
func TestLogDebug(t *testing.T) {

	// Iterate through the test data
	for _, test := range TestData {

		t.Run(test.desc, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.timeFormat, test.logLevel)
			defer teardownTestCase(t)

			// Run the test
			logutil.LogDebug(test.msg)
		})
	}
}

// Test LogError
func TestLogError(t *testing.T) {

	// Iterate through the test data
	for _, test := range TestData {

		t.Run(test.desc, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.timeFormat, test.logLevel)
			defer teardownTestCase(t)

			// Run the test
			logutil.LogError(test.msg)
		})
	}
}

// // Test LogFatal
// func TestLogFatal(t *testing.T) {

// 	// Iterate through the test data
// 	for _, test := range TestData {

// 		t.Run(test.desc, func(t *testing.T) {

// 			// Run routine to setup appropriate log config
// 			teardownTestCase := manageLogConf(t, test.timeFormat, test.logLevel)
// 			defer teardownTestCase(t)

// 			// Run the test
// 			logutil.LogFatal(test.msg)
// 		})
// 	}
// }

// Test LogInfo
func TestLogInfo(t *testing.T) {

	// Iterate through the test data
	for _, test := range TestData {

		t.Run(test.desc, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.timeFormat, test.logLevel)
			defer teardownTestCase(t)

			// Run the test
			logutil.LogInfo(test.msg)
		})
	}
}

// Test LogTrace
func TestLogTrace(t *testing.T) {

	// Iterate through the test data
	for _, test := range TestData {

		t.Run(test.desc, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.timeFormat, test.logLevel)
			defer teardownTestCase(t)

			// Run the test
			logutil.LogTrace(test.msg)
		})
	}
}

// Test LogWarn
func TestLogWarn(t *testing.T) {

	// Iterate through the test data
	for _, test := range TestData {

		t.Run(test.desc, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.timeFormat, test.logLevel)
			defer teardownTestCase(t)

			// Run the test
			logutil.LogWarn(test.msg)
		})
	}
}
