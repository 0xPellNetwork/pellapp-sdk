package result

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// CustomResultManager manages custom result handlers for different message types.
// It allows registering specialized handlers for processing specific message types
// and extracting custom data and digests from them.
type CustomResultManager struct {
	customHandlers map[string]sdktypes.ResultMsgExtractor
}

// NewCustomResultManager creates a new instance of CustomResultManager with
// an initialized map of custom handlers.
func NewCustomResultManager() *CustomResultManager {
	return &CustomResultManager{
		customHandlers: make(map[string]sdktypes.ResultMsgExtractor),
	}
}

// RegisterCustomizedFunc registers a custom result handler for a specific message type.
// The message type is determined by its protobuf URL, and the handler will be called
// when processing results of this message type.
func (r *CustomResultManager) RegisterCustomizedFunc(t proto.Message, f sdktypes.ResultMsgExtractor) {
	r.customHandlers[sdk.MsgTypeURL(t)] = f
}

// WrapServiceResult wraps a result from a protobuf RPC service method call (res proto.Message, err error)
// in a Result object or error. This method takes care of marshaling the res param to
// protobuf and attaching any events on the ctx.EventManager() to the Result.
// If a custom handler is registered for the message type, it will also extract
// custom data and digest from the result.
func (r *CustomResultManager) WrapServiceResult(ctx sdktypes.Context, res proto.Message, err error) (*sdktypes.AvsiResult, error) {
	if err != nil {
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(res)
	if err != nil {
		return nil, err
	}

	var data []byte
	if res != nil {
		data, err = proto.Marshal(res)
		if err != nil {
			return nil, err
		}
	}

	outResult := &sdktypes.AvsiResult{
		Result: &sdktypes.Result{
			Data:         data,
			Events:       ctx.EventManager().AVSIEvents(),
			MsgResponses: []*codectypes.Any{any},
		},
	}

	if resHandler, ok := r.customHandlers[sdk.MsgTypeURL(res)]; ok {
		outResult.CustomData, _ = resHandler.GetData(res)
		outResult.CustomDigest, _ = resHandler.GetDigest(res)
	}

	return outResult, nil
}
