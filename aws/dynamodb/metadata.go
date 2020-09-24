// This file contains all the bits & pieces related to
// fetching metadata about Dynamo DB tables

package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DescribeTable - This function retrieves all the metadata (ARN, item count, keys etc.) about a table
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to
//
//   Example:
//     info, err := DescribeTable(mySession, "fred")
func DescribeTable(sess *session.Session, tableName string) (*dynamodb.DescribeTableOutput, error) {

	// Sanity check
	if tableName == "" {
		return nil, newErrorTableNameNotProvided()
	}

	// Create a basic input structure for the request
	params := &dynamodb.DescribeTableInput{}
	params = params.SetTableName(tableName)

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	result, err := svc.DescribeTable(params)

	// If not ok then bail
	if err != nil {
		return nil, err
	}

	// Return it
	return result, nil
}

// GetTableArn - This function retrieves the Amazon Resource Name (ARN) for the table
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to
//
//   Example:
//     val, err := GetTableArn(mySession, "fred")
func GetTableArn(sess *session.Session, tableName string) (*string, error) {

	// Get the table details
	result, err := DescribeTable(sess, tableName)
	if err != nil {
		return nil, err
	}

	// Check we retrieved something
	if result.Table == nil {
		return nil, newErrorTableDetailsNotProvided()
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
//     count, err := GetTableItemCount(mySession, "fred")
func GetTableItemCount(sess *session.Session, tableName string) (*int64, error) {

	// Get the table details
	result, err := DescribeTable(sess, tableName)
	if err != nil {
		return nil, err
	}

	// Check we retrieved something
	if result.Table == nil {
		return nil, newErrorTableDetailsNotProvided()
	}

	// Extract the item count
	count := result.Table.ItemCount

	// Return it
	return count, nil
}

// GetTableList - This function retrieves a list of available tables
//
//   Parameters:
//     sess: a valid AWS session
//
//   Example:
//     tables, err := GetTableList(mySession)
func GetTableList(sess *session.Session) ([]string, error) {

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	params := &dynamodb.ListTablesInput{}
	result, err := svc.ListTables(params)

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

// TableExists - This function checks if the specified table exists
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to check
//
//   Example:
//     val, err := TableExists(mySession, "fred")
func TableExists(sess *session.Session, tableName string) (bool, error) {

	// Get the table details
	_, err := DescribeTable(sess, tableName)
	if err != nil {

		// Check error details to check if "real" error
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
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
