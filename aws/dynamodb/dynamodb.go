// Package dynamodb provides a simplified api to perform common
// DynamoDB CRUD operations.
//
//   The following AWS GoLang SDK packages are used:
//     * aws
//     * aws/session
//     * service/dynamodb
//     * service/dynamodb/dynamodbattribute
//     * service/dynamodb/expression
package dynamodb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	// Between operator
	Between string = "BT"

	// BeginsWith operator
	BeginsWith string = "BW"

	// Contains operator
	Contains string = "CO"

	// Equals operator
	Equals string = "EQ"

	// GreaterThan operator
	GreaterThan string = "GT"

	// GreaterThanOrEquals operator
	GreaterThanOrEquals string = "GE"

	// In operator
	In string = "IN"

	// LessThan operator
	LessThan string = "LT"

	// LessThanOrEquals operator
	LessThanOrEquals string = "LE"

	// NotEqual operator
	NotEqual string = "NE"
)

// Field - structure used to specify fields for a projection expression
type Field struct {
	Name string
}

// Condition - structure used for key condition & filter expressions
type Condition struct {
	Field    string
	Operator string
	Value    string
}

// To Do:
// Create an item
// Delete an item
// Update an item

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
//     GetTableArn(mySession, "fred")
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
//     GetTableArn(mySession)
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

// NewExpression - This function creates a new query expression object
//
//   Parameters:
//     keys: an array of key condition(s)
//     filters: an array of filter condition(s)
//     projs: an array of field(s)
//
//   Example:
//     NewExpression(keys, filters, projs)
func NewExpression(keys []Condition, filters []Condition, projs []Field) (expression.Expression, error) {

	// Create new query expression
	var emptyExpr expression.Expression
	builder := expression.NewBuilder()

	// Did we get key conditions?
	if keys != nil {

		// Create a key condition expression
		keyExpr, err := newKeyExpression(keys)
		if err != nil {
			return emptyExpr, err
		}

		// Add the key condition expression
		builder = builder.WithKeyCondition(keyExpr)
	}

	// Did we get filter conditions?
	if filters != nil {

		// Create a filter expression
		filtExpr, err := newFilterExpression(filters)
		if err != nil {
			return emptyExpr, err
		}

		// Add the filter expression
		builder = builder.WithFilter(filtExpr)
	}

	// Did we get a projection?
	if projs != nil {

		// Create a projection expression
		projExpr, err := newProjectionExpression(projs)
		if err != nil {
			return emptyExpr, err
		}

		// Add the projection expression
		builder = builder.WithProjection(projExpr)
	}

	// Build the expression
	expr, err := builder.Build()

	// Return it
	return expr, err
}

// newFilterExpression creates a filter expression for use with a query or scan call
func newFilterExpression(filters []Condition) (expression.ConditionBuilder, error) {

	// Iterate records provided
	var firstTime bool = true
	var err error = nil
	var filterExpr expression.ConditionBuilder
	for _, i := range filters {

		// Sanity check
		if i.Field == "" {
			err = errors.New("Field name must be provided for a filter expression")
			return filterExpr, err
		}
		if i.Operator == "" {
			err = errors.New("Operator must be provided for a filter expression")
			return filterExpr, err
		}

		// Build the condition
		var tmpcond expression.ConditionBuilder
		switch strings.ToUpper(i.Operator) {
		case BeginsWith:
			tmpcond = expression.BeginsWith(expression.Name(i.Field), i.Value)
		case Contains:
			tmpcond = expression.Contains(expression.Name(i.Field), i.Value)
		case Equals:
			tmpcond = expression.Name(i.Field).Equal(expression.Value(i.Value))
		case GreaterThan:
			tmpcond = expression.Name(i.Field).GreaterThan(expression.Value(i.Value))
		case GreaterThanOrEquals:
			tmpcond = expression.Name(i.Field).GreaterThanEqual(expression.Value(i.Value))
		case In:
			tmpcond = expression.Name(i.Field).In(expression.Value(i.Value))
		case LessThan:
			tmpcond = expression.Name(i.Field).LessThan(expression.Value(i.Value))
		case LessThanOrEquals:
			tmpcond = expression.Name(i.Field).LessThanEqual(expression.Value(i.Value))
		case NotEqual:
			tmpcond = expression.Name(i.Field).NotEqual(expression.Value(i.Value))
		default:
			err = fmt.Errorf("Operator type %s is not supported by filter expressions", i.Operator)
			return filterExpr, err
		}

		// First condition?
		if firstTime == true {
			filterExpr = tmpcond
			firstTime = false
		} else {
			filterExpr = filterExpr.And(tmpcond)
		}
	}

	// Return it
	return filterExpr, err
}

// newKeyExpression creates a key condition expression for use with a query call
func newKeyExpression(conditions []Condition) (expression.KeyConditionBuilder, error) {

	// Iterate records provided
	var firstTime bool = true
	var err error = nil
	var keyExpr expression.KeyConditionBuilder
	for _, i := range conditions {

		// Sanity check
		if i.Field == "" {
			err = errors.New("Field name must be provided for a key condition expression")
			return keyExpr, err
		}
		if i.Operator == "" {
			err = errors.New("Operator must be provided for a key condition expression")
			return keyExpr, err
		}

		// Build the condition
		var tmpcond expression.KeyConditionBuilder
		switch strings.ToUpper(i.Operator) {
		case BeginsWith:
			tmpcond = expression.Key(i.Field).BeginsWith(i.Value)
		case Equals:
			tmpcond = expression.Key(i.Field).Equal(expression.Value(i.Value))
		case GreaterThan:
			tmpcond = expression.Key(i.Field).GreaterThan(expression.Value(i.Value))
		case GreaterThanOrEquals:
			tmpcond = expression.Key(i.Field).GreaterThanEqual(expression.Value(i.Value))
		case LessThan:
			tmpcond = expression.Key(i.Field).LessThan(expression.Value(i.Value))
		case LessThanOrEquals:
			tmpcond = expression.Key(i.Field).LessThanEqual(expression.Value(i.Value))
		default:
			err = fmt.Errorf("Operator type %s is not supported by key condition expressions", i.Operator)
			return keyExpr, err
		}

		// First condition?
		if firstTime == true {
			keyExpr = tmpcond
			firstTime = false
		} else {
			keyExpr = keyExpr.And(tmpcond)
		}
	}

	// Return it
	return keyExpr, err
}

// newProjectionExpression create a projection condition (ie restricts the fields returned)
func newProjectionExpression(fields []Field) (expression.ProjectionBuilder, error) {

	// Setup
	var firstTime bool = true
	var err error = nil
	var projExpr expression.ProjectionBuilder
	if len(fields) == 0 {
		err = errors.New("Must provide at least one field for a projection expression")
		return projExpr, err
	}

	// Iterate records provided
	for _, i := range fields {

		// Sanity check
		if i.Name == "" {
			err = errors.New("Field name must be provided for a projection expression")
			return projExpr, err
		}

		// Add the field
		if firstTime == true {
			projExpr = expression.NamesList(expression.Name(i.Name))
			firstTime = false
		} else {
			projExpr = expression.AddNames(projExpr, expression.Name(i.Name))
		}
	}

	// Return it
	return projExpr, err
}

// CreateItem - This function adds a new item to the specified table
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to add the item to
//     expr: the expression object to use
//     newItem: the structure containing the new item properties
//
//   Example:
//     CreateItem(mySession, "fred", expr, myStruct)
func CreateItem(sess *session.Session, tableName string, expr expression.Expression, newItem interface{}) error {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}
	if expr.KeyCondition() == nil {
		err := errors.New("A key condition must be provided in the expression")
		return err
	}

	// Create the DynamoDB client
	//svc := dynamodb.New(sess)

	// // Build the query params
	// params := &dynamodb.QueryInput{
	// 	ExpressionAttributeNames:  expr.Names(),
	// 	ExpressionAttributeValues: expr.Values(),
	// 	FilterExpression:          expr.Filter(),
	// 	KeyConditionExpression:    expr.KeyCondition(),
	// 	ProjectionExpression:      expr.Projection(),
	// 	TableName:                 aws.String(tableName),
	// }

	// // Make the call to DynamoDB
	// result, err := svc.Query(params)

	// // If not ok then bail
	// if err != nil {
	// 	return err
	// }

	// // Massage the result(s) & return
	// //err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &castTo)
	// return err
	return nil
}

// QueryItems - This function makes a query call of the specified table to find matching item(s)
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to query
//     expr: the expression object to use
//     castTo: the array definition that results should be returned in
//
//   Example:
//     QueryItems(mySession, "fred", expr, myArray)
func QueryItems(sess *session.Session, tableName string, expr expression.Expression, castTo interface{}) error {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}
	if expr.KeyCondition() == nil {
		err := errors.New("A key condition must be provided in the expression")
		return err
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Build the query params
	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the call to DynamoDB
	result, err := svc.Query(params)

	// If not ok then bail
	if err != nil {
		return err
	}

	// Massage the result(s) & return
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &castTo)
	return err
}

// ScanItems - This function makes a scan call of the specified table to find matching item(s)
//
//   Parameters:
//     sess: a valid AWS session
//     tableName: the name of the table to scan
//     expr: the expression object to use
//     castTo: the array definition that results should be returned in
//
//   Example:
//     ScanItems(mySession, "fred", expr, myArray)
func ScanItems(sess *session.Session, tableName string, expr expression.Expression, castTo interface{}) error {

	// Sanity check
	if tableName == "" {
		err := errors.New("Table name must be provided")
		return err
	}

	// Create the DynamoDB client
	svc := dynamodb.New(sess)

	// Build the query params
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the call to DynamoDB
	result, err := svc.Scan(params)

	// If not ok then bail
	if err != nil {
		return err
	}

	// Massage the result(s) & return
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &castTo)
	return err
}
