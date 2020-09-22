// This file contains all the bits & pieces related to
// creating & deleting tables in Dynamo DB

package dynamodb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// TableAttributes - structure used to represent table attributes
type TableAttributes struct {
	Name    string
	Type    string
	IsKey   bool
	KeyType string
}

// TableConf - structure used to represent table config metadata
type TableConf struct {
	TableName          string
	BillingMode        string
	ReadCapacityUnits  int64
	WriteCapacityUnits int64
}

// CreateTable - This function creates a new table in Dynamo DB
//
//   Parameters:
//     sess: a valid AWS session
//     conf: the configuration metadata for the table
//     attribs: an array of the attributes the table should have
//
//   Example:
//     err := CreateTable(mySession, "fred", myStruct)
func CreateTable(sess *session.Session, conf TableConf, attribs []TableAttributes) error {

	// Sanity check
	if conf.TableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}
	if conf.BillingMode == "" {
		err := errors.New("Billing mode must be provided")
		return err
	}
	if len(attribs) == 0 {
		err := errors.New("Table attributes must be provided")
		return err
	}

	// Setup Dyanmo objects that we need
	var keys []*dynamodb.KeySchemaElement
	var attribdefs []*dynamodb.AttributeDefinition

	// Process the provided attributes
	var havekey = false
	for _, a := range attribs {

		// Sanity checks
		if a.Name == "" {
			err := errors.New("Table attribute must have a name")
			return err
		}
		if a.Type == "" {
			err := errors.New("Table attribute must have a type")
			return err
		}

		// Add the attribute definition
		adef := dynamodb.AttributeDefinition{
			AttributeName: aws.String(a.Name),
			AttributeType: aws.String(a.Type),
		}
		attribdefs = append(attribdefs, &adef)

		// Handle if a key field
		if a.IsKey {

			// Make sure we have a key type
			if a.KeyType == "" {
				err := fmt.Errorf("The key field %s did not include a key type", a.Name)
				return err
			}

			// Add the key element
			kdef := dynamodb.KeySchemaElement{
				AttributeName: aws.String(a.Name),
				KeyType:       aws.String(a.KeyType),
			}
			keys = append(keys, &kdef)
			havekey = true
		}
	}

	// Check we have at least one key
	if havekey == false {
		err := errors.New("No key attributes were provided")
		return err
	}

	// Create a basic input structure for the request
	params := &dynamodb.CreateTableInput{}

	// Add the table name
	params = params.SetTableName(conf.TableName)

	// Add the key elements
	params = params.SetKeySchema(keys)

	// Add the attribute definitions
	params = params.SetAttributeDefinitions(attribdefs)

	// Add the table billing mode
	params = params.SetBillingMode(conf.BillingMode)

	// Add the capacity units if billing mode is provisioned
	if strings.ToUpper(conf.BillingMode) == "PROVISIONED" {
		thruput := dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  &conf.ReadCapacityUnits,
			WriteCapacityUnits: &conf.WriteCapacityUnits,
		}
		params = params.SetProvisionedThroughput(&thruput)
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	_, err := svc.CreateTable(params)

	// Return
	return err
}

// // DeleteItem - This function deletes an item from the specified table
// //
// //   Parameters:
// //     sess: a valid AWS session
// //     tableName: the name of the table to delete the item from
// //     input: the structure containing the key values for the item to be deleted
// //
// //   Example:
// //     err := DeleteItem(mySession, "fred", myStruct)
// func DeleteItem(sess *session.Session, tableName string, input interface{}) error {

// 	// Sanity check
// 	if tableName == "" {
// 		err := errors.New("Table name must be provided")
// 		return err
// 	}

// 	// Create the DynamoDB client
// 	svc := dynamodb.New(sess)

// 	// Marshall the input
// 	item, err := dynamodbattribute.MarshalMap(&input)

// 	// If not ok then bail
// 	if err != nil {
// 		return err
// 	}

// 	// Build the delete params
// 	params := &dynamodb.DeleteItemInput{
// 		Key:       item,
// 		TableName: aws.String(tableName),
// 	}

// 	// Make the call to DynamoDB
// 	_, err = svc.DeleteItem(params)

// 	// Return
// 	return err
// }
