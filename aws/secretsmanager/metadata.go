// This file contains all the bits & pieces related to
// fetching metadata about secrets.

package secretsmanager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// DescribeSecret - This function retrieves the metadata about a secret
//
//   Parameters:
//     sess: a valid AWS session
//     secretName: the name of the secret to fetch
//
//   Example:
//     secretString, err := DescribeSecret(mySession, secretName)
func DescribeSecret(sess *session.Session, secretName string) (*secretsmanager.DescribeSecretOutput, error) {

	// Sanity check
	if secretName == "" {
		return nil, newErrorSecretNameNotProvided()
	}

	// Build the input params
	params := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretName),
	}

	// Create the Secrets Manager client
	svc := secretsmanager.New(sess)

	// Make the call to Secrets Manager
	result, err := svc.DescribeSecret(params)

	// Return the result
	return result, err
}
