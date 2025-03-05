package baseapp

import (
	"github.com/0xPellNetwork/pelldvs-libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	handler "github.com/0xPellNetwork/pellapp-sdk/service"
)

// BaseApp is the main application structure that serves as the foundation
// for dvs applications built on the PellApp-sdk. It manages core
// functionality like message handling, logging, and event indexing.
type BaseApp struct {
	name    string // Name of the application
	version string // Version of the application

	logger log.Logger
	// trace set will return full stack traces for errors in ABCI Log field
	trace bool

	// indexEvents defines the set of events in the form {eventType}.{attributeKey},
	// which informs CometBFT what to index. If empty, all events will be indexed.
	indexEvents map[string]struct{}
	// handlers for DVS services
	msgRouter *handler.MsgRouter
}

// NewBaseApp creates and initializes a new BaseApp instance with the provided parameters.
// It sets up the core components needed for the application to function properly.
func NewBaseApp(
	name string,
	logger log.Logger,
	cdc codec.Codec,
	opts ...func(*BaseApp),
) *BaseApp {
	app := &BaseApp{
		name:      name,
		logger:    logger,
		msgRouter: handler.NewMsgRouter(cdc),
	}

	// apply options
	for _, opt := range opts {
		opt(app)
	}

	return app
}

func (app *BaseApp) GetMsgRouter() *handler.MsgRouter {
	return app.msgRouter
}
