package types

import (
	"context"
	"time"

	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"

	dvstypes "github.com/0xPellNetwork/pellapp-sdk/pelldvs/types"
)

type ContextKeyType string

const (
	ContextKey ContextKeyType = "pkg_context"
)

type Context struct {
	baseCtx                   context.Context
	ms                        storetypes.MultiStore
	eventManager              EventManagerI
	chainID                   int64
	height                    int64
	groupNumbers              []uint32
	groupThresholdPercentages []uint32
	requestData               []byte
	operators                 []*avsitypes.Operator
	validatedResponse         *dvstypes.RequestPostRequestValidatedData
	logger                    log.Logger
}

// Read-only accessors
func (c Context) Context() context.Context { return c.baseCtx }

func (c Context) EventManager() EventManagerI { return c.eventManager }

func (c Context) ChainID() int64 { return c.chainID }

func (c Context) Height() int64 { return c.height }

func (c Context) GroupNumbers() []uint32 { return c.groupNumbers }

func (c Context) GroupThresholdPercentages() []uint32 { return c.groupThresholdPercentages }

func (c Context) RequestData() []byte { return c.requestData }

func (c Context) Operators() []*avsitypes.Operator { return c.operators }

func (c Context) Logger() log.Logger { return c.logger }

func (c Context) ValidatedResponse() *dvstypes.RequestPostRequestValidatedData {
	return c.validatedResponse
}

// MultiStore returns the MultiStore for this context.
func (c Context) MultiStore() storetypes.MultiStore { return c.ms }

// KVStore returns the KV store for a specific store key.
func (c Context) KVStore(key storetypes.StoreKey) storetypes.KVStore {
	return c.ms.GetKVStore(key)
}

// CacheContext returns a new Context with the multi-store cached and a new
// EventManager. The cached context is written to the context when writeCache
// is called. Note, events are automatically emitted on the parent context's
// EventManager when the caller executes the write.
func (c Context) CacheContext() (cc Context, writeCache func()) {
	cms := c.ms.CacheMultiStore()
	cc = c.WithMultiStore(cms).WithEventManager(NewEventManager())

	writeCache = func() {
		c.EventManager().EmitEvents(cc.EventManager().Events())
		cms.Write()
	}

	return cc, writeCache
}

func (c Context) Value(key any) any {
	if key == ContextKey {
		return c
	}
	return c.baseCtx.Value(key)
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.baseCtx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.baseCtx.Done()
}

func (c Context) Err() error {
	return c.baseCtx.Err()
}

type ContextOption func(Context) Context

// todo delete header
func NewContext(baseCtx context.Context, ms storetypes.MultiStore, logger log.Logger, options ...ContextOption) Context {
	ctx := Context{
		baseCtx:      baseCtx,
		ms:           ms,
		logger:       logger,
		eventManager: NewEventManager(),
	}
	for _, option := range options {
		ctx = option(ctx)
	}
	return ctx
}

func (c Context) WithLogger(logger log.Logger) Context {
	c.logger = logger
	return c
}

func (c Context) WithValue(key, value any) Context {
	c.baseCtx = context.WithValue(c.baseCtx, key, value)
	return c
}

// WithContext returns a Context with an updated context.Context.
func (c Context) WithContext(ctx context.Context) Context {
	c.baseCtx = ctx
	return c
}

// WithEventManager returns a Context with an updated event manager
func (c Context) WithEventManager(em EventManagerI) Context {
	c.eventManager = em
	return c
}

func (c Context) WithChainID(chainID int64) Context {
	c.chainID = chainID
	return c
}

func (c Context) WithHeight(height int64) Context {
	c.height = height
	return c
}

func (c Context) WithGroupNumbers(groupNumbers []uint32) Context {
	c.groupNumbers = groupNumbers
	return c
}

func (c Context) WithRequestData(requestData []byte) Context {
	c.requestData = requestData
	return c
}

func (c Context) WithOperator(operators []*avsitypes.Operator) Context {
	c.operators = operators
	return c
}

func (c Context) WithGroupThresholdPercentages(groupThresholdPercentages []uint32) Context {
	c.groupThresholdPercentages = groupThresholdPercentages
	return c
}

func (c Context) WithValidatedResponse(validatedData *dvstypes.RequestPostRequestValidatedData) Context {
	c.validatedResponse = validatedData
	return c
}

// WithMultiStore returns a Context with an updated MultiStore.
func (c Context) WithMultiStore(ms storetypes.MultiStore) Context {
	c.ms = ms
	return c
}

func UnwrapContext(ctx context.Context) Context {
	if sdkCtx, ok := ctx.(Context); ok {
		return sdkCtx
	}
	return ctx.Value(ContextKey).(Context)
}
