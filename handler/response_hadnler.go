package handler

import (
	cosmosrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"

	"github.com/0xPellNetwork/pellapp-sdk/handler/result"
	"github.com/0xPellNetwork/pellapp-sdk/handler/tx"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// ResponseHandler manages gRPC service routing and custom result handling.
// It encapsulates a message router manager and a result handler to process
// incoming requests and generate appropriate responses.
type ResponseHandler struct {
	Mgr           *MsgRouterMgr               // Manages message routing and dispatching
	ResultHandler *result.CustomResultManager // Processes results and generates output data with digest values
}

// NewResponseHandler creates a new ResponseHandler instance that implements the cosmosrpc.Server interface.
// It initializes the message router manager with the provided encoder and result handler.
func NewResponseHandler(encoder tx.MsgEncoder, resultHandler *result.CustomResultManager) cosmosrpc.Server {
	return &ResponseHandler{
		Mgr: NewMsgRouterMgr(
			encoder,
			resultHandler,
		),
		ResultHandler: resultHandler,
	}
}

// RegisterService registers a gRPC service with the ResponseHandler.
// This allows the handler to route incoming requests to the appropriate service implementation.
func (p *ResponseHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}

// InvokeRouterRawByData processes raw binary request data and routes it to the appropriate handler.
// requestData: binary data from processRequestData, for found router and dispatcher
// reqMsg: post-process-response data, attached to context
func (p *ResponseHandler) InvokeRouterRawByData(ctx sdktypes.Context, requestData []byte) (*result.Result, error) {
	return p.Mgr.HandleByData(ctx, requestData)
}

// RegisterResultHandler registers a custom handler for processing results of a specific message type.
func (p *ResponseHandler) RegisterResultHandler(msg proto.Message, handler result.CustomResultHandler) {
	p.ResultHandler.RegisterCustomizedFunc(msg, handler)
}
