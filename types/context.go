package types

import (
	"context"
	"time"

	"github.com/0xPellNetwork/pelldvs-libs/log"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
)

type ContextKeyType string

const (
	ContextKey ContextKeyType = "pkg_context"
)

type Context struct {
	baseCtx                   context.Context
	eventManager              EventManagerI
	chainID                   int64
	height                    int64
	groupNumbers              []uint32
	groupThresholdPercentages []uint32
	requestData               []byte
	operators                 []*avsitypes.Operator
	validatedResponse         *avsitypes.DVSResponse
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

func (c Context) ValidatedResponse() *avsitypes.DVSResponse {
	return c.validatedResponse
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

// todo delete header
func NewContext(baseCtx context.Context) Context {
	return Context{
		baseCtx:      baseCtx,
		eventManager: NewEventManager(),
	}
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

func (c Context) WithValidatedResponse(validatedData *avsitypes.DVSResponse) Context {
	c.validatedResponse = validatedData
	return c
}

func UnwrapContext(ctx context.Context) Context {
	if sdkCtx, ok := ctx.(Context); ok {
		return sdkCtx
	}
	return ctx.Value(ContextKey).(Context)
}
