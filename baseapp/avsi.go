package baseapp

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	"github.com/jinzhu/copier" //nolint:depguard

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// Info returns information about the application.
// It implements the AVSI interface by providing version and last block height.
func (app *BaseApp) Info(ctx context.Context, info *avsitypes.RequestInfo) (*avsitypes.ResponseInfo, error) {
	return &avsitypes.ResponseInfo{
		Version: app.version,
	}, nil
}

// Query handles queries to the application.
// It implements the AVSI interface by returning a successful response with OK code.
func (app *BaseApp) Query(ctx context.Context, req *avsitypes.RequestQuery) (resp *avsitypes.ResponseQuery, err error) {
	return &avsitypes.ResponseQuery{
		Code: avsitypes.CodeTypeOK,
	}, nil
}

// ProcessDVSRequest processes a DVS (Distributed Validation System) request.
// It creates an SDK context with request data and invokes the appropriate request handler.
// Returns the handler's response or an error response if processing fails.
func (app *BaseApp) ProcessDVSRequest(ctx context.Context, req *avsitypes.RequestProcessDVSRequest) (*avsitypes.ResponseProcessDVSRequest, error) {
	sdkCtx := sdktypes.NewContext(ctx, app.cms, app.logger)
	sdkCtx = sdkCtx.WithChainID(req.Request.ChainId).
		WithHeight(req.Request.Height).
		WithGroupNumbers(req.Request.GroupNumbers).
		WithRequestData(req.Request.Data).
		WithGroupThresholdPercentages(req.Request.GroupThresholdPercentages).
		WithOperator(req.Operator)

	resp := &avsitypes.ResponseProcessDVSRequest{}
	res, err := app.msgRouter.InvokeByMsgData(sdkCtx, req.Request.Data)
	if err != nil {
		app.logger.Error("process request error", "err", err)

		_ = copier.Copy(resp, sdktypes.WarpAvsiBaseError(err, res, app.trace))
		return resp, err
	}

	return &avsitypes.ResponseProcessDVSRequest{
		Log:            res.Log,
		Events:         sdktypes.MarkEventsToIndex(res.Events, app.indexEvents),
		Response:       res.CustomData,
		ResponseDigest: res.CustomDigest,
	}, nil
}

// ProcessDVSResponse processes a DVS response after validators have processed a request.
// It creates an SDK context with the original request data and validated response,
// then invokes the appropriate response handler.
func (app *BaseApp) ProcessDVSResponse(ctx context.Context, req *avsitypes.RequestProcessDVSResponse) (*avsitypes.ResponseProcessDVSResponse, error) {
	sdkCtx := sdktypes.NewContext(ctx, app.cms, app.logger)
	sdkCtx = sdkCtx.WithChainID(req.DvsRequest.ChainId).
		WithHeight(req.DvsRequest.Height).
		WithGroupNumbers(req.DvsRequest.GroupNumbers).
		WithRequestData(req.DvsRequest.Data).
		WithGroupThresholdPercentages(req.DvsRequest.GroupThresholdPercentages).
		WithValidatedResponse(req.DvsResponse)

	resp := &avsitypes.ResponseProcessDVSResponse{}
	res, err := app.msgRouter.InvokeByMsgData(sdkCtx, req.DvsRequest.Data)
	if err != nil {
		app.logger.Error("post request error", "err", err)

		_ = copier.Copy(resp, sdktypes.WarpAvsiBaseError(err, res, app.trace))
		return resp, err
	}

	return &avsitypes.ResponseProcessDVSResponse{
		Data:   res.CustomData,
		Log:    res.Log,
		Events: sdktypes.MarkEventsToIndex(res.Events, app.indexEvents),
	}, nil
}

// createQueryContext creates a new sdk.Context for a query, taking as args
// the block height and whether the query needs a proof or not.
func (app *BaseApp) CreateQueryContext() (sdktypes.Context, error) {
	// use custom query multi-store if provided
	qms := app.qms
	if qms == nil {
		qms = app.cms.(storetypes.MultiStore)
	}

	cacheMS := qms.CacheMultiStore()
	// branch the commit multi-store for safety
	ctx := sdktypes.NewContext(context.Background(), cacheMS, app.logger)
	return ctx, nil
}
