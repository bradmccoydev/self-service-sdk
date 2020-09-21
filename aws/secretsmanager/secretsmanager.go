// Package secretsmanager provides a simplified api to perform common
// Secrets Manager CRUD operations.
//
//   The following AWS GoLang SDK packages are used:
//     * aws
//     * aws/session
//     * service/secretsmanager
package secretsmanager

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SecretKeyValue - structure used to specify a secret key-value
type SecretKeyValue struct {
	Key   string
	Value string
}

// getSecret - This function provides a generic routine to retrieve a secret
func getSecret(sess *session.Session, secretName string) (*secretsmanager.GetSecretValueOutput, error) {

	// Create the Secrets Manager client
	svc := secretsmanager.New(sess)

	// Make the call to Secrets Manager
	request := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}
	result, err := svc.GetSecretValue(request)

	// Return the result
	return result, err
}

// GetSecretString - This function retrieves a secret string
//
//   Parameters:
//     sess: a valid AWS session
//     secretName: the name of the secret to fetch
//
//   Example:
//     secretString, err := GetSecretString(mySession, secretName)
func GetSecretString(sess *session.Session, secretName string) (string, error) {

	// Extract the string
	result, err := getSecret(sess, secretName)
	var value string
	if err == nil {
		value = *result.SecretString
	}

	// Return the result
	return value, err
}

// GetSecretKeyValue - This function retrieves a secret key-value pair
//
//   Parameters:
//     sess: a valid AWS session
//     secretName: the name of the secret to fetch
//
//   Example:
//     secretKV, err := GetSecretKeyValue(mySession, secretName)
func GetSecretKeyValue(sess *session.Session, secretName string) (SecretKeyValue, error) {

	result, err := getSecret(sess, secretName)
	var keyval SecretKeyValue
	if err == nil {

		// Parse the result
		str := strings.Trim(string(*result.SecretString), "{}")
		arr := strings.Split(str, ":")

		// Get the key
		keyval.Key = arr[0][1 : len(arr[0])-1]

		// Get the value
		keyval.Value = arr[1][1 : len(arr[1])-1]
	}

	// Return stuff
	return keyval, err
}
