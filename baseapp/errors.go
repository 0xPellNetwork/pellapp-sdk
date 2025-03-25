package baseapp

import (
	errorsmod "cosmossdk.io/errors"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
)

// queryResult returns a ResponseQuery from an error. It will try to parse AVSI
// info from the error.
func queryResult(err error, debug bool) *avsitypes.ResponseQuery {
	space, code, log := errorsmod.ABCIInfo(err, debug)
	return &avsitypes.ResponseQuery{
		Codespace: space,
		Code:      code,
		Log:       log,
	}
}
