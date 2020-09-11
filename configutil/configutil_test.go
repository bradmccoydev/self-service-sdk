package configutil_test

import (
	"os"
	"testing"

	"github.com/bradmccoydev/self-service-sdk/configutil"
	"github.com/bradmccoydev/self-service-sdk/internal"
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
		desc        string
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
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			obj, err := configutil.NewConfig(test.defaults)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				actualVal := obj.GetString("configKey")
				internal.NoError(t, err)
				internal.Equals(t, test.expectedVal, actualVal)
			}
		})
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
		desc             string
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
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			obj, err := configutil.NewConfigFromEnv(test.defaults, test.envVars)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				actualVal := obj.GetString(test.configKeyToCheck)
				internal.NoError(t, err)
				internal.Equals(t, test.expectedVal, actualVal)
			}
		})
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

	// Setup test data
	var tests = []struct {
		desc             string
		defaults         []configutil.DefaultValue
		confFile         string
		confType         string
		confPath         string
		expectErr        bool
		configKeyToCheck string
		expectedVal      string
	}{
		{"No values", nil, "", "", "", true, "", ""},
		{"Valid default config with config file name", validDefault, "fred", "", "", true, "", ""},
		{"Valid default config with config file type", validDefault, "", "TOML", "", true, "", ""},
		{"Valid default config with config file path", validDefault, "", "", "/tmp", true, "", ""},
		{"Valid default config with config file name & type", validDefault, "testdata", "YAML", "", true, "", ""},
		{"Default config with no key", noKeyDefault, "testdata", "YAML", ".", true, "", ""},
		{"Default config with empty key", emptyKeyDefault, "testdata", "YAML", ".", true, "", ""},
		{"Default config with no value", noValueDefault, "testdata", "YAML", ".", false, "configKey", ""},
		{"Default config with empty value", emptyValueDefault, "testdata", "YAML", ".", false, "configKey", ""},
		{"Valid call", validDefault, "testdata", "YAML", ".", false, "configKey", "configVal"},
		{"Valid call", validDefault, "testdata", "YAML", ".", false, "metadata.labels.test", "not-default"},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			obj, err := configutil.NewConfigFromFile(test.defaults, test.confFile, test.confType, test.confPath)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				actualVal := obj.GetString(test.configKeyToCheck)
				internal.NoError(t, err)
				internal.Equals(t, test.expectedVal, actualVal)
			}
		})
	}
}
