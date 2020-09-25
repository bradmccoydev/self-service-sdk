// This file contains all the bits & pieces related to
// error messages for the secretsmanager package.

package secretsmanager

import (
	"errors"
)

func newErrorSecretNameNotProvided() error {
	return errors.New("A secret name must be provided")
}

func newErrorSecretMapNotProvided() error {
	return errors.New("At least one secret key/value pair must be provided")
}

func newErrorSecretStringNotProvided() error {
	return errors.New("A secret string must be provided")
}

func newErrorSecretBinaryAndStringNotProvided() error {
	return errors.New("Either a secret string or binary must be provided")
}
