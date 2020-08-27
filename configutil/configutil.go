package configutil

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	// ConfigTypeJSON defines the JSON config type
	ConfigTypeJSON string = "JSON"
	// ConfigTypeTOML defines the TOML config type
	ConfigTypeTOML string = "TOML"
	// ConfigTypeYAML defines the YAML config type
	ConfigTypeYAML string = "YAML"
)

// ConfigFile structure represents the metadata for a
// configuration file that Viper can load
type ConfigFile struct {
	Name string
	Path string
	Type string
}

// DefaultValue represents a default configuration value
type DefaultValue struct {
	ConfigKey   string
	ConfigValue string
}

// EnvVariable represents an environment variable and
// the target configuration key to place the value
type EnvVariable struct {
	EnvVar    string
	ConfigKey string
}

// CreateConfFileDef creates a new configuration file instance
// that can then be used to instantiate a configuration manager.
//
// Example:
//
//     // Create a config file instance
//     svc := configutil.CreateConfFileDef(fileName, fileType, filePath)
func CreateConfFileDef(fileName string, fileType string, filePath string) (ConfigFile, error) {

	// Sanity checks
	var conf ConfigFile
	var err error = nil
	if fileName == "" {
		err = errors.New("File name must be provided")
		return conf, err
	}
	if fileType == "" {
		err = errors.New("File type must be provided")
		return conf, err
	}
	switch strings.ToUpper(fileType) {
	case ConfigTypeJSON, ConfigTypeTOML, ConfigTypeYAML:
		// These are ok....
	default:
		err = fmt.Errorf("Unsupported configuration file type: %s", fileType)
		return conf, err
	}

	// Return the configuration object
	conf.Name = fileName
	conf.Path = filePath
	conf.Type = fileType
	return conf, err
}

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

// NewConfig creates a new basic instance of the
// configuration manager.
//
// Example:
//
//     // Create an empty config manager instance
//     svc := configutil.NewConfig(defaults)
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

// NewConfigFromEnv creates a new instance of the
// configuration manager using environment variable(s)
//
// Example:
//
//     // Create an empty config manager instance
//     svc := configutil.NewConfigFromEnv(defaults, envVars)
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

// NewConfigFromFile creates a new instance of the
// configuration manager using a config file
//
// Example:
//
//     // Create an empty config manager instance
//     svc := configutil.NewConfigFromFile(defaults, configfile)
func NewConfigFromFile(defaults []DefaultValue, configfile ConfigFile) (*viper.Viper, error) {

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

	// Try & load the config file
	vCfg.SetConfigName(configfile.Name)
	vCfg.SetConfigType(configfile.Type)
	vCfg.AddConfigPath(configfile.Path)
	err = vCfg.ReadInConfig()

	// Return the configuration object
	return vCfg, err
}
