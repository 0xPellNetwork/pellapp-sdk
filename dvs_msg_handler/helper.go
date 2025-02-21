package dvsservermanager

import (
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"

	result "github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/result_handler"
	"github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/tx"
)

type dvsMsgHelper struct {
	lock sync.RWMutex

	cdc     codec.Codec
	encoder tx.MsgEncoder

	ResponseHandler grpc1.Server
	RequestHandler  grpc1.Server
}

// TODO: use dependency auto inject
var helper = &dvsMsgHelper{
	lock: sync.RWMutex{},
}

func InitDvsMsgHelper(cdc codec.Codec) {
	helper.lock.Lock()
	defer helper.lock.Unlock()

	if helper.cdc == nil {
		helper.cdc = cdc
	}

	if helper.encoder == nil {
		helper.encoder = tx.NewDefaultDecoder(helper.cdc)
	}

	if helper.RequestHandler == nil {
		helper.RequestHandler = NewRequestHandler(helper.encoder, result.NewResultCustomizedMgr())
	}

	if helper.ResponseHandler == nil {
		helper.ResponseHandler = NewResponseHandler(helper.encoder, result.NewResultCustomizedMgr())
	}
}

func GetResponseHandler() grpc1.Server {
	return helper.ResponseHandler
}

func GetRequestHandler() grpc1.Server {
	return helper.RequestHandler
}

func GetResponseHandlerSrc() *ResponseHandler {
	return helper.ResponseHandler.(*ResponseHandler)
}

func GetRequestHandlerSrc() *RequestHandler {
	return helper.RequestHandler.(*RequestHandler)
}

func EncodeMsgs(msgs ...sdk.Msg) ([]byte, error) {
	return helper.encoder.EncodeMsgs(msgs...)
}

func DecodeMsg(data []byte) (sdk.Msg, error) {
	tx, err := helper.encoder.Decode(data)
	if err != nil {
		return nil, err
	}
	if len(tx.GetMsgs()) == 0 {
		return nil, fmt.Errorf("DecodeMsg invalid tx")
	}

	return tx.GetMsgs()[0], nil
}
