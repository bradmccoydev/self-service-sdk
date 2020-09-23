// This file contains all the bits & pieces related to
// creating & deleting tables in Dynamo DB

package dynamodb

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	// BillingModePayPerRequest - Pay per request billing mode
	BillingModePayPerRequest string = "PAY_PER_REQUEST"

	// BillingModeProvisioned - Provisioned billing mode
	BillingModeProvisioned string = "PROVISIONED"

	// KeyTypePartition - "partition" key
	KeyTypePartition string = "HASH"

	// KeyTypeSort - "sort" key
	KeyTypeSort string = "RANGE"
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
//     err := CreateTable(mySession, tableConf, tableAttribs)
func CreateTable(sess *session.Session, conf TableConf, attribs []TableAttributes) error {

	// Sanity check
	if conf.TableName == "" {
		return newErrorTableNameNotProvided()
	}
	if conf.BillingMode == "" {
		return newErrorBillingModeNotProvided()
	}
	if len(attribs) == 0 {
		return newErrorTableAttributesNotProvided()
	}

	// Setup Dyanmo objects that we need
	var keys []*dynamodb.KeySchemaElement
	var attribdefs []*dynamodb.AttributeDefinition

	// Process the provided attributes
	var havekey = false
	for _, a := range attribs {

		// Sanity checks
		if a.Name == "" {
			return newErrorTableAttributeNameNotProvided()
		}
		if a.Type == "" {
			return newErrorTableAttributeTypeNotProvided()
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
				return newErrorTableKeyFieldKeyTypeNotProvided(a.Name)
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
		return newErrorTableKeyAttributesNotProvided()
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
	if strings.ToUpper(conf.BillingMode) == BillingModeProvisioned {
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

// DeleteTable - This function deletes the specified table from Dynamo DB
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to delete
//
//   Example:
//     err := DeleteTable(mySession, "fred")
func DeleteTable(sess *session.Session, tableName string) error {

	// Sanity check
	if tableName == "" {
		return newErrorTableNameNotProvided()
	}

	// Build the delete table params
	params := &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	_, err := svc.DeleteTable(params)

	// Return
	return err
}
