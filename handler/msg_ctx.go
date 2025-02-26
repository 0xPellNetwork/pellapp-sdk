package dvsservermanager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ctxDvsRequestKey = "CTX_DVS_POST_RESPONSE"
)

func CtxWithDvsPostResponseData(ctx sdk.Context, postProcessResponseData []byte) sdk.Context {
	return ctx.WithValue(ctxDvsRequestKey, postProcessResponseData)
}

func CtxGetDvsPostResponseData(ctx sdk.Context) ([]byte, bool) {
	value := ctx.Value(ctxDvsRequestKey)
	val, ok := value.([]byte)
	return val, ok
}
