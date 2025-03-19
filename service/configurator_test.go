package service

import (
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/0xPellNetwork/pellapp-sdk/service/result"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// MockMsgEncoderForConfigurator implements tx.MsgEncoder for testing
type MockMsgEncoderForConfigurator struct {
	mock.Mock
}

func (m *MockMsgEncoderForConfigurator) Decode(txBytes []byte) (types.Tx, error) {
	args := m.Called(txBytes)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(types.Tx), args.Error(1)
}

func (m *MockMsgEncoderForConfigurator) Encode(tx types.Tx) ([]byte, error) {
	args := m.Called(tx)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockMsgEncoderForConfigurator) EncodeMsgs(msgs ...types.Msg) ([]byte, error) {
	args := m.Called(msgs)
	return args.Get(0).([]byte), args.Error(1)
}

// MockTxForConfigurator implements types.Tx for testing
type MockTxForConfigurator struct {
	types.Tx
	msgs []types.Msg
}

func (m *MockTxForConfigurator) GetMsgs() []types.Msg {
	return m.msgs
}

// MockProtoMessage implements proto.Message for testing
type MockProtoMessage struct {
	proto.Message
}

func (m *MockProtoMessage) Reset() { *m = MockProtoMessage{} }

func (m *MockProtoMessage) String() string { return "MockProtoMessage" }

func (m *MockProtoMessage) ProtoMessage() {}

// MockResultMsgExtractorForConfigurator implements result.ResultMsgExtractor for testing
type MockResultMsgExtractorForConfigurator struct {
	mock.Mock
}

func (m *MockResultMsgExtractorForConfigurator) GetData(msg proto.Message) ([]byte, error) {
	args := m.Called(msg)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockResultMsgExtractorForConfigurator) GetDigest(msg proto.Message) ([]byte, error) {
	args := m.Called(msg)
	return args.Get(0).([]byte), args.Error(1)
}

// MockServiceHandlerForConfigurator implements a mock service handler for testing
type MockServiceHandlerForConfigurator struct {
	mock.Mock
}

func (m *MockServiceHandlerForConfigurator) TestMethod(ctx context.Context, req *MockMsgForConfigurator) (*MockMsgForConfigurator, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockMsgForConfigurator), args.Error(1)
}

// MockMsgForConfigurator implements types.Msg for testing
type MockMsgForConfigurator struct {
	types.Msg
}

func (m *MockMsgForConfigurator) GetSigners() []types.AccAddress {
	return nil
}

func (m *MockMsgForConfigurator) ValidateBasic() error {
	return nil
}

func TestNewConfigurator(t *testing.T) {
	// Create mock dependencies
	mockEncoder := new(MockMsgEncoderForConfigurator)
	mockResultManager := result.NewCustomResultManager()

	// Create new configurator
	configurator := NewConfigurator(mockEncoder, mockResultManager)

	// Assert configurator is not nil
	assert.NotNil(t, configurator)
}

func TestConfigurator_RegisterService(t *testing.T) {
	// Create mock dependencies
	mockEncoder := new(MockMsgEncoderForConfigurator)
	mockResultManager := result.NewCustomResultManager()
	configurator := NewConfigurator(mockEncoder, mockResultManager)

	// Create a mock service description with a test method
	mockServiceDesc := &grpc.ServiceDesc{
		ServiceName: "TestService",
		HandlerType: (*MockServiceHandlerForConfigurator)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "TestMethod",
				Handler: func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
					return &MockMsgForConfigurator{}, nil
				},
			},
		},
		Streams: []grpc.StreamDesc{},
	}

	// Create a mock handler
	mockHandler := new(MockServiceHandlerForConfigurator)

	// Register the service
	configurator.RegisterService(mockServiceDesc, mockHandler)

	// Verify mock expectations
	mockHandler.AssertExpectations(t)

	// Verify the service was registered (this is an indirect test since we can't directly access the router's internal state)
	// The test passing without panicking indicates successful registration
}

func TestConfigurator_InvokeByMsgData(t *testing.T) {
	// Create mock dependencies
	mockEncoder := new(MockMsgEncoderForConfigurator)
	mockResultManager := result.NewCustomResultManager()
	configurator := NewConfigurator(mockEncoder, mockResultManager)

	// Create test context
	ctx := sdktypes.Context{}

	// Create test data
	testData := []byte("test data")

	// Set up mock expectations
	mockTx := &MockTxForConfigurator{
		msgs: []types.Msg{},
	}
	mockEncoder.On("Decode", testData).Return(mockTx, nil)

	// Test invoking with data
	c, ok := configurator.(*Configurator)
	assert.True(t, ok)
	result, err := c.InvokeByMsgData(ctx, testData)

	// Assert error is not nil (since we haven't registered any handlers)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verify mock expectations
	mockEncoder.AssertExpectations(t)
}

func TestConfigurator_RegisterResultMsgExtractor(t *testing.T) {
	// Create mock dependencies
	mockEncoder := new(MockMsgEncoderForConfigurator)
	mockResultManager := result.NewCustomResultManager()
	configurator := NewConfigurator(mockEncoder, mockResultManager)

	// Create a mock message
	mockMsg := &MockProtoMessage{}

	// Create a mock handler
	mockHandler := new(MockResultMsgExtractorForConfigurator)

	// Register the result message extractor
	c, ok := configurator.(*Configurator)
	assert.True(t, ok)
	c.RegisterResultMsgExtractor(mockMsg, mockHandler)

	// Verify that the handler was registered by checking if it can be retrieved
	// This is an indirect test since we can't directly access the result manager's internal state
	// The test passing without panicking indicates successful registration
}
