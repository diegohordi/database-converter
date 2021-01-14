package errors

import "net/http"

// Application error builder
type ApplicationErrorBuilder struct {
	applicationError *ApplicationError
}

// Build up an application error.
func (builder *ApplicationErrorBuilder) Build() *ApplicationError {
	return builder.applicationError
}

// Creates an error just with the given message and the 500 code as default.
func WithMessageBuilder(message string) *ApplicationErrorBuilder {
	return &ApplicationErrorBuilder{&ApplicationError{message: message, code: http.StatusInternalServerError}}
}

// Creates an error just with the given message and code.
func WithMessageAndCodeBuilder(message string, code int) *ApplicationErrorBuilder {
	return &ApplicationErrorBuilder{&ApplicationError{message: message, code: code}}
}

// Creates an error just with the given source error.
func WithSourceErrorBuilder(sourceError error) *ApplicationErrorBuilder {
	return &ApplicationErrorBuilder{&ApplicationError{message: sourceError.Error(), sourceError: sourceError, code: http.StatusInternalServerError}}
}

// Creates an error just with a custom message and the given source error.
func WithMessageAndSourceErrorBuilder(message string, sourceError error) *ApplicationErrorBuilder {
	return &ApplicationErrorBuilder{&ApplicationError{message: message, sourceError: sourceError, code: http.StatusInternalServerError}}
}


