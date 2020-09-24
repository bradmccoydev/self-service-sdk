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

func newErrorFilterExpressionFieldNameNotProvided() error {
	return errors.New("A field name must be provided for the filter expression")
}

func newErrorFilterExpressionOperatorNotProvided() error {
	return errors.New("An operator must be provided for the filter expression")
}

func newErrorFilterExpressionOperatorNotSupported(op string) error {
	return fmt.Errorf("The operator type %s is not supported by filter expressions", op)
}

func newErrorKeyExpressionFieldNameNotProvided() error {
	return errors.New("A field name must be provided for the key condition expression")
}

func newErrorKeyExpressionOperatorNotProvided() error {
	return errors.New("An operator must be provided for the key condition expression")
}

func newErrorKeyExpressionOperatorNotSupported(op string) error {
	return fmt.Errorf("The operator type %s is not supported by key condition expressions", op)
}

func newErrorKeyExpressionKeyNotProvided() error {
	return errors.New("A key condition must be provided in the expression")
}

func newErrorProjExpressionFieldNameNotProvided() error {
	return errors.New("A field name must be provided for the projection expression")
}

func newErrorProjExpressionFieldsNotProvided() error {
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
