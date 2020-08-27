package logutil

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// TimeFormatUnix defines the Unix time format
	TimeFormatUnix string = "UNIX"
	// TimeFormatUnixMs defines the Unix milliseconds time format
	TimeFormatUnixMs string = "UNIXMS"
	// TimeFormatUnixMicro defines the Unix microseconds time format
	TimeFormatUnixMicro string = "UNIXMICRO"
	// LogLevelDebug defines the debug log level
	LogLevelDebug string = "DEBUG"
	// LogLevelError defines the error log level
	LogLevelError string = "ERROR"
	// LogLevelFatal defines the fatal log level
	LogLevelFatal string = "FATAL"
	// LogLevelInfo defines the info log level
	LogLevelInfo string = "INFO"
	// LogLevelTrace defines the trace log level
	LogLevelTrace string = "TRACE"
	// LogLevelWarn defines the warning log level
	LogLevelWarn string = "WARN"
)

// LogConfig structure represents the logging
// configuration to be used
type LogConfig struct {
	TimeFieldFormat string
	GlobalLogLevel  zerolog.Level
	LogToConsole    bool
}

// CreateLogConfDef creates a new logging configuration instance
// that can then be used to instantiate a logging manager.
//
// Example:
//
//     // Create a config file instance
//     conf := logutil.CreateLogConfDef(timeFormat, logLevel, logToConsole)
func CreateLogConfDef(timeFormat string, logLevel string, logToConsole bool) (*LogConfig, error) {

	// Validate time format
	var conf LogConfig
	var err error = nil
	if timeFormat == "" {
		err = errors.New("Time format must be provided")
		return nil, err
	}
	switch strings.ToUpper(timeFormat) {
	case TimeFormatUnix:
		conf.TimeFieldFormat = zerolog.TimeFormatUnix
	case TimeFormatUnixMs:
		conf.TimeFieldFormat = zerolog.TimeFormatUnixMs
	case TimeFormatUnixMicro:
		conf.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	default:
		err = fmt.Errorf("Unsupported time format: %s", timeFormat)
		return nil, err
	}

	// Validate log level
	if logLevel == "" {
		err = errors.New("Log level must be provided")
		return nil, err
	}
	switch strings.ToUpper(logLevel) {
	case LogLevelDebug:
		conf.GlobalLogLevel = zerolog.DebugLevel
	case LogLevelError:
		conf.GlobalLogLevel = zerolog.ErrorLevel
	case LogLevelFatal:
		conf.GlobalLogLevel = zerolog.FatalLevel
	case LogLevelInfo:
		conf.GlobalLogLevel = zerolog.InfoLevel
	case LogLevelTrace:
		conf.GlobalLogLevel = zerolog.TraceLevel
	case LogLevelWarn:
		conf.GlobalLogLevel = zerolog.WarnLevel
	default:
		err = fmt.Errorf("Unsupported log level: %s", logLevel)
		return nil, err
	}

	// Return the configuration object
	conf.LogToConsole = logToConsole
	return &conf, err
}

// LogDebug creates a new debug level log entry
//
// Example:
//
//     // Create an empty config manager instance
//     logutil.LogLogDebug(msg, config)
func LogDebug(msg string, config LogConfig) {
	setup(config)
	log.Debug().Msg(msg)
}

// LogError creates a new error level log entry
//
// Example:
//
//     // Create an empty config manager instance
//     logutil.LogError(msg, config)
func LogError(msg string, config LogConfig) {
	setup(config)
	log.Error().Msg(msg)
}

// LogInfo creates a new info level log entry
//
// Example:
//
//     // Create an empty config manager instance
//     logutil.LogInfo(msg, config)
func LogInfo(msg string, config LogConfig) {
	setup(config)
	log.Info().Msg(msg)
}

// LogTrace creates a new info level log entry
//
// Example:
//
//     // Create an empty config manager instance
//     logutil.LogTrace(msg, config)
func LogTrace(msg string, config LogConfig) {
	setup(config)
	log.Trace().Msg(msg)
}

// setup sets the required configuration for zerolog
func setup(config LogConfig) {

	// Set time format
	if config.TimeFieldFormat == "" {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	} else {
		zerolog.TimeFieldFormat = config.TimeFieldFormat
	}

	// Set global log level
	zerolog.SetGlobalLevel(config.GlobalLogLevel)

	// Set log output
	if config.LogToConsole == false {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}
