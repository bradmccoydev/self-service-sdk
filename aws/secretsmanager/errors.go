// This file contains all the bits & pieces related to
// error messages for the secretsmanager package.

package secretsmanager

import (
	"errors"
)

func newErrorSecretNameNotProvided() error {
	return errors.New("A secret name must be provided")
}
