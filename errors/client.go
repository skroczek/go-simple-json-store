package errors

import "errors"

var ErrorMissingExtension = errors.New("missing extension")
var ErrorInvalidPath = errors.New("invalid path")

func IsClientError(err error) bool {
	return err == ErrorMissingExtension || err == ErrorInvalidPath
}
