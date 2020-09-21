// This file contains all the bits & pieces related to
// fetching metadata about Dynamo DB tables

package dynamodb

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetTableArn - This function retrieves the Amazon Resource Name (ARN) for the table
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to
//
//   Example:
//     GetTableArn(mySession, "fred")
func GetTableArn(sess *session.Session, tableName string) (*string, error) {

	// Get the table details
	result, err := GetTableDetails(sess, tableName)
	if err != nil {
		return nil, err
	}

	// Check we retrieved something
	if result.Table == nil {
		err := errors.New("Table details were not returned")
		return nil, err
	}

	// Extract the ARN
	arn := result.Table.TableArn

	// Return it
	return arn, nil
}

// GetTableItemCount - This function retrieves the number of items in the table
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to
//
//   Example:
//     GetTableArn(mySession, "fred")
func GetTableItemCount(sess *session.Session, tableName string) (*int64, error) {

	// Get the table details
	result, err := GetTableDetails(sess, tableName)
	if err != nil {
		return nil, err
	}

	// Check we retrieved something
	if result.Table == nil {
		err := errors.New("Table details were not returned")
		return nil, err
	}

	// Extract the item count
	count := result.Table.ItemCount

	// Return it
	return count, nil
}

// GetTableDetails - This function retrieves all the metadata (ARN, item count, keys etc.) about a table
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to
//
//   Example:
//     arn, err := GetTableArn(mySession, "fred")
func GetTableDetails(sess *session.Session, tableName string) (*dynamodb.DescribeTableOutput, error) {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return nil, err
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	request := &dynamodb.DescribeTableInput{}
	request.SetTableName(tableName)
	result, err := svc.DescribeTable(request)

	// If not ok then bail
	if err != nil {
		return nil, err
	}

	// Return it
	return result, nil
}

// GetTableList - This function retrieves a list of available tables
//
//   Parameters:
//     sess: a valid AWS session
//
//   Example:
//     tables, err := GetTableArn(mySession)
func GetTableList(sess *session.Session) ([]string, error) {

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	request := &dynamodb.ListTablesInput{}
	result, err := svc.ListTables(request)

	// If not ok then bail
	if err != nil {
		return nil, err
	}

	// Populate array for returning
	var response []string
	for _, i := range result.TableNames {
		response = append(response, *i)
	}

	// Return it
	return response, nil
}
