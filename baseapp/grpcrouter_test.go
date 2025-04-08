package baseapp

import (
	"context"
	"testing"

	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/runtime/protoiface"

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// MockService is a simple gRPC service for testing
type MockService struct{}

// MockRequest is a simple request for testing
type MockRequest struct{}

// MockResponse is a simple response for testing
type MockResponse struct{}

// MockMethodHandler is a simple method handler for testing
func MockMethodHandler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	return &MockResponse{}, nil
}

func TestNewGRPCQueryRouter(t *testing.T) {
	router := NewGRPCQueryRouter()

	assert.NotNil(t, router)
	assert.NotNil(t, router.routes)
	assert.NotNil(t, router.hybridHandlers)
	assert.Empty(t, router.routes)
	assert.Empty(t, router.hybridHandlers)
}

func TestGRPCQueryRouter_Route(t *testing.T) {
	router := NewGRPCQueryRouter()

	// Create a mock handler function with a specific behavior we can verify
	expectedValue := []byte("test response")
	mockHandler := func(ctx sdktypes.Context, req *avsitypes.RequestQuery) (*avsitypes.ResponseQuery, error) {
		return &avsitypes.ResponseQuery{Value: expectedValue}, nil
	}

	// Test route not found
	handler := router.Route("/not/exists")
	assert.Nil(t, handler)

	// Test adding and finding a route
	fqName := "/mock.Service/MockMethod"
	router.routes[fqName] = mockHandler

	// Get the handler and verify it exists
	handler = router.Route(fqName)
	assert.NotNil(t, handler)

	// Execute the handler and verify it behaves as expected
	// This tests functionality instead of comparing function pointers
	ctx := sdktypes.Context{}
	req := &avsitypes.RequestQuery{}
	resp, err := handler(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedValue, resp.Value)
}

func TestGRPCQueryRouter_HybridHandlerByRequestName(t *testing.T) {
	router := NewGRPCQueryRouter()

	// Manually add a hybrid handler for testing
	testName := "test.Request"
	testHandler := func(ctx context.Context, req, resp protoiface.MessageV1) error {
		return nil
	}
	router.hybridHandlers[testName] = []func(ctx context.Context, req, resp protoiface.MessageV1) error{testHandler}

	// Test retrieving existing handler
	handlers := router.HybridHandlerByRequestName(testName)
	assert.NotNil(t, handlers)
	assert.Len(t, handlers, 1)

	// Test retrieving non-existent handler
	handlers = router.HybridHandlerByRequestName("non.existent")
	assert.Empty(t, handlers)
}

func TestGRPCQueryRouter_SetInterfaceRegistry(t *testing.T) {
	router := NewGRPCQueryRouter()

	// Initially, the codec should be nil
	assert.Nil(t, router.cdc)
	assert.Nil(t, router.binaryCodec)

	// Set interface registry
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	router.SetInterfaceRegistry(interfaceRegistry)

	// After setting, codec should be non-nil
	assert.NotNil(t, router.cdc)
	assert.NotNil(t, router.binaryCodec)
	assert.IsType(t, codec.NewProtoCodec(interfaceRegistry).GRPCCodec(), router.cdc)
	assert.IsType(t, codec.NewProtoCodec(interfaceRegistry), router.binaryCodec)

	// Verify service data contains the reflection service
	require.GreaterOrEqual(t, len(router.serviceData), 1)
}
