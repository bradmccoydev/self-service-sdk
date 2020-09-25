package dynamodb_test

import (
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/bradmccoydev/self-service-sdk/internal"
	"github.com/bradmccoydev/self-service-sdk/sdk/aws/dynamodb"
)

// CreateTestTableItem
func CreateTestTableItem() (string, error) {

	// Ensure the test table exists
	err := CreateTableIfNotExists(TestTableConf)
	if err != nil {
		log.Fatal(err)
	}

	// Create an item
	sess := internal.CreateAwsSession(true)
	key := "test-" + time.Now().Format(time.RFC3339)
	item := TestTableFullItem{Name: key, Description: "Blah blah blah"}
	err = dynamodb.CreateItem(sess, TestTableConf.TableName, item)
	if err != nil {
		return "", err
	}
	return key, nil
}

// Test CreateItem
func TestCreateItem(t *testing.T) {

	// Setup backend
	createerr := CreateTableIfNotExists(TestTableConf)
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup input test data
	itemKey := "test-" + time.Now().Format(time.RFC3339)
	emptyInput := TestTableFullItem{}
	noKey := TestTableFullItem{Description: "Fred"}
	emptyKey := TestTableFullItem{Name: "", Description: "Fred"}
	validInput := TestTableFullItem{Name: itemKey, Description: "Something"}

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		input     TestTableFullItem
		expectErr bool
	}{
		{"No inputs", false, "", emptyInput, true},
		{"Just session", true, "", emptyInput, true},
		{"Session & invalid table name", true, TestTableNameInvalid, emptyInput, true},
		{"Session & valid table name", true, TestTableNameValid, emptyInput, true},
		{"Session, valid table name & no key", true, TestTableNameValid, noKey, true},
		{"Session, valid table name & empty key", true, TestTableNameValid, emptyKey, true},
		{"Valid input", true, TestTableNameValid, validInput, false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			var sess *session.Session
			if test.validSess {
				sess = internal.CreateAwsSession(true)
			} else {
				sess = internal.CreateAwsSession(false)
			}
			err := dynamodb.CreateItem(sess, test.tableName, test.input)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test DeleteItem
func TestDeleteItem(t *testing.T) {

	// Setup backend
	itemKey, createerr := CreateTestTableItem()
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup input test data
	emptyInput := TestTableKeys{}
	emptyKey := TestTableKeys{Name: ""}
	validInput := TestTableKeys{Name: itemKey}

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		input     TestTableKeys
		expectErr bool
	}{
		{"No inputs", false, "", emptyInput, true},
		{"Just session", true, "", emptyInput, true},
		{"Session & invalid table name", true, TestTableNameInvalid, emptyInput, true},
		{"Session & valid table name", true, TestTableNameValid, emptyInput, true},
		{"Session, valid table name & empty key", true, TestTableNameValid, emptyKey, true},
		{"Valid input", true, TestTableNameValid, validInput, false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			var sess *session.Session
			if test.validSess {
				sess = internal.CreateAwsSession(true)
			} else {
				sess = internal.CreateAwsSession(false)
			}
			err := dynamodb.DeleteItem(sess, test.tableName, test.input)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test QueryItems
func TestQueryItems(t *testing.T) {

	// Setup backend
	itemKey, createerr := CreateTestTableItem()
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup response array
	var response *[]TestTableFullItem

	// Setup key condition test data
	var emptyKey []dynamodb.Condition
	invalidKeyField := []dynamodb.Condition{{Field: TestTableKeyFieldInvalid, Operator: "EQ", Value: "123"}}
	invalidKeyVal := []dynamodb.Condition{{Field: TestTableKeyFieldValid, Operator: "EQ", Value: "123"}}
	validKey := []dynamodb.Condition{{Field: TestTableKeyFieldValid, Operator: "EQ", Value: itemKey}}

	// Setup expression test data
	var emptyExpression expression.Expression
	emptyKeyExpr, _ := dynamodb.NewExpression(emptyKey, nil, nil)
	invalidKeyFieldExpr, _ := dynamodb.NewExpression(invalidKeyField, nil, nil)
	invalidKeyValueExpr, _ := dynamodb.NewExpression(invalidKeyVal, nil, nil)
	validKeyExpr, _ := dynamodb.NewExpression(validKey, nil, nil)

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		expr      expression.Expression
		expectErr bool
	}{
		{"No inputs", false, "", emptyExpression, true},
		{"Just session", true, "", emptyExpression, true},
		{"Session & invalid table name", true, TestTableNameInvalid, emptyExpression, true},
		{"Session & valid table name", true, TestTableNameValid, emptyKeyExpr, true},
		{"Session, valid table name & empty key", true, TestTableNameValid, emptyExpression, true},
		{"Session, valid table name & invalid key field", true, TestTableNameValid, invalidKeyFieldExpr, true},
		{"Session, valid table name & invalid key value", true, TestTableNameValid, invalidKeyValueExpr, false},
		{"Session, valid table name & valid key", true, TestTableNameValid, validKeyExpr, false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			var sess *session.Session
			if test.validSess {
				sess = internal.CreateAwsSession(true)
			} else {
				sess = internal.CreateAwsSession(false)
			}
			err := dynamodb.QueryItems(sess, test.tableName, test.expr, &response)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test ScanItems
func TestScanItems(t *testing.T) {

	// Setup backend
	itemKey, createerr := CreateTestTableItem()
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup response array
	var response *[]TestTableFullItem

	// Setup filter test data
	invalidFilter := []dynamodb.Condition{{Field: TestTableKeyFieldInvalid, Operator: "EQ", Value: "123"}}
	validFilter := []dynamodb.Condition{{Field: TestTableKeyFieldValid, Operator: "EQ", Value: itemKey}}

	// Setup expression test data
	var emptyExpression expression.Expression
	invalidFilterExpr, _ := dynamodb.NewExpression(nil, invalidFilter, nil)
	validFilterExpr, _ := dynamodb.NewExpression(nil, validFilter, nil)

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		expr      expression.Expression
		expectErr bool
	}{
		{"No inputs", false, "", emptyExpression, true},
		{"Just session", true, "", emptyExpression, true},
		{"Session & invalid table name", true, TestTableNameInvalid, emptyExpression, true},
		{"Session & valid table name", true, TestTableNameValid, emptyExpression, false},
		{"Session, valid table name & invalid filter", true, TestTableNameValid, invalidFilterExpr, false},
		{"Session, valid table name & valid filter", true, TestTableNameValid, validFilterExpr, false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			var sess *session.Session
			if test.validSess {
				sess = internal.CreateAwsSession(true)
			} else {
				sess = internal.CreateAwsSession(false)
			}
			err := dynamodb.ScanItems(sess, test.tableName, test.expr, &response)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}

// Test UpdateItem
func TestUpdateItem(t *testing.T) {

	// Setup backend
	itemKey, createerr := CreateTestTableItem()
	if createerr != nil {
		log.Fatal(createerr)
	}

	// Setup keys test data
	noKey := TestTableKeys{}
	emptyKey := TestTableKeys{Name: ""}
	validKey := TestTableKeys{Name: itemKey}

	// Setup input test data
	emptyInput := TestTableUpdateItem{}
	validInput := TestTableUpdateItem{Description: "Updated", Addtion: "Added"}

	// Setup test data
	tests := []struct {
		desc      string
		validSess bool
		tableName string
		keys      TestTableKeys
		input     TestTableUpdateItem
		expectErr bool
	}{
		{"No inputs", false, "", noKey, emptyInput, true},
		{"Just session", true, "", noKey, emptyInput, true},
		{"Session & invalid table name", true, TestTableNameInvalid, noKey, emptyInput, true},
		{"Session & valid table name", true, TestTableNameValid, noKey, emptyInput, true},
		{"Session, valid table name & no key", true, TestTableNameValid, noKey, emptyInput, true},
		{"Session, valid table name & empty key", true, TestTableNameValid, emptyKey, emptyInput, true},
		{"Session, valid table name & empty input", true, TestTableNameValid, validKey, emptyInput, false},
		{"Valid input", true, TestTableNameValid, validKey, validInput, false},
	}

	// Iterate through the test data
	for _, test := range tests {

		t.Run(test.desc, func(t *testing.T) {

			// Run the test
			var sess *session.Session
			if test.validSess {
				sess = internal.CreateAwsSession(true)
			} else {
				sess = internal.CreateAwsSession(false)
			}
			err := dynamodb.UpdateItem(sess, test.tableName, test.keys, test.input)
			if test.expectErr {
				internal.HasError(t, err)
			} else {
				internal.NoError(t, err)
			}
		})
	}
}
