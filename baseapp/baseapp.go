package baseapp

import (
	"github.com/0xPellNetwork/pelldvs-libs/log"
)

type BaseApp struct {
	logger log.Logger

	// trace set will return full stack traces for errors in ABCI Log field
	trace bool

	// indexEvents defines the set of events in the form {eventType}.{attributeKey},
	// which informs CometBFT what to index. If empty, all events will be indexed.
	indexEvents map[string]struct{}
}

func NewBaseApp(
	logger log.Logger,
) *BaseApp {
	app := &BaseApp{
		logger: logger,
	}
	return app
}

// Trace returns the boolean value for logging error stack traces.
func (app *BaseApp) Trace() bool {
	return app.trace
}

func (app *BaseApp) SetIndexEvents(ie []string) {
	app.indexEvents = make(map[string]struct{}, len(ie))
	for _, e := range ie {
		app.indexEvents[e] = struct{}{}
	}
}

func (app *BaseApp) SetTrace(trace bool) {
	app.trace = trace
}
