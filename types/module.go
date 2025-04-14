package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// AppModule is a marker interface for all app modules in the system
type AppModule interface {
	// IsAppModule is a dummy method to tag a struct as implementing an AppModule.
	IsAppModule()
}

// HasServices is the extension interface that modules should implement to register
// implementations of services defined in .proto files.
type BasicModule interface {
	AppModule

	// Name returns the module's name
	Name() string

	// RegisterServices registers module services with the configurator
	RegisterServices(Configurator)

	// RegisterInterfaces registers module interfaces with the registry
	RegisterInterfaces(types.InterfaceRegistry)

	// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the bank module. client.Context
	RegisterGRPCGatewayRoutes(gogogrpc.ClientConn, *runtime.ServeMux)

	// RegisterQueryServices registers module query services with the router
	RegisterQueryServices(router gogogrpc.Server)
}

// MsgResultExtractor interface for modules that need to provide custom result extraction
type MsgResultExtractor interface {
	// RegisterResultMsgExtractor registers result message extractors with the configurator
	RegisterResultMsgExtractors(Configurator)
}

// ModuleManager defines a module manager that provides the high level utility
// for managing and executing operations for a group of modules
type ModuleManager struct {
	Modules map[string]any
}

// NewManager creates a new Manager object based on the provided modules
func NewManager(modules ...BasicModule) *ModuleManager {
	moduleMap := make(map[string]any)

	for _, module := range modules {
		moduleMap[module.Name()] = module
	}

	return &ModuleManager{
		Modules: moduleMap,
	}
}

// RegisterServices calls RegisterServices on all modules
func (m *ModuleManager) RegisterServices(c Configurator) {
	for _, module := range m.Modules {
		module.(BasicModule).RegisterServices(c)
	}
}

// RegisterInterfaces calls RegisterInterfaces on all modules
func (m *ModuleManager) RegisterInterfaces(ir types.InterfaceRegistry) {
	for _, module := range m.Modules {
		module.(BasicModule).RegisterInterfaces(ir)
	}
}

// RegisterGRPCGatewayRoutes registers all module rest routes
func (m *ModuleManager) RegisterGRPCGatewayRoutes(clientCtx gogogrpc.ClientConn, rtr *runtime.ServeMux) {
	for _, module := range m.Modules {
		module.(BasicModule).RegisterGRPCGatewayRoutes(clientCtx, rtr)
	}
}

// RegisterResultMsgExtractors calls RegisterResultMsgExtractor on modules
// that implement the MsgResultExtractor interface
func (m *ModuleManager) RegisterResultMsgExtractors(c Configurator) {
	for _, module := range m.Modules {
		if msgExtractor, ok := module.(MsgResultExtractor); ok {
			msgExtractor.RegisterResultMsgExtractors(c)
		}
	}
}

// RegisterQueryServices calls RegisterQueryServices on all modules
func (m *ModuleManager) RegisterQueryServices(router gogogrpc.Server) {
	for _, module := range m.Modules {
		module.(BasicModule).RegisterQueryServices(router)
	}
}
