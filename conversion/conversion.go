package conversion

import "database-converter/errors"

/**
Starts the conversion.
*/
func Start() *errors.ApplicationError{
	if err:= testConnections(); err != nil {
		return err
	}
	return nil
}
