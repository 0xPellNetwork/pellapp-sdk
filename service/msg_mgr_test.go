package service

import (
	"context"
	"fmt"
	"testing"

	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	protov2 "google.golang.org/protobuf/proto"

	"github.com/0xPellNetwork/pellapp-sdk/proto/test"
	"github.com/0xPellNetwork/pellapp-sdk/service/result"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// MockMsgEncoderForMsgMgr implements tx.MsgEncoder for testing
type MockMsgEncoderForMsgMgr struct {
	DecodeFunc     func([]byte) (types.Tx, error)
	EncodeFunc     func(types.Tx) ([]byte, error)
	EncodeMsgsFunc func(...types.Msg) ([]byte, error)
}

func (m *MockMsgEncoderForMsgMgr) Decode(data []byte) (types.Tx, error) {
	if m.DecodeFunc != nil {
		return m.DecodeFunc(data)
	}
	return nil, nil
}

func (m *MockMsgEncoderForMsgMgr) Encode(tx types.Tx) ([]byte, error) {
	if m.EncodeFunc != nil {
		return m.EncodeFunc(tx)
	}
	return nil, nil
}

func (m *MockMsgEncoderForMsgMgr) EncodeMsgs(msgs ...types.Msg) ([]byte, error) {
	if m.EncodeMsgsFunc != nil {
		return m.EncodeMsgsFunc(msgs...)
	}
	return nil, nil
}

// MockMsg implements sdk.Msg for testing
type MockMsg struct {
	*test.TestMsg
}

func NewMockMsg(typeURL string) *MockMsg {
	return &MockMsg{
		TestMsg: &test.TestMsg{
			TypeUrl: typeURL,
		},
	}
}

func (m *MockMsg) Reset() { m.TestMsg.Reset() }

func (m *MockMsg) String() string { return m.TypeUrl }

func (m *MockMsg) ProtoMessage() {}

// MockService implements a test gRPC service
type MockService struct {
	HandlerFunc func(context.Context, any) (any, error)
}

func (s *MockService) TestMethod(ctx context.Context, req any) (any, error) {
	if s.HandlerFunc != nil {
		return s.HandlerFunc(ctx, req)
	}
	return nil, nil
}

func TestNewMsgRouterMgr(t *testing.T) {
	encoder := &MockMsgEncoderForMsgMgr{}
	resultHandler := result.NewCustomResultManager()

	router := NewMsgRouterMgr(encoder, resultHandler)
	assert.NotNil(t, router)
	assert.NotNil(t, router.Router)
	assert.Equal(t, encoder, router.encoder)
	assert.Equal(t, resultHandler, router.resultHandler)
}

func TestRegisterMsgHandler(t *testing.T) {
	encoder := &MockMsgEncoderForMsgMgr{}
	resultHandler := result.NewCustomResultManager()
	router := NewMsgRouterMgr(encoder, resultHandler)

	// Create a mock service description
	sd := &grpc.ServiceDesc{
		ServiceName: "test.service",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "TestMethod",
				Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
					msg := NewMockMsg("/test.service/TestMethod")
					if err := dec(msg); err != nil {
						return nil, err
					}
					return msg, nil
				},
			},
		},
	}

	mockService := &MockService{}
	err := router.RegisterMsgHandler(sd, sd.Methods[0], mockService)
	require.NoError(t, err)

	// Verify handler was registered
	msg := NewMockMsg("/test.service/TestMethod")
	handler, found := router.GetHandler(sdktypes.Context{}, msg)
	assert.True(t, found)
	assert.NotNil(t, handler)
}

func TestGetHandlerByData(t *testing.T) {
	encoder := &MockMsgEncoderForMsgMgr{
		DecodeFunc: func(data []byte) (types.Tx, error) {
			return &MockTxForMsgMgr{
				Msgs: []types.Msg{NewMockMsg("/test.service/TestMethod")},
			}, nil
		},
	}
	resultHandler := result.NewCustomResultManager()
	router := NewMsgRouterMgr(encoder, resultHandler)

	// Register a handler
	sd := &grpc.ServiceDesc{
		ServiceName: "test.service",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "TestMethod",
				Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
					msg := NewMockMsg("/test.service/TestMethod")
					if err := dec(msg); err != nil {
						return nil, err
					}
					return msg, nil
				},
			},
		},
	}

	mockService := &MockService{}
	err := router.RegisterMsgHandler(sd, sd.Methods[0], mockService)
	require.NoError(t, err)

	// Test getting handler by data
	handler, err := router.GetHandlerByData([]byte("test data"))
	require.NoError(t, err)
	assert.NotNil(t, handler)
}

// MockTxForMsgMgr implements types.Tx for testing
type MockTxForMsgMgr struct {
	Msgs []types.Msg
}

func (m *MockTxForMsgMgr) GetMsgs() []types.Msg {
	return m.Msgs
}

func (m *MockTxForMsgMgr) GetMsgsV2() ([]protov2.Message, error) {
	msgs := make([]protov2.Message, len(m.Msgs))
	for i, msg := range m.Msgs {
		msgs[i] = msg.(protov2.Message)
	}
	return msgs, nil
}

func TestHandleByData(t *testing.T) {
	encoder := &MockMsgEncoderForMsgMgr{
		DecodeFunc: func(data []byte) (types.Tx, error) {
			return &MockTxForMsgMgr{
				Msgs: []types.Msg{NewMockMsg("/test.service/TestMethod")},
			}, nil
		},
	}
	resultHandler := result.NewCustomResultManager()
	router := NewMsgRouterMgr(encoder, resultHandler)

	// Register a handler
	sd := &grpc.ServiceDesc{
		ServiceName: "test.service",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "TestMethod",
				Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
					msg := NewMockMsg("/test.service/TestMethod")
					if err := dec(msg); err != nil {
						return nil, err
					}
					// Return a valid proto.Message
					return NewMockMsg("/test.service/TestMethod"), nil
				},
			},
		},
	}

	mockService := &MockService{}
	err := router.RegisterMsgHandler(sd, sd.Methods[0], mockService)
	require.NoError(t, err)

	// Create a context with necessary values
	ctx := sdktypes.NewContext(context.Background())

	// Test handling data
	result, err := router.HandleByData(ctx, []byte("test data"))
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestDVSResponseHandler(t *testing.T) {
	encoder := &MockMsgEncoderForMsgMgr{}
	resultHandler := result.NewCustomResultManager()
	router := NewMsgRouterMgr(encoder, resultHandler)

	// Create a mock service description with both request and response handlers
	sd := &grpc.ServiceDesc{
		ServiceName: "test.service",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "TestMethod",
				Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
					msg := NewMockMsg("/test.service/TestMethod")
					if err := dec(msg); err != nil {
						return nil, err
					}
					return msg, nil
				},
			},
			{
				MethodName: "TestMethodDVSResponsHandler",
				Handler: func(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
					msg := NewMockMsg("/test.service/TestMethod")
					if err := dec(msg); err != nil {
						return nil, err
					}
					return msg, nil
				},
			},
		},
	}

	mockService := &MockService{}
	// Register both handlers
	err := router.RegisterMsgHandler(sd, sd.Methods[0], mockService)
	require.NoError(t, err)
	err = router.RegisterMsgHandler(sd, sd.Methods[1], mockService)
	require.NoError(t, err)

	// Test request handler
	msg := NewMockMsg("/test.service/TestMethod")
	ctx := sdktypes.NewContext(context.Background())
	handler, found := router.GetHandler(ctx, msg)
	assert.True(t, found)
	assert.NotNil(t, handler)

	// Test response handler
	ctx = sdktypes.NewContext(context.Background())
	validatedData := &avsitypes.DVSResponse{
		Data: []byte("test data"),
	}
	ctx = ctx.WithValidatedResponse(validatedData)
	handler, found = router.GetHandler(ctx, msg)
	assert.True(t, found)
	assert.NotNil(t, handler)

	// Verify that the handlers are different by checking the router map directly
	requestKey := types.MsgTypeURL(msg)
	responseKey := fmt.Sprintf(DVSResponseHandlerKeyPattern, requestKey)

	requestHandler, requestFound := router.Router[requestKey]
	responseHandler, responseFound := router.Router[responseKey]

	assert.True(t, requestFound)
	assert.True(t, responseFound)
	assert.NotNil(t, requestHandler)
	assert.NotNil(t, responseHandler)
}
