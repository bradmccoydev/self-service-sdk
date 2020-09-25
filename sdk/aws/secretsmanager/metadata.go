// This file contains all the bits & pieces related to
// fetching metadata about secrets.

package secretsmanager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

// GetSecretArn - This function retrieves the Amazon Resource Name (ARN) for the secret
//
//   Parameters:
//     sess: a valid AWS session
//     secretName: the name of the secret to get the ARN for
//
//   Example:
//     val, err := GetSecretArn(mySession, "fred")
func GetSecretArn(sess *session.Session, secretName string) (*string, error) {

	// Get the table details
	result, err := DescribeSecret(sess, secretName)
	if err != nil {
		return nil, err
	}

	// Extract the ARN
	arn := result.ARN

	// Return it
	return arn, err
}

// SecretExists - This function checks if the specified secret exists
//
//   Parameters:
//     sess: a valid AWS session
//     secretName: the name of the secret to check
//
//   Example:
//     val, err := SecretExists(mySession, "fred")
func SecretExists(sess *session.Session, secretName string) (bool, error) {

	// Get the table details
	_, err := DescribeSecret(sess, secretName)
	if err != nil {

		// Check error details to check if "real" error
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceNotFoundException:
				return false, nil
			default:
				return false, err
			}
		}

		//if err.
		return false, err
	}

	// Return it
	return true, nil
}
