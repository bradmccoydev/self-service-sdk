package configutil_test

import (
	"testing"

	"github.com/bradmccoydev/self-service-sdk/configutil"
	"github.com/stretchr/testify/assert"
)

// Test NewConfig
func TestNewConfig(t *testing.T) {

	// Setup test default values
	validDefault := []configutil.DefaultValue{{
		ConfigKey:   "configKey",
		ConfigValue: "configVal",
	}}
	noKeyDefault := []configutil.DefaultValue{{
		ConfigValue: "configVal",
	}}
	emptyKeyDefault := []configutil.DefaultValue{{
		ConfigKey:   "",
		ConfigValue: "configVal",
	}}
	noValueDefault := []configutil.DefaultValue{{
		ConfigKey: "configKey",
	}}
	emptyValueDefault := []configutil.DefaultValue{{
		ConfigKey:   "configKey",
		ConfigValue: "",
	}}

	// Setup test data
	var tests = []struct {
		name        string
		defaults    []configutil.DefaultValue
		expectErr   bool
		expectedVal string
	}{
		{"No default configs", nil, false, ""},
		{"Valid default config", validDefault, false, "configVal"},
		{"Default config with no key", noKeyDefault, true, ""},
		{"Default config with empty key", emptyKeyDefault, true, ""},
		{"Default config with no value", noValueDefault, false, ""},
		{"Default config with empty value", emptyValueDefault, false, ""},
	}

	// Iterate through the test data
	assert := assert.New(t)
	for _, test := range tests {

		// Run each test
		obj, err := configutil.NewConfig(test.defaults)
		if test.expectErr {
			assert.NotNil(err, test.name)
		} else {
			actualVal := obj.GetString("configKey")
			assert.Nil(err, test.name)
			assert.Equal(test.expectedVal, actualVal, test.name)
		}
	}
}

// Test NewConfigFile
func TestNewConfigFile(t *testing.T) {

	// Setup test data
	var tests = []struct {
		name        string
		fileName    string
		fileType    string
		filePath    string
		expectErr   bool
		expectedVal string
	}{
		{"No values", "", "", "", true, ""},
		{"Just fileName", "fred", "", "", true, ""},
		{"With fileName & filePath", "fred", "", "fred", true, ""},
		{"With all - invalid fileType", "fred", "fred", "fred", true, ""},
		{"With all - valid fileType", "fred", "json", "fred", false, ""},
	}

	// Iterate through the test data
	assert := assert.New(t)
	for _, test := range tests {

		// Run each test
		_, err := configutil.NewConfigFile(test.fileName, test.fileType, test.filePath)
		if test.expectErr {
			assert.NotNil(err, test.name)
		} else {
			assert.Nil(err, test.name)
		}
	}
}

// Test NewConfigFromEnv
func NewConfigFromEnv(t *testing.T) {

	// Setup test default values
	validDefault := []configutil.DefaultValue{
		{
			ConfigKey:   "configKey",
			ConfigValue: "configVal",
		},
		{
			ConfigKey:   "shell",
			ConfigValue: "default",
		}}
	noKeyDefault := []configutil.DefaultValue{{
		ConfigValue: "configVal",
	}}
	emptyKeyDefault := []configutil.DefaultValue{{
		ConfigKey:   "",
		ConfigValue: "configVal",
	}}
	noValueDefault := []configutil.DefaultValue{{
		ConfigKey: "configKey",
	}}
	emptyValueDefault := []configutil.DefaultValue{{
		ConfigKey:   "configKey",
		ConfigValue: "",
	}}

	// Setup test environment variable values
	var emptyConfFile configutil.EnvVariable
	validConfFile, _ := configutil.NewConfigFile("testdata", "yaml", ".")
	missingConfFile, _ := configutil.NewConfigFile("testdata2", "yaml", ".")
	invalidConfFileType, _ := configutil.NewConfigFile("testdata", "json", ".")

	// Setup test data
	var tests = []struct {
		name             string
		defaults         []configutil.DefaultValue
		envVars          []configutil.EnvVariable
		expectErr        bool
		configKeyToCheck string
		expectedVal      string
	}{
		{"No values", nil, nil, false, "", ""},
		{"Default config with no key", noKeyDefault, emptyConfFile, true, "", ""},
		{"Default config with empty key", emptyKeyDefault, emptyConfFile, true, "", ""},
		{"Default config with no value", noValueDefault, validConfFile, false, "configKey", ""},
		{"Default config with empty value", emptyValueDefault, validConfFile, false, "configKey", ""},
		{"Valid default config with empty config file", validDefault, emptyConfFile, true, "", ""},
		{"Valid default config with missing config file", validDefault, missingConfFile, true, "", ""},
		{"Valid default config with invalid config file type", validDefault, invalidConfFileType, true, "", ""},
		{"Valid call", validDefault, validConfFile, false, "configKey", "configVal"},
		{"Valid call", validDefault, validConfFile, false, "metadata.labels.test", "not-default"},
	}

	// Iterate through the test data
	assert := assert.New(t)
	for _, test := range tests {

		// Run each test
		obj, err := configutil.NewConfigFromFile(test.defaults, test.confFile)
		if test.expectErr {
			assert.NotNil(err, test.name)
		} else {
			actualVal := obj.GetString(test.configKeyToCheck)
			assert.Nil(err, test.name)
			assert.Equal(test.expectedVal, actualVal, test.name)
		}
	}
}

// Test NewConfigFromFile
func TestNewConfigFromFile(t *testing.T) {

	// Setup test default values
	validDefault := []configutil.DefaultValue{
		{
			ConfigKey:   "configKey",
			ConfigValue: "configVal",
		},
		{
			ConfigKey:   "metadata.labels.test",
			ConfigValue: "default",
		}}
	noKeyDefault := []configutil.DefaultValue{{
		ConfigValue: "configVal",
	}}
	emptyKeyDefault := []configutil.DefaultValue{{
		ConfigKey:   "",
		ConfigValue: "configVal",
	}}
	noValueDefault := []configutil.DefaultValue{{
		ConfigKey: "configKey",
	}}
	emptyValueDefault := []configutil.DefaultValue{{
		ConfigKey:   "configKey",
		ConfigValue: "",
	}}

	// Setup test config file values
	var emptyConfFile configutil.ConfigFile
	validConfFile, _ := configutil.NewConfigFile("testdata", "yaml", ".")
	missingConfFile, _ := configutil.NewConfigFile("testdata2", "yaml", ".")
	invalidConfFileType, _ := configutil.NewConfigFile("testdata", "json", ".")

	// Setup test data
	var tests = []struct {
		name             string
		defaults         []configutil.DefaultValue
		confFile         configutil.ConfigFile
		expectErr        bool
		configKeyToCheck string
		expectedVal      string
	}{
		{"No values", nil, emptyConfFile, true, "", ""},
		{"Default config with no key", noKeyDefault, emptyConfFile, true, "", ""},
		{"Default config with empty key", emptyKeyDefault, emptyConfFile, true, "", ""},
		{"Default config with no value", noValueDefault, validConfFile, false, "configKey", ""},
		{"Default config with empty value", emptyValueDefault, validConfFile, false, "configKey", ""},
		{"Valid default config with empty config file", validDefault, emptyConfFile, true, "", ""},
		{"Valid default config with missing config file", validDefault, missingConfFile, true, "", ""},
		{"Valid default config with invalid config file type", validDefault, invalidConfFileType, true, "", ""},
		{"Valid call", validDefault, validConfFile, false, "configKey", "configVal"},
		{"Valid call", validDefault, validConfFile, false, "metadata.labels.test", "not-default"},
	}

	// Iterate through the test data
	assert := assert.New(t)
	for _, test := range tests {

		// Run each test
		obj, err := configutil.NewConfigFromFile(test.defaults, test.confFile)
		if test.expectErr {
			assert.NotNil(err, test.name)
		} else {
			actualVal := obj.GetString(test.configKeyToCheck)
			assert.Nil(err, test.name)
			assert.Equal(test.expectedVal, actualVal, test.name)
		}
	}
}
