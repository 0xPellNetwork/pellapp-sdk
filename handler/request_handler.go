package dvsservermanager

import (
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"

	result "github.com/0xPellNetwork/pellapp-sdk/handler/result"
	"github.com/0xPellNetwork/pellapp-sdk/handler/tx"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// RequestHandler manages gRPC requests and message routing
// It's responsible for dispatching incoming requests to appropriate handlers and managing results
type RequestHandler struct {
	Mgr           *MsgRouterMgr               // Message router manager that routes messages to corresponding handlers
	ResultHandler *result.CustomResultManager // Processes results and generates output data with digest values
}

// NewRequestHandler creates a new RequestHandler instance implementing the grpc1.Server interface
func NewRequestHandler(encoder tx.MsgEncoder, resultHandler *result.CustomResultManager) grpc1.Server {
	return &RequestHandler{
		Mgr: NewMsgRouterMgr(
			encoder,
			resultHandler,
		),
		ResultHandler: resultHandler,
	}
}

// RegisterService registers a gRPC service to the router manager
func (p *RequestHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}

// InvokeRouterRawByData invokes the router handler with raw byte data
func (p *RequestHandler) InvokeRouterRawByData(ctx sdktypes.Context, data []byte) (*result.Result, error) {
	res, err := p.Mgr.HandleByData(ctx, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// RegisterResultHandler registers a custom handler for a specific message type
func (p *RequestHandler) RegisterResultHandler(msg proto.Message, handler result.CustomResultHandler) {
	p.ResultHandler.RegisterCustomizedFunc(msg, handler)
}
