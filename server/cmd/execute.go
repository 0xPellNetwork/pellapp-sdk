package cmd

import (
	"context"

	cmtcli "github.com/0xPellNetwork/pelldvs/libs/cli"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/0xPellNetwork/pellapp-sdk/client"
	"github.com/0xPellNetwork/pellapp-sdk/client/flags"
	"github.com/0xPellNetwork/pellapp-sdk/server"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// Execute executes the root command of an application. It handles creating a
// server context object with the appropriate server and client objects injected
// into the underlying stdlib Context. It also handles adding core CLI flags,
// specifically the logging flags. It returns an error upon execution failure.
func Execute(rootCmd *cobra.Command, envPrefix, defaultHome string) error {
	// Create and set a client.Context on the command's Context. During the pre-run
	// of the root command, a default initialized client.Context is provided to
	// seed child command execution with values such as AccountRetriever, Keyring,
	// and a PellDVS RPC. This requires the use of a pointer reference when
	// getting and setting the client.Context. Ideally, we utilize
	// https://github.com/spf13/cobra/pull/1118.
	ctx := CreateExecuteContext(context.Background())

	rootCmd.PersistentFlags().String(flags.FlagLogLevel, zerolog.InfoLevel.String(), "The logging level (trace|debug|info|warn|error|fatal|panic|disabled or '*:<level>,<key>:<level>')")
	// NOTE: The default logger is only checking for the "json" value, any other value will default to plain text.
	rootCmd.PersistentFlags().String(flags.FlagLogFormat, "plain", "The logging format (json|plain)")
	rootCmd.PersistentFlags().Bool(flags.FlagLogNoColor, false, "Disable colored logs")

	executor := cmtcli.PrepareBaseCmd(rootCmd, envPrefix, defaultHome)
	return executor.ExecuteContext(ctx)
}

// CreateExecuteContext returns a base Context with server and client context
// values initialized.
func CreateExecuteContext(ctx context.Context) context.Context {
	srvCtx := server.NewDefaultContext()
	ctx = context.WithValue(ctx, sdktypes.ClientContextKey, &client.Context{})
	ctx = context.WithValue(ctx, sdktypes.ServerContextKey, srvCtx)

	return ctx
}
