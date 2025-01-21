package dvsservermanager

import (
	"context"
	"fmt"

	result "github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/result_handler"
	"github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler/tx"

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
)

type MsgHandler func(ctx sdktypes.Context, msg sdk.Msg) (*result.Result, error)

// MsgRouterMgr defines router for dvs server
type MsgRouterMgr struct {
	Router                 map[string]MsgHandler
	encoder                tx.MsgEncoder
	findRouterTypeNameFunc func(msg sdk.Msg) string // ONLY FOR router dispatcher; register use sdk.MsgTypeURL
	resultHandler          *result.ResultCustomizedMgr
}

func NewMsgRouterMgr(
	encoder tx.MsgEncoder,
	resultHandler *result.ResultCustomizedMgr,
) *MsgRouterMgr {
	return &MsgRouterMgr{
		Router:  map[string]MsgHandler{},
		encoder: encoder,
		findRouterTypeNameFunc: func(msg sdk.Msg) string {
			return sdk.MsgTypeURL(msg)
		},
		resultHandler: resultHandler,
	}
}

// RegisterMsgHandler
// inspire by github.com/cosmos/cosmos-sdk@v0.50.9/baseapp/msg_service_router.go:120 MsgServiceRouter.registerMsgServiceHandler
func (m *MsgRouterMgr) RegisterMsgHandler(sd *grpc.ServiceDesc, method grpc.MethodDesc, handler interface{}) error {
	fqMethod := fmt.Sprintf("/%s/%s", sd.ServiceName, method.MethodName)
	methodHandler := method.Handler

	var requestTypeName string

	// NOTE: This is how we pull the concrete request type for each handler for registering in the InterfaceRegistry.
	// This approach is maybe a bit hacky, but less hacky than reflecting on the handler object itself.
	// We use a no-op interceptor to avoid actually calling into the handler itself.
	_, _ = methodHandler(nil, context.Background(), func(i interface{}) error {
		msg, ok := i.(sdk.Msg)
		if !ok {
			// We panic here because there is no other alternative and the app cannot be initialized correctly
			// this should only happen if there is a problem with code generation in which case the app won't
			// work correctly anyway.
			panic(fmt.Errorf("unable to register service method %s: %T does not implement sdk.Msg", fqMethod, i))
		}

		requestTypeName = sdk.MsgTypeURL(msg)
		return nil
	}, noopInterceptor)

	// requestTypeName register check
	if _, ok := m.Router[requestTypeName]; !ok {
		m.Router[requestTypeName] = func(ctx sdktypes.Context, msg sdk.Msg) (*result.Result, error) {
			// ctx = ctx.WithEventManager(sdk.NewEventManager())
			interceptor := func(goCtx context.Context, _ interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				goCtx = context.WithValue(goCtx, sdktypes.ContextKey, ctx)
				return handler(goCtx, msg)
			}

			res, err := methodHandler(handler, ctx, noopDecoder, interceptor)
			if err != nil {
				return nil, err
			}

			resMsg, ok := res.(proto.Message)
			if !ok {
				return nil, fmt.Errorf("expecting proto.Message, got %T", resMsg)
			}

			return m.resultHandler.WrapServiceResult(ctx, resMsg, err)
		}
	} else {
		log.Warn("duplicate existing handler for %s", requestTypeName)
	}

	return nil
}

func (m *MsgRouterMgr) GetHandler(msg sdk.Msg) MsgHandler {
	return m.Router[m.findRouterTypeNameFunc(msg)]
}

func (m *MsgRouterMgr) GetHandlerByData(data []byte) MsgHandler {
	msgTx, err := m.encoder.Decode(data)
	if err != nil {
		return nil
	}
	for _, msg := range msgTx.GetMsgs() {
		msgType := m.findRouterTypeNameFunc(msg)
		if handler, ok := m.Router[msgType]; ok {
			return handler
		}
	}

	return nil
}

func (m *MsgRouterMgr) HandleByData(ctx sdktypes.Context, data []byte) (*result.Result, error) {
	msgTx, err := m.encoder.Decode(data)
	if err != nil {
		return nil, err
	}
	for _, msg := range msgTx.GetMsgs() {
		msgType := m.findRouterTypeNameFunc(msg)
		if handler, ok := m.Router[msgType]; ok {
			return handler(ctx, msg)
		}
	}

	return nil, fmt.Errorf("no handler found for %s", msgTx.GetMsgs())
}

func noopDecoder(_ interface{}) error { return nil }

func noopInterceptor(_ context.Context, _ interface{}, _ *grpc.UnaryServerInfo, _ grpc.UnaryHandler) (interface{}, error) {
	return nil, nil
}

// RegisterServiceRouter helper for registering service router
func RegisterServiceRouter(routerMgr *MsgRouterMgr, sd *grpc.ServiceDesc, handler interface{}) {
	for _, method := range sd.Methods {
		err := routerMgr.RegisterMsgHandler(sd, method, handler)
		if err != nil {
			panic(err)
		}
	}
}
