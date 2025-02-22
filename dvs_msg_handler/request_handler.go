package dvsservermanager

import (
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"

	result "github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/result_handler"
	"github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/tx"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

type RequestHandler struct {
	Mgr           *MsgRouterMgr
	ResultHandler *result.ResultCustomizedMgr
}

func NewRequestHandler(encoder tx.MsgEncoder, resultHandler *result.ResultCustomizedMgr) grpc1.Server {
	return &RequestHandler{
		Mgr: NewMsgRouterMgr(
			encoder,
			resultHandler,
		),
		ResultHandler: resultHandler,
	}
}

func (p *RequestHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}

func (p *RequestHandler) InvokeRouterRawByData(ctx sdktypes.Context, data []byte) (*result.Result, error) {
	res, err := p.Mgr.HandleByData(ctx, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *RequestHandler) RegisterResultHandler(msg proto.Message, handler result.ResultCustomizedIFace) {
	p.ResultHandler.RegisterCustomizedFunc(msg, handler)
}
