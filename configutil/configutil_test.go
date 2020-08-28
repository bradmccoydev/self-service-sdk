package configutil_test

import (
	"os"
	"testing"

	"github.com/bradmccoydev/self-service-sdk/configutil"
	"github.com/stretchr/testify/assert"
)

// manageTestEnvVar handles creation/desctruction of an
// environment variable for the test case NewConfigFromEnv
func manageTestEnvVar(t *testing.T) func(t *testing.T) {
	t.Log("Setup for NewConfigFromEnv: creating environment variable")
	os.Setenv("TEST_VAR_FRED", "BlahBlahBlah")
	return func(t *testing.T) {
		t.Log("Teardown for NewConfigFromEnv: removing environment variable")
		os.Unsetenv("TEST_VAR_FRED")
	}
}

// Test CreateConfFileDef
func TestCreateConfFileDef(t *testing.T) {

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
		_, err := configutil.CreateConfFileDef(test.fileName, test.fileType, test.filePath)
		if test.expectErr {
			assert.NotNil(err, test.name)
		} else {
			assert.Nil(err, test.name)
		}
	}
}

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

// Test NewConfigFromEnv
func TestNewConfigFromEnv(t *testing.T) {

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
	var emptyEnvVars []configutil.EnvVariable
	missingEnvVar := []configutil.EnvVariable{{
		ConfigKey: "configKey",
	}}
	invalidEnvVar := []configutil.EnvVariable{{
		EnvVar:    "FrEdd1EB0y",
		ConfigKey: "configKey",
	}}
	validEnvVar := []configutil.EnvVariable{{
		EnvVar:    "TEST_VAR_FRED",
		ConfigKey: "shell",
	}}

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
		{"Default config with no key", noKeyDefault, nil, true, "", ""},
		{"Default config with empty key", emptyKeyDefault, nil, true, "", ""},
		{"Default config with no value", noValueDefault, nil, false, "configKey", ""},
		{"Default config with empty value", emptyValueDefault, nil, false, "configKey", ""},
		{"Valid default config with no env vars", validDefault, nil, false, "configKey", "configVal"},
		{"Valid default config with empty env vars array", validDefault, emptyEnvVars, false, "configKey", "configVal"},
		{"Valid default config with missing env var", validDefault, missingEnvVar, true, "", ""},
		{"Valid default config with invalid env var", validDefault, invalidEnvVar, false, "configKey", "configVal"},
		{"Valid call", validDefault, validEnvVar, false, "shell", "BlahBlahBlah"},
	}

	// Run setup routine to configure test env var
	teardownTestCase := manageTestEnvVar(t)
	defer teardownTestCase(t)

	// Iterate through the test data
	assert := assert.New(t)
	for _, test := range tests {

		// Run each test
		obj, err := configutil.NewConfigFromEnv(test.defaults, test.envVars)
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
	validConfFile, _ := configutil.CreateConfFileDef("testdata", "yaml", ".")
	missingConfFile, _ := configutil.CreateConfFileDef("testdata2", "yaml", ".")
	invalidConfFileType, _ := configutil.CreateConfFileDef("testdata", "json", ".")

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
