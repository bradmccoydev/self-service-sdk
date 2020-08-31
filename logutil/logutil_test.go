package logutil_test

import (
	"testing"

	"github.com/bradmccoydev/self-service-sdk/internal"
	"github.com/bradmccoydev/self-service-sdk/logutil"
)

// TestCase structure for managing test data
type TestCase struct {
	name      string
	msg       string
	genConfig bool
	expectErr bool
}

// TestData contains the individual test cases to be performed
var TestData = [...]TestCase{
	{"No config", "No config", false, false},
	{"Valid config but no message", "", true, false},
	{"Valid config with message", "Valid config with message", true, false},
}

// manageLogConf handles creation of zerolog configuration
// for the test case NewConfigFromEnv
func manageLogConf(t *testing.T, genConfig bool) func(t *testing.T) {
	t.Logf("Setup zerolog configuration for test: %s", string(t.Name()))
	if genConfig == true {
		logConf, _ := logutil.CreateLogConfDef(logutil.TimeFormatUnix, logutil.LogLevelTrace, false)
		logutil.NewWithConfig(logConf)
	} else {
		logutil.New()
	}
	return func(t *testing.T) {
		t.Log("Teardown for NewConfigFromEnv: removing environment variable")
		//os.Unsetenv("TEST_VAR_FRED")
	}
}

// Test CreateLogConfDef
func TestCreateLogConfDef(t *testing.T) {

	// Setup test data
	tests := []struct {
		name         string
		timeFormat   string
		logLevel     string
		logToConsole bool
		expectErr    bool
	}{
		{"No values", "", "", false, true},
		{"Invalid time format", "fred", "", false, true},
		{"Valid time format & invalid log level", logutil.TimeFormatUnix, "fred", false, true},
		{"All valid", logutil.TimeFormatUnix, logutil.LogLevelInfo, false, false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			// Run the test
			_, err := logutil.CreateLogConfDef(test.timeFormat, test.logLevel, test.logToConsole)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test LogDebug
func TestLogDebug(t *testing.T) {

	// Iterate through the test data
	for _, test := range TestData {

		t.Run(test.name, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.genConfig)
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

		t.Run(test.name, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.genConfig)
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

// 		t.Run(test.name, func(t *testing.T) {

// 			// Run routine to setup appropriate log config
// 			teardownTestCase := manageLogConf(t, test.genConfig)
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

		t.Run(test.name, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.genConfig)
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

		t.Run(test.name, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.genConfig)
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

		t.Run(test.name, func(t *testing.T) {

			// Run routine to setup appropriate log config
			teardownTestCase := manageLogConf(t, test.genConfig)
			defer teardownTestCase(t)

			// Run the test
			logutil.LogWarn(test.msg)
		})
	}
}
