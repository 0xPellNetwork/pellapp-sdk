package types

import (
	proto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
)

// ConfiguratorInterface defines functionality for service configuration and message handling
type Configurator interface {
	// RegisterService registers a gRPC service to the router
	RegisterService(sd *grpc.ServiceDesc, handler any)

	// InvokeByMsgData invokes the router handler with raw byte data
	InvokeByMsgData(ctx Context, data []byte) (*AvsiResult, error)

	// RegisterResultMsgExtractor registers a custom handler for a specific message type
	RegisterResultMsgExtractor(msg proto.Message, handler ResultMsgExtractor)
}
