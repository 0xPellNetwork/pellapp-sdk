package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino" // Import amino.proto file for reflection
	"google.golang.org/grpc"

	"github.com/0xPellNetwork/pellapp-sdk/client"
	"github.com/0xPellNetwork/pellapp-sdk/server/config"
	"github.com/0xPellNetwork/pellapp-sdk/server/types"
)

// NewGRPCServer returns a correctly configured and initialized gRPC server.
// Note, the caller is responsible for starting the server. See StartGRPCServer.
func NewGRPCServer(clientCtx client.Context, app types.Application, cfg config.GRPCConfig) (*grpc.Server, error) {
	maxSendMsgSize := cfg.MaxSendMsgSize
	if maxSendMsgSize == 0 {
		maxSendMsgSize = config.DefaultGRPCMaxSendMsgSize
	}

	maxRecvMsgSize := cfg.MaxRecvMsgSize
	if maxRecvMsgSize == 0 {
		maxRecvMsgSize = config.DefaultGRPCMaxRecvMsgSize
	}

	grpcSrv := grpc.NewServer(
		grpc.ForceServerCodec(codec.NewProtoCodec(clientCtx.InterfaceRegistry).GRPCCodec()),
		grpc.MaxSendMsgSize(maxSendMsgSize),
		grpc.MaxRecvMsgSize(maxRecvMsgSize),
	)

	app.RegisterGRPCServer(grpcSrv)
	return grpcSrv, nil
}

// StartGRPCServer starts the provided gRPC server on the address specified in cfg.
//
// Note, this creates a blocking process if the server is started successfully.
// Otherwise, an error is returned. The caller is expected to provide a Context
// that is properly canceled or closed to indicate the server should be stopped.
func StartGRPCServer(ctx context.Context, logger log.Logger, cfg config.GRPCConfig, grpcSrv *grpc.Server) error {
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to listen on address %s: %w", cfg.Address, err)
	}

	errCh := make(chan error)

	// Start the gRPC in an external goroutine as Serve is blocking and will return
	// an error upon failure, which we'll send on the error channel that will be
	// consumed by the for block below.
	go func() {
		logger.Info("starting gRPC server...", "address", cfg.Address)
		errCh <- grpcSrv.Serve(listener)
	}()

	// Start a blocking select to wait for an indication to stop the server or that
	// the server failed to start properly.
	select {
	case <-ctx.Done():
		// The calling process canceled or closed the provided context, so we must
		// gracefully stop the gRPC server.
		logger.Info("stopping gRPC server...", "address", cfg.Address)
		grpcSrv.GracefulStop()

		return nil

	case err := <-errCh:
		logger.Error("failed to start gRPC server", "err", err)
		return err
	}
}
