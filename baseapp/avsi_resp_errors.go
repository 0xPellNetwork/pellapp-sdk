package baseapp

import (
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"

	errorsmod "github.com/0xPellNetwork/pellapp-sdk/errors"
)

func responseProcessDVSRequestWithEvents(err error, events []avsitypes.Event, debug bool) *avsitypes.ResponseProcessDVSRequest {
	space, code, log := errorsmod.AVSIInfo(err, debug)
	return &avsitypes.ResponseProcessDVSRequest{
		Code:      code,
		Log:       log,
		Events:    events,
		Codespace: space,
	}
}

func responseProcessDVSResponseWithEvents(err error, events []avsitypes.Event, debug bool) *avsitypes.ResponseProcessDVSResponse {
	space, code, log := errorsmod.AVSIInfo(err, debug)
	return &avsitypes.ResponseProcessDVSResponse{
		Code:      code,
		Log:       log,
		Events:    events,
		Codespace: space,
	}
}
