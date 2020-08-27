package logutil_test

import (
	"testing"

	"github.com/bradmccoydev/self-service-sdk/logutil"
	"github.com/stretchr/testify/assert"
)

// Test CreateLogConfDef
func TestCreateLogConfDef(t *testing.T) {

	// Setup test data
	var tests = []struct {
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
	assert := assert.New(t)
	for _, test := range tests {

		// Run each test
		_, err := logutil.CreateLogConfDef(test.timeFormat, test.logLevel, test.logToConsole)
		if test.expectErr {
			assert.NotNil(err, test.name)
		} else {
			assert.Nil(err, test.name)
		}
	}
}

// Test LogInfo
func TestLogInfo(t *testing.T) {

	// Setup test default values
	var emptyConfig logutil.LogConfig
	validConfig, _ := logutil.CreateLogConfDef(logutil.TimeFormatUnix, logutil.LogLevelInfo, false)

	// Setup test data
	var tests = []struct {
		name      string
		msg       string
		config    logutil.LogConfig
		expectErr bool
	}{
		{"Empty config", "Empty config", emptyConfig, true},
		{"Valid config but no message", "", *validConfig, false},
		{"Valid config with message", "Valid config with message", *validConfig, false},
	}

	// Iterate through the test data
	//assert := assert.New(t)
	for _, test := range tests {

		// Run each test
		logutil.LogInfo(test.name, test.config)
		// obj, err := logutil.LogInfo(test.name, *test.config)
		// if test.expectErr {
		// 	assert.NotNil(err, test.name)
		// } else {
		// 	assert.Nil(err, test.name)
		// }
	}
}
