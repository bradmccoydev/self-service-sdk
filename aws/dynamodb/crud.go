// This file contains all the bits & pieces related to
// creating, reading updating & deleting items from
// Dynamo DB tables

package dynamodb

import (
	"errors"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// CreateItem - This function adds a new item to the specified table
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to add the item to
//     input: the structure containing the new item properties
//
//   Example:
//     err := CreateItem(mySession, "fred", myStruct)
func CreateItem(sess *session.Session, tableName string, input interface{}) error {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}

	// Marshall the input
	item, err := dynamodbattribute.MarshalMap(&input)

	// If not ok then bail
	if err != nil {
		return err
	}

	// Build the input params
	params := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	_, err = svc.PutItem(params)

	// Return
	return err
}

// DeleteItem - This function deletes an item from the specified table
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to delete the item from
//     input: the structure containing the key values for the item to be deleted
//
//   Example:
//     err := DeleteItem(mySession, "fred", myStruct)
func DeleteItem(sess *session.Session, tableName string, input interface{}) error {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}

	// Marshall the input
	item, err := dynamodbattribute.MarshalMap(&input)

	// If not ok then bail
	if err != nil {
		return err
	}

	// Build the delete params
	params := &dynamodb.DeleteItemInput{
		Key:       item,
		TableName: aws.String(tableName),
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	_, err = svc.DeleteItem(params)

	// Return
	return err
}

// QueryItems - This function makes a query call of the specified table to find matching item(s)
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to query
//     expr: the expression object to use
//     response: the array definition that the results should be returned in
//
//   Example:
//     err := QueryItems(mySession, "fred", expr, myArray)
func QueryItems(sess *session.Session, tableName string, expr expression.Expression, response interface{}) error {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}
	if expr.KeyCondition() == nil {
		err := errors.New("A key condition must be provided in the expression")
		return err
	}

	// Build the query params
	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	result, err := svc.Query(params)

	// If not ok then bail
	if err != nil {
		return err
	}

	// Massage the result(s) & return
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &response)
	return err
}

// ScanItems - This function makes a scan call of the specified table to find matching item(s)
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to scan
//     expr: the expression object to use
//     response: the array definition that results should be returned in
//
//   Example:
//     err := ScanItems(mySession, "fred", expr, myArray)
func ScanItems(sess *session.Session, tableName string, expr expression.Expression, response interface{}) error {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}

	// Build the query params
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	result, err := svc.Scan(params)

	// If not ok then bail
	if err != nil {
		return err
	}

	// Massage the result(s) & return
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &response)
	return err
}

// UpdateItem - This function updates an item in the specified table
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to update
//     keys: the structure containing the item keys
//     input: the structure containing the item properties to update
//
//   Example:
//     err := UpdateItem(mySession, "fred", myStruct)
func UpdateItem(sess *session.Session, tableName string, keys interface{}, input interface{}) error {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}

	// Marshall the keys
	itemKeys, err := dynamodbattribute.MarshalMap(&keys)

	// If not ok then bail
	if err != nil {
		return err
	}

	// We need to iterate through the input interface fields
	var update expression.UpdateBuilder
	u := reflect.ValueOf(&input).Elem()
	t := u.Type()
	for i := 0; i < u.NumField(); i++ {

		// Check if it is empty
		f := u.Field(i)
		if !reflect.DeepEqual(f.Interface(), reflect.Zero(f.Type()).Interface()) {

			// Get json field name
			jsonFieldName := t.Field(i).Name
			if jsonTag := t.Field(i).Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
				if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
					jsonFieldName = jsonTag[:commaIdx]
				}
			}
			// Build the update definition
			update = update.Set(expression.Name(jsonFieldName), expression.Value(f.Interface()))
		}
	}

	// Create an update expression
	builder := expression.NewBuilder().WithUpdate(update)
	expression, err := builder.Build()

	// If not ok then bail
	if err != nil {
		return err
	}

	// Build the update params
	params := &dynamodb.UpdateItemInput{
		Key:                       itemKeys,
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		UpdateExpression:          expression.Update(),
		ReturnValues:              aws.String("UPDATED_NEW"),
		TableName:                 aws.String(tableName),
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Make the call to DynamoDB
	_, err = svc.UpdateItem(params)

	// Return
	return err
}
