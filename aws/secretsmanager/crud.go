// This file contains all the bits & pieces related to
// creating, reading updating & deleting secrets from
// Secrets Manager

package secretsmanager

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	guuid "github.com/google/uuid"
)

// secretDetails - structure used to manage secret details
type secretDetails struct {
	name      string
	desc      string
	secBin    []byte
	secString string
}

// SecretKeyValue - structure used to specify a secret key-value
type SecretKeyValue struct {
	Key   string
	Value string
}

// createSecret - This function creates a secret
func createSecret(sess *session.Session, secret secretDetails) error {

	// Sanity check
	if len(secret.secBin) == 0 && secret.secString == "" {
		return newErrorSecretBinaryAndStringNotProvided()
	}

	// Generate a UUID
	id := guuid.New().String()

	// Build the input params
	params := &secretsmanager.CreateSecretInput{
		ClientRequestToken: aws.String(id),
		Description:        aws.String(secret.desc),
		Name:               aws.String(secret.name),
		SecretString:       aws.String(secret.secString),
		SecretBinary:       secret.secBin,
	}

	// Create the Secrets Manager client
	svc := secretsmanager.New(sess)

	// Make the call to Secrets Manager
	_, err := svc.CreateSecret(params)

	// Return the result
	return err
}

// CreateSecretKeyValue - This function creates a secret string
//
//   Parameters:
//     sess: a valid AWS session
//     name: the name of the secret to create
//     description: the description for the secret
//     secret: a hashmap of key/value secret pairs
//
//   Example:
//     err := CreateSecretKeyValue(mySession, secretName, secretDesc, secMap)
func CreateSecretKeyValue(sess *session.Session, name string, description string, secret map[string]string) error {

	// Sanity check
	if name == "" {
		return newErrorSecretNameNotProvided()
	}
	if len(secret) == 0 {
		return newErrorSecretMapNotProvided()
	}

	// Convert map to json
	jsonByte, err := json.Marshal(secret)
	if err != nil {
		return err
	}

	// Call the routine to create the secret
	secJSON := string(jsonByte)
	secDetails := secretDetails{
		name:      name,
		desc:      description,
		secString: secJSON,
	}
	err = createSecret(sess, secDetails)

	// Return the result
	return err
}

// CreateSecretString - This function creates a secret string
//
//   Parameters:
//     sess: a valid AWS session
//     name: the name of the secret to create
//     description: the description for the secret
//     secret: the secret string
//
//   Example:
//     err := CreateSecretString(mySession, secretName, secretDesc, secretString)
func CreateSecretString(sess *session.Session, name string, description string, secret string) error {

	// Sanity check
	if name == "" {
		return newErrorSecretNameNotProvided()
	}
	if secret == "" {
		return newErrorSecretStringNotProvided()
	}

	// Call the routine to create the secret
	secDetails := secretDetails{
		name:      name,
		desc:      description,
		secString: secret,
	}
	err := createSecret(sess, secDetails)

	// Return the result
	return err
}

// DeleteSecret - This function deletes a secret
//
//   Parameters:
//     sess: a valid AWS session
//     secretName: the name of the secret to delete
//
//   Example:
//     err := DeleteSecret(mySession, secretName)
func DeleteSecret(sess *session.Session, secretName string) error {

	// Sanity check
	if secretName == "" {
		return newErrorSecretNameNotProvided()
	}

	// Build the input params
	params := &secretsmanager.DeleteSecretInput{
		SecretId: aws.String(secretName),
	}

	// Create the Secrets Manager client
	svc := secretsmanager.New(sess)

	// Make the call to Secrets Manager
	_, err := svc.DeleteSecret(params)

	// Return the result
	return err
}

// getSecret - This function provides a generic routine to retrieve a secret
func getSecret(sess *session.Session, secretName string) (*secretsmanager.GetSecretValueOutput, error) {

	// Sanity check
	if secretName == "" {
		return nil, newErrorSecretNameNotProvided()
	}

	// Build the input params
	params := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	// Create the Secrets Manager client
	svc := secretsmanager.New(sess)

	// Make the call to Secrets Manager
	result, err := svc.GetSecretValue(params)

	// Return the result
	return result, err
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
