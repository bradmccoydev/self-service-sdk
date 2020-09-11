// Package configutil provides a framework for handling configuration parameters.
// It is based on Viper (https://github.com/spf13/viper) and supports setting of
// configuration parameters from a file &/or environment variables.
package configutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	configTypeJSON string = "JSON"
	configTypeTOML string = "TOML"
	configTypeYAML string = "YAML"
)

// DefaultValue - structure to represent the metadata for a configuration default value
type DefaultValue struct {
	ConfigKey   string
	ConfigValue string
}

// EnvVariable - structure to represent an environment variable and
// the target configuration key to place the value
type EnvVariable struct {
	EnvVar    string
	ConfigKey string
}

// Load default values into the configuration
func loadDefaults(config *viper.Viper, defaults []DefaultValue) error {

	// Iterate through each k/v default value
	for _, v := range defaults {
		if v.ConfigKey == "" {
			err := errors.New("Empty config key provided in default values")
			return err
		}
		config.SetDefault(v.ConfigKey, v.ConfigValue)
	}
	return nil
}

// NewConfig - This function creates a new basic instance of a config manager.
//
//   Parameters:
//     defaults: an array of key/value default values
//
//   Example:
//     cfg := configutil.NewConfig(defaults)
func NewConfig(defaults []DefaultValue) (*viper.Viper, error) {

	// Instantiate Viper
	vCfg := viper.New()

	// Handle any default values
	var err error = nil
	if len(defaults) > 0 {
		err = loadDefaults(vCfg, defaults)
	}

	// Return the configuration object
	return vCfg, err
}

// NewConfigFromEnv - This function creates a new instance of the
// config manager and loads in the specified environment variable(s)
//
//   Parameters:
//     defaults: an array of key/value default values
//     envVars: an array of environment variables to load
//
//   Example:
//     cfg := configutil.NewConfig(defaults, envVars)
func NewConfigFromEnv(defaults []DefaultValue, envVars []EnvVariable) (*viper.Viper, error) {

	// Instantiate Viper
	vCfg := viper.New()

	// Handle any default values
	var err error = nil
	if len(defaults) > 0 {
		err = loadDefaults(vCfg, defaults)
		if err != nil {
			return nil, err
		}
	}

	// Try & load the environment variables
	for _, n := range envVars {

		// Sanity check
		key := n.ConfigKey
		env := n.EnvVar
		if env == "" {
			err := errors.New("Empty environment variable provided")
			return nil, err
		}

		// Ok bind the environment variable
		if key == "" {
			err = vCfg.BindEnv(env)
		} else {
			err = vCfg.BindEnv(key, env)
		}
		if err != nil {
			return nil, err
		}
	}

	// Return the configuration object
	return vCfg, err
}

// NewConfigFromFile - This function creates a new instance of the
// configuration manager using a config file
//
//   Parameters:
//     defaults: an array of key/value default values
//     fileName: the name of the config file to load
//     fileType: the type (JSON, TOML, YAML) of the config file
//     filePath: the path to the config file
//
//   Example:
//     cfg := configutil.NewConfig(defaults, "fred", "yaml", "/tmp")
func NewConfigFromFile(defaults []DefaultValue, fileName string, fileType string, filePath string) (*viper.Viper, error) {

	// Sanity checks
	var err error = nil
	if fileName == "" {
		err = errors.New("File name must be provided")
		return nil, err
	}
	if fileType == "" {
		err = errors.New("File type must be provided")
		return nil, err
	}
	switch strings.ToUpper(fileType) {
	case configTypeJSON, configTypeTOML, configTypeYAML:
		// These are ok....
	default:
		err = fmt.Errorf("Unsupported configuration file type: %s", fileType)
		return nil, err
	}

	// Instantiate Viper
	vCfg := viper.New()

	// Handle any default values
	if len(defaults) > 0 {
		err = loadDefaults(vCfg, defaults)
		if err != nil {
			return nil, err
		}
	}

	// Try & load the config file
	typeLower := strings.ToLower(fileType)
	vCfg.SetConfigName(fileName)
	vCfg.SetConfigType(typeLower)
	vCfg.AddConfigPath(filePath)
	err = vCfg.ReadInConfig()

	// Return the configuration object
	return vCfg, err
}
