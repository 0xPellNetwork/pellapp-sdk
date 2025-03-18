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

// MsgRouter handles message encoding/decoding and routing for DVS services
type MsgRouter struct {
	cdc          codec.Codec
	encoder      tx.MsgEncoder
	configurator cosmosrpc.Server
}

// Init initializes the DvsMsgHandlers with codec and creates default handlers if not set
func NewMsgRouter(cdc codec.Codec) *MsgRouter {
	encoder := tx.NewDefaultDecoder(cdc)

	return &MsgRouter{
		cdc:          cdc,
		encoder:      encoder,
		configurator: NewConfigurator(encoder, result.NewCustomResultManager()),
	}
}

// InvokeByMsgData routes raw byte data to the configurator
func (h *MsgRouter) InvokeByMsgData(sdkCtx sdktypes.Context, data []byte) (*sdktypes.AvsiResult, error) {
	return h.configurator.(*Configurator).InvokeByMsgData(sdkCtx, data)
}

// GetConfigurator returns the configurator
func (h *MsgRouter) GetConfigurator() *Configurator {
	return h.configurator.(*Configurator)
}

// EncodeMsgs encodes SDK messages into bytes using the configured encoder
func (h *MsgRouter) EncodeMsgs(msgs ...sdk.Msg) ([]byte, error) {
	return h.encoder.EncodeMsgs(msgs...)
}

// DecodeMsg decodes bytes into an SDK message using the configured encoder
func (h *MsgRouter) DecodeMsg(data []byte) (sdk.Msg, error) {
	tx, err := h.encoder.Decode(data)
	if err != nil {
		return nil, err
	}
	if len(tx.GetMsgs()) == 0 {
		return nil, fmt.Errorf("DecodeMsg invalid tx")
	}

	return tx.GetMsgs()[0], nil
}
