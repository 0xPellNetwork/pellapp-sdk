package client

import (
	stdcontext "context"
	"encoding/json"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
	"sigs.k8s.io/yaml"
)

var _ gogogrpc.ClientConn = Context{}

// Context
type Context struct {
	HomeDir      string
	OutputFormat string
	Output       io.Writer

	GRPCClient        *grpc.ClientConn
	Codec             codec.Codec
	InterfaceRegistry codectypes.InterfaceRegistry

	// CmdContext is the context.Context from the Cobra command.
	CmdContext stdcontext.Context
}

// NewContext creates a new Context
func NewContext() Context {
	return Context{}
}

// WithCodec returns a copy of the Context with an updated Codec.
func (ctx Context) WithCodec(m codec.Codec) Context {
	ctx.Codec = m
	return ctx
}

// WithInterfaceRegistry returns the context with an updated InterfaceRegistry
func (ctx Context) WithInterfaceRegistry(interfaceRegistry codectypes.InterfaceRegistry) Context {
	ctx.InterfaceRegistry = interfaceRegistry
	return ctx
}

// WithCodec sets the Client connection for the context
func (ctx Context) WithGRPCClient(grpcClient *grpc.ClientConn) Context {
	ctx.GRPCClient = grpcClient
	return ctx
}

// WithOutputFormat returns a copy of the context with an updated OutputFormat field.
func (ctx Context) WithOutputFormat(format string) Context {
	ctx.OutputFormat = format
	return ctx
}

// WithOutput returns a copy of the context with an updated output writer (e.g. stdout).
func (ctx Context) WithOutput(w io.Writer) Context {
	ctx.Output = w
	return ctx
}

// WithHomeDir returns a copy of the Context with HomeDir set.
func (ctx Context) WithHomeDir(dir string) Context {
	if dir != "" {
		ctx.HomeDir = dir
	}
	return ctx
}

// WithCmdContext returns a copy of the context with an updated context.Context,
// usually set to the cobra cmd context.
func (ctx Context) WithCmdContext(c stdcontext.Context) Context {
	ctx.CmdContext = c
	return ctx
}

func (ctx Context) Invoke(stdctx stdcontext.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	return ctx.GRPCClient.Invoke(stdctx, method, args, reply, opts...)
}

func (ctx Context) NewStream(stdctx stdcontext.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return ctx.GRPCClient.NewStream(stdctx, desc, method, opts...)
}

// PrintProto outputs toPrint to the ctx.Output based on ctx.OutputFormat which is
// either text or json. If text, toPrint will be YAML encoded. Otherwise, toPrint
// will be JSON encoded using ctx.Codec. An error is returned upon failure.
func (ctx Context) PrintProto(toPrint proto.Message) error {
	// always serialize JSON initially because proto json can't be directly YAML encoded
	out, err := ctx.Codec.MarshalJSON(toPrint)
	if err != nil {
		return err
	}
	return ctx.printOutput(out)
}

// PrintRaw is a variant of PrintProto that doesn't require a proto.Message type
// and uses a raw JSON message. No marshaling is performed.
func (ctx Context) PrintRaw(toPrint json.RawMessage) error {
	return ctx.printOutput(toPrint)
}

func (ctx Context) printOutput(out []byte) error {
	var err error
	if ctx.OutputFormat == "text" {
		out, err = yaml.JSONToYAML(out)
		if err != nil {
			return err
		}
	}

	writer := ctx.Output
	if writer == nil {
		writer = os.Stdout
	}

	_, err = writer.Write(out)
	if err != nil {
		return err
	}

	if ctx.OutputFormat != "text" {
		// append new-line for formats besides YAML
		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
