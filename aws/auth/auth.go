package auth

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession creates an AWS session using the defaults
//
// Example:
//
//     // Create a new session using the defaults
//     sess := session.NewSession()
func NewSession() (*session.Session, error) {

	sess, err := session.NewSession()
	return sess, err
}

// NewSessionWithConfig creates an AWS session
// using the provided configuration
//
// Example:
//
//     // Create a new session using the provided config
//     sess := session.NewSessionWithConfig()
func NewSessionWithConfig(conf aws.Config) (*session.Session, error) {

	sess, err := session.NewSession(&conf)
	return sess, err
}

// NewSessionWithOptions creates an AWS session
// using the provided options
//
// Example:
//
//     // Create a new session using the provided options
//     sess := session.NewSessionWithOptions(options)
func NewSessionWithOptions(opts session.Options) (*session.Session, error) {

	sess, err := session.NewSessionWithOptions(opts)
	return sess, err
}
