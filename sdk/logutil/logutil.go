// Package logutil provides a framework for handling logging.
// It is based on Zerolog (https://github.com/rs/zerolog).
package logutil

import (
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

// LogDebug creates a new debug level log entry
//
//   Example:
//     logutil.LogLogDebug(msg)
func LogDebug(msg string) {
	log.Debug().Msg(msg)
}

// LogError creates a new error level log entry
//
//   Example:
//     logutil.LogError(msg)
func LogError(msg string) {
	log.Error().Msg(msg)
}

// LogFatal creates a new fatal level log entry
//
//   Example:
//     logutil.LogFatal(msg)
func LogFatal(msg string) {
	log.Fatal().Msg(msg)
}

// LogInfo creates a new info level log entry
//
//   Example:
//     logutil.LogInfo(msg)
func LogInfo(msg string) {
	log.Info().Msg(msg)
}

// LogTrace creates a new trace level log entry
//
//   Example:
//     logutil.LogTrace(msg)
func LogTrace(msg string) {
	log.Trace().Msg(msg)
}

// LogWarn creates a new warning level log entry
//
//   Example:
//     logutil.LogWarn(msg)
func LogWarn(msg string) {
	log.Warn().Msg(msg)
}

// New - This function creates a new logger
//
//   Parameters:
//     timeFormat: the format to use for the log timestamps. Defaults to UNIX
//     logLevel: the logging level (ie log messages must be a higher level to be written). Defaults to INFO
//     logToConsole: not currently implemented
//
//   Example:
//     logutil.New("UNIXMS", "INFO", false)
func New(timeFormat string, logLevel string, logToConsole bool) {

	// Handle time format
	switch strings.ToUpper(timeFormat) {
	case TimeFormatUnix:
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	case TimeFormatUnixMs:
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	case TimeFormatUnixMicro:
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	default:
		// Unknown format - so we default
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}

	// Handle log level
	if logLevel == "" {
		logLevel = LogLevelInfo
	}
	switch strings.ToUpper(logLevel) {
	case LogLevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LogLevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LogLevelFatal:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case LogLevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LogLevelTrace:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case LogLevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		// Unknown log level - so we default
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set log output
	if logToConsole == false {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}
