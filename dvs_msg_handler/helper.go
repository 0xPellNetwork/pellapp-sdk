package dvsservermanager

import (
	"fmt"
	"sync"

	result "github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/result_handler"
	"github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/tx"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
)

type dvsMsgHelper struct {
	lock sync.RWMutex

	cdc     codec.Codec
	encoder tx.MsgEncoder

	PostProcessRequestHandler grpc1.Server
	ProcessRequestHandler     grpc1.Server
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

	if helper.ProcessRequestHandler == nil {
		helper.ProcessRequestHandler = NewProcessRequestHandler(helper.encoder, result.NewResultCustomizedMgr())
	}
	if helper.PostProcessRequestHandler == nil {
		helper.PostProcessRequestHandler = NewPostProcessRequestHandler(helper.encoder, result.NewResultCustomizedMgr())
	}
}

func GetPostProcessRequestHandler() grpc1.Server {
	return helper.PostProcessRequestHandler
}

func GetProcessRequestHandler() grpc1.Server {
	return helper.ProcessRequestHandler
}

func GetPostProcessRequestHandlerSrc() *PostProcessRequestHandler {
	return helper.PostProcessRequestHandler.(*PostProcessRequestHandler)
}

func GetProcessRequestHandlerSrc() *ProcessRequestHandler {
	return helper.ProcessRequestHandler.(*ProcessRequestHandler)
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
