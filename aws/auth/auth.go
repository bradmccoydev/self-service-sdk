// Package auth provides a simplified api to perform common
// AWS authentication operations.
//
//   The following AWS GoLang SDK packages are used:
//     * aws
//     * aws/session
package auth

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession - This function creates an AWS session using the defaults
//
//   Example:
//     sess, err := session.NewSession()
func NewSession() (*session.Session, error) {

	sess, err := session.NewSession()
	return sess, err
}

// NewSessionWithConfig - This function creates an AWS session
// using the provided configuration
//
//   Example:
//     sess, err := session.NewSessionWithConfig(awsConf)
func NewSessionWithConfig(conf aws.Config) (*session.Session, error) {

	sess, err := session.NewSession(&conf)
	return sess, err
}

// NewSessionWithOptions - This function creates an AWS session
// using the provided options
//
//   Example:
//     sess, err := session.NewSessionWithOptions(options)
func NewSessionWithOptions(opts session.Options) (*session.Session, error) {

	sess, err := session.NewSessionWithOptions(opts)
	return sess, err
}
