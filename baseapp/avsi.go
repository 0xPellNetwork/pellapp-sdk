package baseapp

import (
	"context"

	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"

	dm "github.com/0xPellNetwork/pellapp-sdk/dvs_msg_handler"
	dvstypes "github.com/0xPellNetwork/pellapp-sdk/pelldvs/types"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

func (app *BaseApp) Info(ctx context.Context, info *avsitypes.RequestInfo) (*avsitypes.ResponseInfo, error) {
	return &avsitypes.ResponseInfo{
		Version:         app.version,
		LastBlockHeight: 0,
	}, nil
}

func (app *BaseApp) Query(ctx context.Context, query *avsitypes.RequestQuery) (*avsitypes.ResponseQuery, error) {
	return &avsitypes.ResponseQuery{
		Code: avsitypes.CodeTypeOK,
	}, nil
}

func (app *BaseApp) ProcessDVSRequest(ctx context.Context, req *avsitypes.RequestProcessDVSRequest) (*avsitypes.ResponseProcessDVSRequest, error) {
	sdkCtx := sdktypes.NewContext(ctx)
	sdkCtx = sdkCtx.WithChainID(req.Request.ChainId).
		WithHeight(req.Request.Height).
		WithGroupNumbers(req.Request.GroupNumbers).
		WithRequestData(req.Request.Data).
		WithGroupThresholdPercentages(req.Request.GroupThresholdPercentages).
		WithOperator(req.Operator)

	handlerSrc := dm.GetRequestHandlerSrc()
	res, err := handlerSrc.InvokeRouterRawByData(sdkCtx, req.Request.Data)
	if err != nil {
		app.logger.Error("process request error", "err", err)
		return responseProcessDVSRequestWithEvents(err, sdktypes.MarkEventsToIndex(res.Events, app.indexEvents), app.trace), err
	}

	return &avsitypes.ResponseProcessDVSRequest{
		Log:            res.Log,
		Events:         sdktypes.MarkEventsToIndex(res.Events, app.indexEvents),
		Response:       res.CustomData,
		ResponseDigest: res.CustomDigest,
	}, nil
}

func (app *BaseApp) ProcessDVSResponse(ctx context.Context, req *avsitypes.RequestProcessDVSResponse) (*avsitypes.ResponseProcessDVSResponse, error) {
	sdkCtx := sdktypes.NewContext(ctx)
	sdkCtx = sdkCtx.WithChainID(req.DvsRequest.ChainId).
		WithHeight(req.DvsRequest.Height).
		WithGroupNumbers(req.DvsRequest.GroupNumbers).
		WithRequestData(req.DvsRequest.Data).
		WithGroupThresholdPercentages(req.DvsRequest.GroupThresholdPercentages).
		WithValidatedResponse(dvstypes.NewValidatedResponse(req.DvsResponse))

	handlerSrc := dm.GetResponseHandlerSrc()
	res, err := handlerSrc.InvokeRouterRawByData(sdkCtx, req.DvsRequest.Data)
	if err != nil {
		app.logger.Error("post request error", "err", err)
		return responseProcessDVSResponseWithEvents(err, sdktypes.MarkEventsToIndex(res.Events, app.indexEvents), app.trace), err
	}

	return &avsitypes.ResponseProcessDVSResponse{}, nil
}
