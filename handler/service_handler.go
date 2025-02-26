package handler

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosrpc "github.com/cosmos/gogoproto/grpc"

	result "github.com/0xPellNetwork/pellapp-sdk/handler/result"
	"github.com/0xPellNetwork/pellapp-sdk/handler/tx"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// DvsMsgHandlers handles message encoding/decoding and routing for DVS services
type DvsMsgHandlers struct {
	cdc      codec.Codec
	encoder  tx.MsgEncoder
	handlers Handlers
}

// Handlers contains GRPC servers for handling requests and responses
type Handlers struct {
	ResponseHandler cosmosrpc.Server
	RequestHandler  cosmosrpc.Server
}

// Init initializes the DvsMsgHandlers with codec and creates default handlers if not set
func NewDvsMsgHandlers(cdc codec.Codec) *DvsMsgHandlers {
	encoder := tx.NewDefaultDecoder(cdc)

	return &DvsMsgHandlers{
		cdc:     cdc,
		encoder: encoder,
		handlers: Handlers{
			ResponseHandler: NewResponseHandler(encoder, result.NewCustomResultManager()),
			RequestHandler:  NewRequestHandler(encoder, result.NewCustomResultManager()),
		},
	}
}

// ResponseHandlerInvokeRouterRawByData routes raw byte data to the response handler
func (h *DvsMsgHandlers) ResponseHandlerInvokeRouterRawByData(sdkCtx sdktypes.Context, data []byte) (*result.Result, error) {
	return h.handlers.ResponseHandler.(*ResponseHandler).InvokeRouterRawByData(sdkCtx, data)
}

// RequestHandlerInvokeRouterRawByData routes raw byte data to the request handler
func (h *DvsMsgHandlers) RequestHandlerInvokeRouterRawByData(sdkCtx sdktypes.Context, data []byte) (*result.Result, error) {
	return h.handlers.RequestHandler.(*RequestHandler).InvokeRouterRawByData(sdkCtx, data)
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
