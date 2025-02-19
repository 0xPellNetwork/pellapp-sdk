package dvsservermanager

import (
	result "github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/result_handler"
	"github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/tx"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"

	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
)

type PostProcessRequestHandler struct {
	Mgr           *MsgRouterMgr
	ResultHandler *result.ResultCustomizedMgr
}

func NewPostProcessRequestHandler(encoder tx.MsgEncoder, resultHandler *result.ResultCustomizedMgr) grpc1.Server {
	return &PostProcessRequestHandler{
		Mgr: NewMsgRouterMgr(
			encoder,
			resultHandler,
		),
		ResultHandler: resultHandler,
	}
}

func (p *PostProcessRequestHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}

// InvokeRouterRawByData
// requestData: binary data from processRequestData, for found router and dispatcher
// reqMsg: post-process-response data, attached to context
func (p *PostProcessRequestHandler) InvokeRouterRawByData(ctx sdktypes.Context, requestData []byte) (*result.Result, error) {
	return p.Mgr.HandleByData(ctx, requestData)
}

func (p *PostProcessRequestHandler) RegisterResultHandler(msg proto.Message, handler result.ResultCustomizedIFace) {
	p.ResultHandler.RegisterCustomizedFunc(msg, handler)
}
