package errors

import (
	"encoding/json"
)

// Common struct for the Application errors
type ApplicationError struct {
	sourceError error
	message     string
	code        int
}

// Gets the source error from this application error.
func (err *ApplicationError) SourceError() error {
	return err.sourceError
}

// Gets the message from this application error.
func (err *ApplicationError) Message() string {
	return err.message
}

// Gets the code from this application error.
func (err *ApplicationError) Code() int {
	return err.code
}

func (err *ApplicationError) Error() string {
	return err.message
}

func (err *ApplicationError) MarshalJSON() ([]byte, error) {
	value, jsonErr := json.Marshal(struct {
		Code        int         `json:"code"`
		Message     string      `json:"message"`
		SourceError interface{} `json:"source"`
	}{
		Code:        err.code,
		Message:     err.message,
		SourceError: err.sourceError,
	})
	if jsonErr != nil {
		return nil, jsonErr
	}
	return value, nil
}
