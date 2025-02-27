package types

import proto "github.com/cosmos/gogoproto/proto"

type ResultCustomizedIFace interface {
	GetData(proto.Message) ([]byte, error)
	GetDigest(proto.Message) ([]byte, error)
}

type DvsResult struct {
	*Result
	CustomData   []byte
	CustomDigest []byte
}
