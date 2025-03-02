package service

import (
	cosmosrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"

	"github.com/0xPellNetwork/pellapp-sdk/service/result"
	"github.com/0xPellNetwork/pellapp-sdk/service/tx"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// RequestHandler manages gRPC requests and message routing
// It's responsible for dispatching incoming requests to appropriate handlers and managing results
type RequestHandler struct {
	Mgr           *MsgRouterMgr               // Message router manager that routes messages to corresponding handlers
	ResultHandler *result.CustomResultManager // Processes results and generates output data with digest values
}

// NewRequestHandler creates a new RequestHandler instance implementing the cosmosrpc.Server interface
func NewRequestHandler(encoder tx.MsgEncoder, resultHandler *result.CustomResultManager) cosmosrpc.Server {
	return &RequestHandler{
		Mgr: NewMsgRouterMgr(
			encoder,
			resultHandler,
		),
		ResultHandler: resultHandler,
	}
}

// RegisterService registers a gRPC service to the router manager
func (p *RequestHandler) RegisterService(sd *grpc.ServiceDesc, handler any) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}

// InvokeByMsgData invokes the router handler with raw byte data
func (p *RequestHandler) InvokeByMsgData(ctx sdktypes.Context, data []byte) (*result.Result, error) {
	return p.Mgr.HandleByData(ctx, data)
}

// RegisterResultMsgExtractor registers a custom handler for a specific message type
func (p *RequestHandler) RegisterResultMsgExtractor(msg proto.Message, handler result.ResultMsgExtractor) {
	p.ResultHandler.RegisterCustomizedFunc(msg, handler)
}
