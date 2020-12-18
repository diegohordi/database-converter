package errors

import (
	"fmt"
	"log"
)

type ApplicationError struct {
	error error
	message string
	code int
}

func BuildApplicationError(err error, message string, code int) *ApplicationError{
	return &ApplicationError{
		error:   err,
		message: message,
		code:    code,
	}
}

func (err *ApplicationError) GetError() error {
	return err.error
}

func (err *ApplicationError) GetMessage() string {
	return err.message
}

func (err *ApplicationError) GetCode() int {
	return err.code
}

func (err *ApplicationError) ThrowPanic() {
	log.Fatal(fmt.Errorf(err.GetMessage(), ": %s \n", err.GetError()))
}