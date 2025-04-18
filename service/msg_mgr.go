package service

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/log"
	"google.golang.org/grpc"

	"github.com/0xPellNetwork/pellapp-sdk/service/result"
	"github.com/0xPellNetwork/pellapp-sdk/service/tx"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

const DVSResponsHandler = "DVSResponsHandler"

const DVSResponseHandlerKeyPattern = "%sResp"

// MsgHandler is a function type that handles SDK messages and returns a result or error
type MsgHandler func(ctx sdktypes.Context, msg sdk.Msg) (*sdktypes.AvsiResult, error)

// defaultMsgKeyFunc is a default function for getting the message key
func defaultMsgKeyFunc(msg sdk.Msg) string {
	return sdk.MsgTypeURL(msg)
}

// MsgRouterMgr defines router for dvs server
type MsgRouterMgr struct {
	Router        map[string]MsgHandler
	encoder       tx.MsgEncoder
	calcMsgKey    func(msg sdk.Msg) string // ONLY FOR router dispatcher; register use sdk.MsgTypeURL
	resultHandler *result.CustomResultManager
}

// NewMsgRouterMgr creates a new message router manager with the provided encoder and result handler
func NewMsgRouterMgr(
	encoder tx.MsgEncoder,
	resultHandler *result.CustomResultManager,
) *MsgRouterMgr {
	return &MsgRouterMgr{
		Router:        map[string]MsgHandler{},
		encoder:       encoder,
		calcMsgKey:    defaultMsgKeyFunc,
		resultHandler: resultHandler,
	}
}

// RegisterMsgHandler registers a gRPC service method as a message handler
// Inspired by github.com/cosmos/cosmos-sdk@v0.50.9/baseapp/msg_service_router.go:120 MsgServiceRouter.registerMsgServiceHandler
func (m *MsgRouterMgr) RegisterMsgHandler(sd *grpc.ServiceDesc, method grpc.MethodDesc, handler any) error {
	fqMethod := fmt.Sprintf("/%s/%s", sd.ServiceName, method.MethodName)

	var requestTypeName string

	// NOTE: This is how we pull the concrete request type for each handler for registering in the InterfaceRegistry.
	// This approach is maybe a bit hacky, but less hacky than reflecting on the handler object itself.
	// We use a no-op interceptor to avoid actually calling into the handler itself.
	_, _ = method.Handler(nil, context.Background(), func(i any) error {
		msg, ok := i.(sdk.Msg)
		if !ok {
			// We panic here because there is no other alternative and the app cannot be initialized correctly
			// this should only happen if there is a problem with code generation in which case the app won't
			// work correctly anyway.
			panic(fmt.Errorf("unable to register service method %s: %T does not implement sdk.Msg", fqMethod, i))
		}

		requestTypeName = sdk.MsgTypeURL(msg)

		if strings.Contains(method.MethodName, DVSResponsHandler) {
			requestTypeName = fmt.Sprintf(DVSResponseHandlerKeyPattern, requestTypeName)
		}

		return nil
	}, noopInterceptor)

	// requestTypeName register check
	if _, ok := m.Router[requestTypeName]; !ok {
		m.Router[requestTypeName] = func(ctx sdktypes.Context, msg sdk.Msg) (*sdktypes.AvsiResult, error) {
			// ctx = ctx.WithEventManager(sdk.NewEventManager())
			interceptor := func(goCtx context.Context, _ any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
				goCtx = context.WithValue(goCtx, sdktypes.ContextKey, ctx)
				return handler(goCtx, msg)
			}

			res, err := method.Handler(handler, ctx, noopDecoder, interceptor)
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

// GetHandler returns the handler for a specific message type
func (m *MsgRouterMgr) GetHandler(ctx sdktypes.Context, msg sdk.Msg) (MsgHandler, bool) {
	key := m.calcMsgKey(msg)
	if ctx.ValidatedResponse() != nil {
		key = fmt.Sprintf(DVSResponseHandlerKeyPattern, key)
	}

	handler, found := m.Router[key]
	return handler, found
}

// GetHandlerByData decodes the message data and returns the appropriate handler
func (m *MsgRouterMgr) GetHandlerByData(data []byte) (MsgHandler, error) {
	msgTx, err := m.encoder.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	for _, msg := range msgTx.GetMsgs() {
		msgType := m.calcMsgKey(msg)
		if handler, ok := m.Router[msgType]; ok {
			return handler, nil
		}
	}

	return nil, fmt.Errorf("no handler found for message types in transaction")
}

// HandleByData decodes the message data, finds the appropriate handler, and processes the message
func (m *MsgRouterMgr) HandleByData(ctx sdktypes.Context, data []byte) (*sdktypes.AvsiResult, error) {
	msgTx, err := m.encoder.Decode(data)
	if err != nil {
		return nil, err
	}

	for _, msg := range msgTx.GetMsgs() {
		handler, found := m.GetHandler(ctx, msg)
		if found {
			return handler(ctx, msg)
		}
	}

	return nil, fmt.Errorf("no handler found for %s", msgTx.GetMsgs())
}

// noopDecoder is a no-operation decoder used during handler registration
func noopDecoder(_ any) error { return nil }

// noopInterceptor is a no-operation interceptor used during handler registration
func noopInterceptor(_ context.Context, _ any, _ *grpc.UnaryServerInfo, _ grpc.UnaryHandler) (any, error) {
	return nil, nil
}

// RegisterServiceRouter is a helper function for registering all methods of a gRPC service
func RegisterServiceRouter(routerMgr *MsgRouterMgr, sd *grpc.ServiceDesc, handler any) {
	for _, method := range sd.Methods {
		if err := routerMgr.RegisterMsgHandler(sd, method, handler); err != nil {
			panic(err)
		}
	}
}
