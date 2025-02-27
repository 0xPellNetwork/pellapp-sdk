package baseapp

import (
	"context"
	"fmt"

	"github.com/0xPellNetwork/pellapp-sdk/utils/tx"

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"
)

type MsgHandler func(ctx sdktypes.Context, msg sdk.Msg) (*sdktypes.DvsResult, error)

// MsgServiceRouter manages message routing and handling in the application
type MsgServiceRouter struct {
	// routes maps message type URLs to their request handlers
	routes map[string]MsgHandler
	// responseRoutes maps message type URLs to their response handlers
	responseRoutes map[string]sdktypes.DvsResponseHandler
	// customizers maps message type URLs to their result customization handlers
	customizers map[string]sdktypes.ResultCustomizedIFace
	// encoder handles message encoding and decoding
	encoder tx.MsgEncoder
}

// NewMsgServiceRouter creates a new MsgServiceRouter.
func NewMsgServiceRouter(
	encoder tx.MsgEncoder,
) *MsgServiceRouter {
	return &MsgServiceRouter{
		routes:      map[string]MsgHandler{},
		customizers: make(map[string]sdktypes.ResultCustomizedIFace),
		encoder:     encoder,
	}
}

// RegisterMsgHandler
// inspire by github.com/cosmos/cosmos-sdk@v0.50.9/baseapp/msg_service_router.go:120 MsgServiceRouter.registerMsgServiceHandler
func (m *MsgServiceRouter) RegisterMsgHandler(sd *grpc.ServiceDesc, method grpc.MethodDesc, handler interface{}) error {
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
	if _, ok := m.routes[requestTypeName]; !ok {
		m.routes[requestTypeName] = func(ctx sdktypes.Context, msg sdk.Msg) (*sdktypes.DvsResult, error) {
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

			return m.WrapServiceResult(ctx, resMsg, err)
		}
	} else {
		log.Warn("duplicate existing handler for %s", requestTypeName)
	}

	return nil
}

func (r *MsgServiceRouter) RegisterCustomizedFunc(t proto.Message, f sdktypes.ResultCustomizedIFace) {
	r.customizers[sdk.MsgTypeURL(t)] = f
}

func (r *MsgServiceRouter) WrapServiceResult(ctx sdktypes.Context, res proto.Message, err error) (*sdktypes.DvsResult, error) {
	if err != nil {
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(res)
	if err != nil {
		return nil, err
	}

	var data []byte
	if res != nil {
		data, err = proto.Marshal(res)
		if err != nil {
			return nil, err
		}
	}

	outResult := &sdktypes.DvsResult{
		Result: &sdktypes.Result{
			Data:         data,
			Events:       ctx.EventManager().AVSIEvents(),
			MsgResponses: []*codectypes.Any{any},
		},
	}

	if resHandler, ok := r.customizers[sdk.MsgTypeURL(res)]; ok {
		outResult.CustomData, _ = resHandler.GetData(res)
		outResult.CustomDigest, _ = resHandler.GetDigest(res)
	}
	return outResult, nil
}

func (m *MsgServiceRouter) Handler(msg sdk.Msg) MsgHandler {
	return m.routes[sdk.MsgTypeURL(msg)]
}

func (m *MsgServiceRouter) RequestHandleByData(ctx sdktypes.Context, data []byte) (*sdktypes.DvsResult, error) {
	msgTx, err := m.encoder.Decode(data)
	if err != nil {
		return nil, err
	}
	for _, msg := range msgTx.GetMsgs() {
		if handler, ok := m.routes[sdk.MsgTypeURL(msg)]; ok {
			return handler(ctx, msg)
		}
	}
	return nil, fmt.Errorf("no request handler found for %s", msgTx.GetMsgs())
}

// RegisterServiceRouter for registering service router
func RegisterServiceRouter(routerMgr *MsgServiceRouter, sd *grpc.ServiceDesc, handler interface{}) {
	for _, method := range sd.Methods {
		err := routerMgr.RegisterMsgHandler(sd, method, handler)
		if err != nil {
			panic(err)
		}
	}
}

func (r *MsgServiceRouter) RegisterResponseHandlerFunc(t proto.Message, f sdktypes.DvsResponseHandler) {
	r.responseRoutes[sdk.MsgTypeURL(t)] = f
}

func (m *MsgServiceRouter) ResponseHandleByData(ctx sdktypes.Context, data []byte) error {
	msgTx, err := m.encoder.Decode(data)
	if err != nil {
		return err
	}
	for _, msg := range msgTx.GetMsgs() {
		if handler, ok := m.responseRoutes[sdk.MsgTypeURL(msg)]; ok {
			return handler.ResponseHandler(ctx, msg)
		}
	}
	return fmt.Errorf("no response handler found for %s", msgTx.GetMsgs())
}

func noopDecoder(_ interface{}) error { return nil }

func noopInterceptor(_ context.Context, _ interface{}, _ *grpc.UnaryServerInfo, _ grpc.UnaryHandler) (interface{}, error) {
	return nil, nil
}
