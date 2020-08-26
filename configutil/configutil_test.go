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
		{"Valid default config", validDefault, false, ""},
		{"Default config with no key", noKeyDefault, true, ""},
		{"Default config with empty key", emptyKeyDefault, true, ""},
		{"Default config with no value", noValueDefault, false, ""},
		{"Default config with empty value", emptyValueDefault, false, ""},
	}

	// Iterate through the test data
	assert := assert.New(t)
	for _, test := range tests {

		// Run each test
		_, err := configutil.NewConfig(test.defaults)
		if test.expectErr {
			assert.NotNil(err, test.name)
		} else {
			assert.Nil(err, test.name)
		}
	}
}
