package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type DvsResponseHandler interface {
	ResponseHandler(ctx Context, msg sdk.Msg) error
}
