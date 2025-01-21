package resulthandler

import (
	sdk "github.com/0xPellNetwork/pellapp-sdk/types"

	"github.com/cosmos/gogoproto/proto"
)

type ResultCustomizedIFace interface {
	GetData(proto.Message) ([]byte, error)
	GetDigest(proto.Message) ([]byte, error)
}

type Result struct {
	*sdk.Result
	CustomData   []byte
	CustomDigest []byte
}
