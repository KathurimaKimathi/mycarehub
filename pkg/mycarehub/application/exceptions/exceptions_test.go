package exceptions_test

import (
	"fmt"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"

	"github.com/stretchr/testify/assert"
)

func TestCustomErrors(t *testing.T) {
	err := exceptions.UserNotFoundError(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.NormalizeMSISDNError(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.PinNotFoundError(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.PinMismatchError(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.InvalidFlavourDefinedErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.SaveUserPinError(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.EmptyUserIDErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.ProfileNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.InvalidatePinErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.ResetPinErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.PINExpiredErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.EmptyInputErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.PINErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.NotOptedInErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.NotActiveErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.InvalidContactTypeErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.NoContactsErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.ContactNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.GenerateTempPINErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.InvalidResetPinPayloadErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = exceptions.ExpiredPinErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.LoginCountUpdateErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.LoginTimeUpdateErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.NexAllowedLOginTimeErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.SendSMSErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.FailedToUpdateItemErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.ItemNotFoundErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = exceptions.InputValidationErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.EncryptionErr(fmt.Errorf("error"))
	assert.NotNil(t, err)

	err = exceptions.FailedToSaveItemErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.GeneratePinErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
	err = exceptions.ValidatePINDigitsErr(fmt.Errorf("error"))
	assert.NotNil(t, err)
}
