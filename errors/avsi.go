package errors

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	// SuccessAVSICode declares an AVSI response use 0 to signal that the
	// processing was successful and no error is returned.
	SuccessAVSICode uint32 = 0

	// All unclassified errors that do not provide an AVSI code are clubbed
	// under an internal error code and a generic message instead of
	// detailed error string.
	internalAVSICodespace        = UndefinedCodespace
	internalAVSICode      uint32 = 1
)

// AVSIInfo returns the AVSI error information as consumed by the tendermint
// client. Returned codespace, code, and log message should be used as a AVSI response.
// Any error that does not provide AVSICode information is categorized as error
// with code 1, codespace UndefinedCodespace
// When not running in a debug mode all messages of errors that do not provide
// AVSICode information are replaced with generic "internal error". Errors
// without an AVSICode information as considered internal.
func AVSIInfo(err error, debug bool) (codespace string, code uint32, log string) {
	if errIsNil(err) {
		return "", SuccessAVSICode, ""
	}

	encode := defaultErrEncoder
	if debug {
		encode = debugErrEncoder
	}

	code, codespace = avsiInfo(err)
	log = encode(err)
	return
}

// The debugErrEncoder encodes the error with a stacktrace.
func debugErrEncoder(err error) string {
	return fmt.Sprintf("%+v", err)
}

func defaultErrEncoder(err error) string {
	return err.Error()
}

// avsiInfo tests if given error contains an AVSI code and returns the value of
// it if available. This function is testing for the causer interface as well
// and unwraps the error.
func avsiInfo(err error) (code uint32, codespace string) {
	if errIsNil(err) {
		return SuccessAVSICode, ""
	}

	var customErr *Error

	if errors.As(err, &customErr) {
		code = customErr.AVSICode()
		codespace = customErr.Codespace()
	} else {
		code = internalAVSICode
		codespace = internalAVSICodespace
	}

	return
}

// errIsNil returns true if value represented by the given error is nil.
//
// Most of the time a simple == check is enough. There is a very narrowed
// spectrum of cases (mostly in tests) where a more sophisticated check is
// required.
func errIsNil(err error) bool {
	if err == nil {
		return true
	}

	if val := reflect.ValueOf(err); val.Kind() == reflect.Ptr {
		return val.IsNil()
	}

	return false
}
