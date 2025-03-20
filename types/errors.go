package types

import (
	sdkerrors "github.com/0xPellNetwork/pellapp-sdk/errors"

	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
)

// AvsiBaseError contains common fields for AVSI error responses
type AvsiBaseError struct {
	Code      uint32
	Log       string
	Events    []avsitypes.Event
	Codespace string
}

// WarpAvsiBaseError creates a standardized error structure from an error and result
func WarpAvsiBaseError(err error, res *AvsiResult, debug bool) *AvsiBaseError {
	// Extract error details
	space, code, log := sdkerrors.AVSIInfo(err, debug)

	// Handle potentially nil result
	var events []avsitypes.Event
	if res != nil {
		events = res.Events
	}

	return &AvsiBaseError{
		Code:      code,
		Log:       log,
		Events:    events,
		Codespace: space,
	}
}
