package types

import (
	"context"

	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
)

// AVSI is an interface that enables any finite, deterministic state machine
// to be driven by a blockchain-based replication engine via the AVSI.
type AVSI interface {
	Info(context.Context, *avsitypes.RequestInfo) (*avsitypes.ResponseInfo, error)    // Return application info
	Query(context.Context, *avsitypes.RequestQuery) (*avsitypes.ResponseQuery, error) // Query for state

	//dvs connection
	ProcessDVSRequest(context.Context, *avsitypes.RequestProcessDVSRequest) (*avsitypes.ResponseProcessDVSRequest, error)
	ProcessDVSResponse(context.Context, *avsitypes.RequestProcessDVSResponse) (*avsitypes.ResponseProcessDVSResponse, error)
}
