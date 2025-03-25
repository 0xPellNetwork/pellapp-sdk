package baseapp

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	crypto "github.com/0xPellNetwork/pelldvs/proto/pelldvs/crypto"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/jinzhu/copier" //nolint:depguard
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// Supported AVSI Query prefixes and paths
const (
	QueryPathApp    = "app"
	QueryPathCustom = "custom"
	QueryPathP2P    = "p2p"
	QueryPathStore  = "store"

	QueryPathBroadcastTx = "/cosmos.tx.v1beta1.Service/BroadcastTx"
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
	defer func() {
		if r := recover(); r != nil {
			resp = queryResult(errorsmod.Wrapf(sdkerrors.ErrPanic, "%v", r), app.trace)
		}
	}()

	// when a client did not provide a req height, manually inject the latest
	if req.Height == 0 {
		req.Height = app.LastBlockHeight()
	}

	telemetry.IncrCounter(1, "req", "count")
	telemetry.IncrCounter(1, "req", req.Path)
	start := telemetry.Now()
	defer telemetry.MeasureSince(start, req.Path)

	if req.Path == QueryPathBroadcastTx {
		return queryResult(errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "can't route a broadcast tx message"), app.trace), nil
	}

	// handle gRPC routes first rather than calling splitPath because '/' characters
	// are used as part of gRPC paths
	if grpcHandler := app.grpcQueryRouter.Route(req.Path); grpcHandler != nil {
		return app.handleQueryGRPC(grpcHandler, req), nil
	}

	path := SplitAVSIQueryPath(req.Path)
	if len(path) == 0 {
		return queryResult(errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "no req path provided"), app.trace), nil
	}

	switch path[0] {
	case QueryPathStore:
		resp = handleQueryStore(app, path, *req)
	default:
		resp = queryResult(errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "unknown req path"), app.trace)
	}

	return resp, nil
}

// SplitAVSIQueryPath splits a string path using the delimiter '/'.
//
// e.g. "this/is/funny" becomes []string{"this", "is", "funny"}
func SplitAVSIQueryPath(requestPath string) (path []string) {
	path = strings.Split(requestPath, "/")

	// first element is empty string
	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}

	return path
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

func gRPCErrorToSDKError(err error) error {
	status, ok := grpcstatus.FromError(err)
	if !ok {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	switch status.Code() {
	case codes.NotFound:
		return errorsmod.Wrap(sdkerrors.ErrKeyNotFound, err.Error())

	case codes.InvalidArgument:
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())

	case codes.FailedPrecondition:
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())

	case codes.Unauthenticated:
		return errorsmod.Wrap(sdkerrors.ErrUnauthorized, err.Error())

	default:
		return errorsmod.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}
}

func (app *BaseApp) handleQueryGRPC(handler GRPCQueryHandler, req *avsitypes.RequestQuery) *avsitypes.ResponseQuery {
	ctx, err := app.CreateQueryContext(req.Height, req.Prove)
	if err != nil {
		return queryResult(err, app.trace)
	}

	resp, err := handler(ctx, req)
	if err != nil {
		resp = queryResult(gRPCErrorToSDKError(err), app.trace)
		resp.Height = req.Height
		return resp
	}

	return resp
}

// SetGRPCQueryRouter sets the GRPCQueryRouter of the BaseApp.
func (app *BaseApp) SetGRPCQueryRouter(grpcQueryRouter *GRPCQueryRouter) {
	app.grpcQueryRouter = grpcQueryRouter
}

// CreateQueryContext creates a new sdktypes.Context for a query, taking as args
// the block height and whether the query needs a proof or not.
func (app *BaseApp) CreateQueryContext(height int64, prove bool) (sdktypes.Context, error) {
	ctx := sdktypes.NewContext(context.Background(), app.cms, app.logger)
	return ctx, nil
}

// Commit commits the current state of the application.
func (app *BaseApp) Commit() error {
	app.cms.Commit()
	return nil
}

func handleQueryStore(app *BaseApp, path []string, req avsitypes.RequestQuery) *avsitypes.ResponseQuery {
	// "/store" prefix for store queries
	queryable, ok := app.cms.(storetypes.Queryable)
	if !ok {
		return queryResult(errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "multi-store does not support queries"), app.trace)
	}

	req.Path = "/" + strings.Join(path[1:], "/")

	if req.Height <= 1 && req.Prove {
		return queryResult(
			errorsmod.Wrap(
				sdkerrors.ErrInvalidRequest,
				"cannot query with proof when height <= 1; please provide a valid height",
			), app.trace)
	}

	sdkReq := storetypes.RequestQuery(req)
	resp, err := queryable.Query(&sdkReq)
	if err != nil {
		return queryResult(err, app.trace)
	}
	resp.Height = req.Height

	proofs := crypto.ProofOps{Ops: make([]crypto.ProofOp, 0)}
	if resp.ProofOps != nil {
		for i, proof := range proofs.Ops {
			// convert to AVSI proof ops
			proofs.Ops[i] = crypto.ProofOp{
				Type: proof.Type,
				Key:  proof.Key,
				Data: proof.Data,
			}
		}
	} else {
		proofs.Ops = nil
	}

	avsiResp := avsitypes.ResponseQuery{
		Code:      resp.Code,
		Log:       resp.Log,
		Info:      resp.Info,
		Index:     resp.Index,
		Key:       resp.Key,
		Value:     resp.Value,
		ProofOps:  &proofs,
		Height:    resp.Height,
		Codespace: resp.Codespace,
	}

	return &avsiResp
}
