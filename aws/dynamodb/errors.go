// This file contains all the bits & pieces related to
// error messages for this package.

package dynamodb

import (
	"errors"
	"fmt"
)

/***
Expression errors
***/

func newErrorExpressionFieldNameNotProvided() error {
	return errors.New("A field name must be provided for the expression")
}

func newErrorExpressionKeyNotProvided() error {
	return errors.New("A key condition must be provided in the expression")
}

func newErrorExpressionOperatorNotProvided() error {
	return errors.New("An operator must be provided for the expression")
}

func newErrorExpressionFilterOperatorNotSupported(op string) error {
	return fmt.Errorf("The operator type %s is not supported by filter expressions", op)
}

func newErrorExpressionKeyCondOperatorNotSupported(op string) error {
	return fmt.Errorf("The operator type %s is not supported by key condition expressions", op)
}

func newErrorExpressionProjFieldsNotProvided() error {
	return fmt.Errorf("At least one field must be provided for a projection expression")
}

/***
Table errors
***/

func newErrorBillingModeNotProvided() error {
	return errors.New("Billing mode must be provided")
}

func newErrorTableAttributeNameNotProvided() error {
	return errors.New("Table attribute must have a name")
}

func newErrorTableAttributeTypeNotProvided() error {
	return errors.New("Table attribute must have a type")
}

func newErrorTableAttributesNotProvided() error {
	return errors.New("Table attributes must be provided")
}

func newErrorTableKeyAttributesNotProvided() error {
	return errors.New("No key attributes were provided")
}

func newErrorTableKeyFieldKeyTypeNotProvided(keyName string) error {
	return fmt.Errorf("The key field %s did not include a key type", keyName)
}

func newErrorTableDetailsNotProvided() error {
	return errors.New("Table details were not returned")
}

func newErrorTableNameNotProvided() error {
	return errors.New("Table name must be provided")
}
