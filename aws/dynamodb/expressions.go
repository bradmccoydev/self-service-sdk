// This file contains all the bits & pieces related to
// handling of Dynamo DB expressions (filters, keys,
// projections)

package dynamodb

import (
	"strings"

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

// NewExpression - This function creates a new query expression object
//
//   Parameters:
//     keys: an array of key condition(s)
//     filters: an array of filter condition(s)
//     projs: an array of field(s)
//
//   Example:
//     expr, err := NewExpression(keys, filters, projs)
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
			return filterExpr, newErrorFilterExpressionFieldNameNotProvided()
		}
		if i.Operator == "" {
			return filterExpr, newErrorFilterExpressionOperatorNotProvided()
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
			return filterExpr, newErrorFilterExpressionOperatorNotSupported(i.Operator)
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
			return keyExpr, newErrorKeyExpressionFieldNameNotProvided()
		}
		if i.Operator == "" {
			return keyExpr, newErrorKeyExpressionOperatorNotProvided()
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
			return keyExpr, newErrorKeyExpressionOperatorNotSupported(i.Operator)
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
		return projExpr, newErrorProjExpressionFieldsNotProvided()
	}

	// Iterate records provided
	for _, i := range fields {

		// Sanity check
		if i.Name == "" {
			return projExpr, newErrorProjExpressionFieldNameNotProvided()
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
