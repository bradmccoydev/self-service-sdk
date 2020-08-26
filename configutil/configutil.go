package configutil

import (
	"errors"

	"github.com/spf13/viper"
)

// DefaultValue represents a default configuration value
type DefaultValue struct {
	ConfigKey   string
	ConfigValue string
}

func loadDefaults(config *viper.Viper, defaults []DefaultValue) error {

	// Iterate through each k/v default value
	for _, v := range defaults {
		if v.ConfigKey == "" {
			err := errors.New("Empty config key provided")
			return err
		}
		config.SetDefault(v.ConfigKey, v.ConfigValue)
	}
	return nil
}

// NewConfig creates a new basic instance of the
// configuration manager.
//
// Example:
//
//     // Create an empty config manager instance
//     svc := configutil.NewConfig(defaults)
func NewConfig(defaults []DefaultValue) (*viper.Viper, error) {

	vCfg := viper.New()

	// Handle any default values
	var err error = nil
	if len(defaults) > 0 {
		err = loadDefaults(vCfg, defaults)
	}

	// Return the configuration object
	return vCfg, err
}

// NewConfigFromFile creates a new instance of the
// configuration manager using a config file
//
// Example:
//
//     // Create an empty config manager instance
//     svc := configutil.NewConfigFromFile(defaults, configfile)
func NewConfigFromFile(defaults []DefaultValue) (*viper.Viper, error) {

	vCfg := viper.New()

	// Handle any default values
	var err error = nil
	if len(defaults) > 0 {
		err = loadDefaults(vCfg, defaults)
	}

	// Return the configuration object
	return vCfg, err
}
