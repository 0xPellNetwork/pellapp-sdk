package service

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosrpc "github.com/cosmos/gogoproto/grpc"

	"github.com/0xPellNetwork/pellapp-sdk/service/result"
	"github.com/0xPellNetwork/pellapp-sdk/service/tx"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// DvsMsgHandlers handles message encoding/decoding and routing for DVS services
type DvsMsgHandlers struct {
	cdc       codec.Codec
	encoder   tx.MsgEncoder
	processor cosmosrpc.Server
}

// Init initializes the DvsMsgHandlers with codec and creates default handlers if not set
func NewDvsMsgHandlers(cdc codec.Codec) *DvsMsgHandlers {
	encoder := tx.NewDefaultDecoder(cdc)

	return &DvsMsgHandlers{
		cdc:       cdc,
		encoder:   encoder,
		processor: NewProcessor(encoder, result.NewCustomResultManager()),
	}
}

// InvokeByMsgData routes raw byte data to the processor
func (h *DvsMsgHandlers) InvokeByMsgData(sdkCtx sdktypes.Context, data []byte) (*result.Result, error) {
	return h.processor.(*Processor).InvokeByMsgData(sdkCtx, data)
}

// EncodeMsgs encodes SDK messages into bytes using the configured encoder
func (h *DvsMsgHandlers) EncodeMsgs(msgs ...sdk.Msg) ([]byte, error) {
	return h.encoder.EncodeMsgs(msgs...)
}

// DecodeMsg decodes bytes into an SDK message using the configured encoder
func (h *DvsMsgHandlers) DecodeMsg(data []byte) (sdk.Msg, error) {
	tx, err := h.encoder.Decode(data)
	if err != nil {
		return nil, err
	}
	if len(tx.GetMsgs()) == 0 {
		return nil, fmt.Errorf("DecodeMsg invalid tx")
	}

	return tx.GetMsgs()[0], nil
}
