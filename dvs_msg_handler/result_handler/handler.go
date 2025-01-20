package resulthandler

import (
	sdktypes "github.com/pelldvs/pellapp-sdk/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gogoproto/proto"
)

type ResultCustomizedMgr struct {
	customizedResMap map[string]ResultCustomizedIFace
}

func NewResultCustomizedMgr() *ResultCustomizedMgr {
	return &ResultCustomizedMgr{
		customizedResMap: make(map[string]ResultCustomizedIFace),
	}
}

func (r *ResultCustomizedMgr) RegisterCustomizedFunc(t proto.Message, f ResultCustomizedIFace) {
	r.customizedResMap[sdk.MsgTypeURL(t)] = f
}

// WrapServiceResult wraps a result from a protobuf RPC service method call (res proto.Message, err error)
// in a Result object or error. This method takes care of marshaling the res param to
// protobuf and attaching any events on the ctx.EventManager() to the Result.
func (r *ResultCustomizedMgr) WrapServiceResult(ctx sdktypes.Context, res proto.Message, err error) (*Result, error) {
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

	outResult := &Result{
		Result: &sdktypes.Result{
			Data:         data,
			Events:       ctx.EventManager().AVSIEvents(),
			MsgResponses: []*codectypes.Any{any},
		},
	}

	if resHandler, ok := r.customizedResMap[sdk.MsgTypeURL(res)]; ok {
		outResult.CustomData, _ = resHandler.GetData(res)
		outResult.CustomDigest, _ = resHandler.GetDigest(res)
	}

	return outResult, nil
}
