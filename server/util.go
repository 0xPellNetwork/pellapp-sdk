package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/0xPellNetwork/pelldvs-libs/log"
	cmtcmd "github.com/0xPellNetwork/pelldvs/cmd/pelldvs/commands"
	pelldvscfg "github.com/0xPellNetwork/pelldvs/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/rs/zerolog"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/0xPellNetwork/pellapp-sdk/baseapp"
	"github.com/0xPellNetwork/pellapp-sdk/client/flags"
	"github.com/0xPellNetwork/pellapp-sdk/server/config"
	"github.com/0xPellNetwork/pellapp-sdk/server/types"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// server context
type Context struct {
	Viper  *viper.Viper
	Config *pelldvscfg.Config
	Logger log.Logger
}

func NewDefaultContext() *Context {
	return NewContext(
		viper.New(),
		pelldvscfg.DefaultConfig(),
		log.NewLogger(os.Stdout),
	)
}

func NewContext(v *viper.Viper, config *pelldvscfg.Config, logger log.Logger) *Context {
	return &Context{v, config, logger}
}

func bindFlags(basename string, cmd *cobra.Command, v *viper.Viper) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("bindFlags failed: %v", r)
		}
	}()

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		err = v.BindEnv(f.Name, fmt.Sprintf("%s_%s", basename, strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))))
		if err != nil {
			panic(err)
		}

		err = v.BindPFlag(f.Name, f)
		if err != nil {
			panic(err)
		}

		// Apply the viper config value to the flag when the flag is not set and
		// viper has a value.
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			err = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			if err != nil {
				panic(err)
			}
		}
	})

	return err
}

// InterceptConfigsPreRunHandler is identical to InterceptConfigsAndCreateContext
// except it also sets the server context on the command and the server logger.
func InterceptConfigsPreRunHandler(cmd *cobra.Command, customAppConfigTemplate string, customAppConfig interface{}, pellDVSConfig *pelldvscfg.Config) error {
	serverCtx, err := InterceptConfigsAndCreateContext(cmd, customAppConfigTemplate, customAppConfig, pellDVSConfig)
	if err != nil {
		return err
	}

	// overwrite default server logger
	logger, err := CreateSDKLogger(serverCtx, cmd.OutOrStdout())
	if err != nil {
		return err
	}
	serverCtx.Logger = logger.With(log.ModuleKey, "server")

	serverCtx.Logger.Info("PellDVS configuration",
		"Pell config", fmt.Sprintf("%+v", serverCtx.Config.Pell),
	)

	// set server context
	return SetCmdServerContext(cmd, serverCtx)
}

// InterceptConfigsAndCreateContext performs a pre-run function for the root daemon
// application command. It will create a Viper literal and a default server
// Context. The server PellDVS configuration will either be read and parsed
// or created and saved to disk, where the server Context is updated to reflect
// the PellDVS configuration. It takes custom app config template and config
// settings to create a custom PellDVS configuration. If the custom template
// is empty, it uses default-template provided by the server. The Viper literal
// is used to read and parse the application configuration. Command handlers can
// fetch the server Context to get the PellDVS configuration or to get access
// to Viper.
func InterceptConfigsAndCreateContext(cmd *cobra.Command, customAppConfigTemplate string, customAppConfig interface{}, cmtConfig *pelldvscfg.Config) (*Context, error) {
	serverCtx := NewDefaultContext()

	// Get the executable name and configure the viper instance so that environmental
	// variables are checked based off that name. The underscore character is used
	// as a separator.
	executableName, err := os.Executable()
	if err != nil {
		return nil, err
	}

	basename := path.Base(executableName)

	// configure the viper instance
	if err := serverCtx.Viper.BindPFlags(cmd.Flags()); err != nil {
		return nil, err
	}
	if err := serverCtx.Viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return nil, err
	}

	serverCtx.Viper.SetEnvPrefix(basename)
	serverCtx.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	serverCtx.Viper.AutomaticEnv()

	// intercept configuration files, using both Viper instances separately
	config, err := interceptConfigs(serverCtx.Viper, customAppConfigTemplate, customAppConfig, cmtConfig)
	if err != nil {
		return nil, err
	}

	// return value is a PellDVS configuration object
	serverCtx.Config = config
	if err = bindFlags(basename, cmd, serverCtx.Viper); err != nil {
		return nil, err
	}

	return serverCtx, nil
}

// CreateSDKLogger creates a the default SDK logger.
// It reads the log level and format from the server context.
func CreateSDKLogger(ctx *Context, out io.Writer) (log.Logger, error) {
	var opts []log.Option
	if ctx.Viper.GetString(flags.FlagLogFormat) == flags.OutputFormatJSON {
		opts = append(opts, log.OutputJSONOption())
	}
	opts = append(opts,
		log.ColorOption(!ctx.Viper.GetBool(flags.FlagLogNoColor)),
		// We use PellDVS flag (cmtcli.TraceFlag) for trace logging.
		//log.TraceOption(ctx.Viper.GetBool(FlagTrace)),
	)

	// check and set filter level or keys for the logger if any
	logLvlStr := ctx.Viper.GetString(flags.FlagLogLevel)
	if logLvlStr == "" {
		return log.NewLogger(out, opts...), nil
	}

	logLvl, err := zerolog.ParseLevel(logLvlStr)
	switch {
	case err != nil:
		// If the log level is not a valid zerolog level, then we try to parse it as a key filter.
		filterFunc, err := log.ParseLogLevel(logLvlStr)
		if err != nil {
			return nil, err
		}

		opts = append(opts, log.FilterOption(filterFunc))
	default:
		opts = append(opts, log.LevelOption(logLvl))
	}

	return log.NewLogger(out, opts...), nil
}

// GetServerContextFromCmd returns a Context from a command or an empty Context
// if it has not been set.
func GetServerContextFromCmd(cmd *cobra.Command) *Context {
	if v := cmd.Context().Value(sdktypes.ServerContextKey); v != nil {
		serverCtxPtr := v.(*Context)
		return serverCtxPtr
	}

	return NewDefaultContext()
}

// SetCmdServerContext sets a command's Context value to the provided argument.
// If the context has not been set, set the given context as the default.
func SetCmdServerContext(cmd *cobra.Command, serverCtx *Context) error {
	v := cmd.Context().Value(sdktypes.ServerContextKey)
	if v == nil {
		v = serverCtx
	}

	serverCtxPtr := v.(*Context)
	*serverCtxPtr = *serverCtx

	return nil
}

// interceptConfigs parses and updates a PellDVS configuration file or
// creates a new one and saves it. It also parses and saves the application
// configuration file. The PellDVS configuration file is parsed given a root
// Viper object, whereas the application is parsed with the private package-aware
// viperCfg object.
func interceptConfigs(rootViper *viper.Viper, customAppTemplate string, customConfig interface{}, cmtConfig *pelldvscfg.Config) (*pelldvscfg.Config, error) {
	rootDir := rootViper.GetString(flags.FlagHome)
	configPath := filepath.Join(rootDir, "config")
	cmtCfgFile := filepath.Join(configPath, "config.toml")

	conf := cmtConfig

	switch _, err := os.Stat(cmtCfgFile); {
	case os.IsNotExist(err):
		pelldvscfg.EnsureRoot(rootDir)

		if err = conf.ValidateBasic(); err != nil {
			return nil, fmt.Errorf("error in config file: %w", err)
		}

		defaultCometCfg := pelldvscfg.DefaultConfig()
		if conf.RPC.PprofListenAddress == defaultCometCfg.RPC.PprofListenAddress {
			conf.RPC.PprofListenAddress = "localhost:6060"
		}

		pelldvscfg.WriteConfigFile(cmtCfgFile, conf)

	case err != nil:
		return nil, err

	default:
		rootViper.SetConfigType("toml")
		rootViper.SetConfigName("config")
		rootViper.AddConfigPath(configPath)

		if err := rootViper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read in %s: %w", cmtCfgFile, err)
		}
	}

	// Read into the configuration whatever data the viper instance has for it.
	// This may come from the configuration file above but also any of the other
	// sources viper uses.
	if err := rootViper.Unmarshal(conf); err != nil {
		return nil, err
	}

	conf.SetRoot(rootDir)

	appCfgFilePath := filepath.Join(configPath, "app.toml")
	if _, err := os.Stat(appCfgFilePath); os.IsNotExist(err) {
		if customAppTemplate != "" {
			config.SetConfigTemplate(customAppTemplate)

			if err = rootViper.Unmarshal(&customConfig); err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", appCfgFilePath, err)
			}

			config.WriteConfigFile(appCfgFilePath, customConfig)
		} else {
			appConf, err := config.ParseConfig(rootViper)
			if err != nil {
				return nil, fmt.Errorf("failed to parse %s: %w", appCfgFilePath, err)
			}

			config.WriteConfigFile(appCfgFilePath, appConf)
		}
	}

	rootViper.SetConfigType("toml")
	rootViper.SetConfigName("app")
	rootViper.AddConfigPath(configPath)

	if err := rootViper.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to merge configuration: %w", err)
	}

	return conf, nil
}

// add server commands
func AddCommands(rootCmd *cobra.Command, defaultNodeHome string, appCreator types.AppCreator, addStartFlags types.ModuleInitFlags) {
	pelldvsCmds := &cobra.Command{
		Use:   "dvs",
		Short: "PellDVS subcommands",
	}

	pelldvsCmds.AddCommand(
		cmtcmd.VersionCmd,
		cmtcmd.ShowNodeIDCmd,
	)

	startCmd := StartCmd(appCreator, defaultNodeHome)
	addStartFlags(startCmd)

	rootCmd.AddCommand(
		startCmd,
		pelldvsCmds,
		version.NewVersionCommand(),
	)
}

// AddCommandsWithStartCmdOptions adds server commands with the provided StartCmdOptions.
func AddCommandsWithStartCmdOptions(rootCmd *cobra.Command, defaultNodeHome string, appCreator types.AppCreator, opts StartCmdOptions) {
	startCmd := StartCmdWithOptions(appCreator, defaultNodeHome, opts)

	rootCmd.AddCommand(
		startCmd,
		version.NewVersionCommand(),
	)
}

// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// TODO there must be a better way to get external IP
func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if skipInterface(iface) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ip := addrToIP(addr)
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

// ListenForQuitSignals listens for SIGINT and SIGTERM. When a signal is received,
// the cleanup function is called, indicating the caller can gracefully exit or
// return.
//
// Note, the blocking behavior of this depends on the block argument.
// The caller must ensure the corresponding context derived from the cancelFn is used correctly.
func ListenForQuitSignals(g *errgroup.Group, block bool, cancelFn context.CancelFunc, logger log.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	f := func() {
		sig := <-sigCh
		cancelFn()

		logger.Info("caught signal", "signal", sig.String())
	}

	if block {
		g.Go(func() error {
			f()
			return nil
		})
	} else {
		go f()
	}
}

// GetAppDBBackend gets the backend type to use for the application DBs.
func GetAppDBBackend(opts types.AppOptions) dbm.BackendType {
	rv := cast.ToString(opts.Get("app-db-backend"))
	if len(rv) == 0 {
		rv = cast.ToString(opts.Get("db_backend"))
	}

	// Cosmos SDK has migrated to cosmos-db which does not support all the backends which tm-db supported
	if rv == "cleveldb" || rv == "badgerdb" || rv == "boltdb" {
		panic(fmt.Sprintf("invalid app-db-backend %q, use %q, %q, %q instead",
			rv, dbm.GoLevelDBBackend, dbm.PebbleDBBackend, dbm.RocksDBBackend),
		)
	}

	if len(rv) != 0 {
		return dbm.BackendType(rv)
	}

	return dbm.GoLevelDBBackend
}

func skipInterface(iface net.Interface) bool {
	if iface.Flags&net.FlagUp == 0 {
		return true // interface down
	}

	if iface.Flags&net.FlagLoopback != 0 {
		return true // loopback interface
	}

	return false
}

func addrToIP(addr net.Addr) net.IP {
	var ip net.IP

	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip
}

func openDB(rootDir string, backendType dbm.BackendType) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("application", backendType, dataDir)
}

func openTraceWriter(traceWriterFile string) (w io.WriteCloser, err error) {
	if traceWriterFile == "" {
		return
	}
	return os.OpenFile(
		traceWriterFile,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0o666,
	)
}

// DefaultBaseappOptions returns the default baseapp options provided by the Cosmos SDK
func DefaultBaseappOptions(appOpts types.AppOptions) []func(*baseapp.BaseApp) {
	return []func(*baseapp.BaseApp){}
}
