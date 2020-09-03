package dynamodb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
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

// TableName structure
type TableName struct {
	TableName string
}

// Field structure manages fields for a projection
type Field struct {
	Name string
}

// Filter structure manages filter expressions
type Filter struct {
	Field    string
	Operator string
	Value    string
}

// To Do:
// Create a projection expression
// Perform table scan
// - with all options
// - with just filter
// - with just projection
// Get a specific item
// Update an item
// Delete an item

// GetTableArn - retrieves the ARN for the table
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

// GetTableItemCount - retrieves the number of items
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

// GetTableDetails - retrieves the details for a table
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

// GetTableList - retrieves a list of tables
func GetTableList(sess *session.Session) ([]TableName, error) {

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
	var response []TableName
	for _, i := range result.TableNames {
		response = append(response, TableName{
			*i,
		})
	}

	// Return it
	return response, nil
}

// NewFilterExpression handles creation of a filter expression
func NewFilterExpression(filters []Filter) (expression.ConditionBuilder, error) {

	// Iterate records provided
	var firstTime bool = true
	var err error = nil
	var filterExpr expression.ConditionBuilder
	for _, i := range filters {

		// Sanity check
		if i.Field == "" {
			err = errors.New("Field name must be provided")
			return filterExpr, err
		}
		if i.Operator == "" {
			err = errors.New("Operator must be provided")
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
			err = fmt.Errorf("Unsupported operator type: %s", i.Operator)
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

// NewProjectionExpression handles creation of a projection expression
func NewProjectionExpression(fields []Field) (expression.ProjectionBuilder, error) {

	// Setup
	var firstTime bool = true
	var err error = nil
	var projExpr expression.ProjectionBuilder
	if len(fields) == 0 {
		err = errors.New("Must provide at least one field")
		return projExpr, err
	}

	// Iterate records provided
	for _, i := range fields {

		// Sanity check
		if i.Name == "" {
			err = errors.New("Field name must be provided")
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
