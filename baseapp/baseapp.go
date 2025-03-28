package baseapp

import (
	"fmt"

	cosmoslog "cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	"github.com/0xPellNetwork/pellapp-sdk/service"
	"github.com/0xPellNetwork/pellapp-sdk/types"
)

type (
	// StoreLoader defines a customizable function to control how we load the
	// CommitMultiStore from disk. This is useful for state migration, when
	// loading a datastore written with an older version of the software. In
	// particular, if a module changed the substore key name (or removed a substore)
	// between two versions of the software.
	StoreLoader func(ms storetypes.CommitMultiStore) error
)

// BaseApp is the main application structure that serves as the foundation
// for dvs applications built on the PellApp-sdk. It manages core
// functionality like message handling, logging, and event indexing.
type BaseApp struct {
	name    string // Name of the application
	version string // Version of the application

	logger log.Logger
	// trace set will return full stack traces for errors in AVSI Log field
	trace bool

	db          dbm.DB                      // common DB backend
	cms         storetypes.CommitMultiStore // Main (uncached) state
	qms         storetypes.MultiStore       // Optional alternative multistore for querying only.
	storeLoader StoreLoader                 // function to handle store loading, may be overridden with SetStoreLoader()

	// flag for sealing options and parameters to a BaseApp
	sealed bool
	// indexEvents defines the set of events in the form {eventType}.{attributeKey},
	// which informs CometBFT what to index. If empty, all events will be indexed.
	indexEvents map[string]struct{}
	// handlers for DVS services
	msgRouter       *service.MsgRouter
	grpcQueryRouter *GRPCQueryRouter // router for redirecting gRPC query calls

	anteHandler types.AnteHandler
}

// NewBaseApp creates and initializes a new BaseApp instance with the provided parameters.
// It sets up the core components needed for the application to function properly.
func NewBaseApp(
	name string,
	logger log.Logger,
	db dbm.DB,
	cdc codec.Codec,
	opts ...func(*BaseApp),
) *BaseApp {
	clogger := cosmoslog.NewCustomLogger(*(logger.Impl().(*zerolog.Logger)))
	app := &BaseApp{
		name:        name,
		logger:      logger,
		msgRouter:   service.NewMsgRouter(cdc),
		cms:         store.NewCommitMultiStore(db, clogger, storemetrics.NewNoOpMetrics()), // by default we use a no-op metric gather in store
		storeLoader: DefaultStoreLoader,
	}

	// apply options
	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (app *BaseApp) GetMsgRouter() *service.MsgRouter {
	return app.msgRouter
}

func (app *BaseApp) GRPCQueryRouter() *GRPCQueryRouter { return app.grpcQueryRouter }

func (app *BaseApp) SetAnteHandler(ah types.AnteHandler) {
	if app.sealed {
		panic("Cannot call SetAnteHandler: baseapp already sealed")
	}

	app.anteHandler = ah
}

func (app *BaseApp) Sealed() {
	if app.sealed {
		panic("Cannot call SetAnteHandler: baseapp already sealed")
	}

	app.sealed = true
}

// SetQueryMultiStore set a alternative MultiStore implementation to support grpc query service.
func (app *BaseApp) SetQueryMultiStore(ms storetypes.MultiStore) {
	app.qms = ms
}

// SetupDefaultQueryStore sets the default query store as a cached version of the main store
func (app *BaseApp) SetupDefaultQueryStore() error {
	if app.cms == nil {
		return fmt.Errorf("commit multistore not initialized")
	}

	app.SetQueryMultiStore(app.cms.CacheMultiStore())
	return nil
}

// MountStore mounts a store to the provided key in the BaseApp multistore,
// using the default DB.
func (app *BaseApp) MountStore(key storetypes.StoreKey, typ storetypes.StoreType) {
	app.cms.MountStoreWithDB(key, typ, app.db)
}

// LastCommitID returns the last CommitID of the multistore.
func (app *BaseApp) LastCommitID() storetypes.CommitID {
	return app.cms.LastCommitID()
}

// LastBlockHeight returns the last committed block height.
func (app *BaseApp) LastBlockHeight() int64 {
	return app.cms.LatestVersion()
}

// DefaultStoreLoader will be used by default and loads the latest version
func DefaultStoreLoader(ms storetypes.CommitMultiStore) error {
	return ms.LoadLatestVersion()
}

// CommitMultiStore returns the root multi-store.
// App constructor can use this to access the `cms`.
// UNSAFE: must not be used during the avsi life cycle.
func (app *BaseApp) CommitMultiStore() storetypes.CommitMultiStore {
	return app.cms
}

// HasQueryMultiStore returns if the Query MultiStore was set.
func (app *BaseApp) HasQueryMultiStore() bool {
	return app.qms != nil
}

// QueryMultiStore returns the QueryMultiStore for GRPC query services
func (app *BaseApp) QueryMultiStore() storetypes.MultiStore {
	if app.HasQueryMultiStore() {
		return app.qms
	}
	return app.cms
}

// SetCMS sets the CommitMultiStore for the BaseApp.
func (app *BaseApp) SetCMS(cms storetypes.CommitMultiStore) {
	if app.sealed {
		panic("SetCMS() on sealed BaseApp")
	}
	app.cms = cms
}
