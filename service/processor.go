package service

import (
	cosmosrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"

	"github.com/0xPellNetwork/pellapp-sdk/service/result"
	"github.com/0xPellNetwork/pellapp-sdk/service/tx"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// Processor manages gRPC requests and message routing
// It's responsible for dispatching incoming requests to appropriate handlers and managing results
type Processor struct {
	Router           *MsgRouterMgr               // Message router manager that routes messages to corresponding handlers
	ResultManager *result.CustomResultManager // Processes results and generates output data with digest values
}

// NewProcessor creates a new RequestHandler instance implementing the cosmosrpc.Server interface
func NewProcessor(encoder tx.MsgEncoder, resultHandler *result.CustomResultManager) cosmosrpc.Server {
	return &Processor{
		Router: NewMsgRouterMgr(
			encoder,
			resultHandler,
		),
		ResultManager: resultHandler,
	}
}

// RegisterService registers a gRPC service to the router manager
func (p *Processor) RegisterService(sd *grpc.ServiceDesc, handler any) {
	RegisterServiceRouter(p.Router, sd, handler)
}

// InvokeByMsgData invokes the router handler with raw byte data
func (p *Processor) InvokeByMsgData(ctx sdktypes.Context, data []byte) (*result.Result, error) {
	return p.Router.HandleByData(ctx, data)
}

// RegisterResultMsgExtractor registers a custom handler for a specific message type
func (p *Processor) RegisterResultMsgExtractor(msg proto.Message, handler result.ResultMsgExtractor) {
	p.ResultManager.RegisterCustomizedFunc(msg, handler)
}
