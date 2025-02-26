package result

import (
	"github.com/cosmos/gogoproto/proto"

	sdk "github.com/0xPellNetwork/pellapp-sdk/types"
)

// Result extends the sdk.Result structure by adding custom data and digest fields
// Used to carry additional information when processing dvs results
type Result struct {
	*sdk.Result         // Embedded original SDK result
	CustomData   []byte // Custom data, can be used to store application-specific result data
	CustomDigest []byte // Custom digest, can be used for verification or indexing purposes
}

// CustomResultHandler defines an interface for handling custom result data
// Types implementing this interface can generate custom data and digests from a proto.Message
type CustomResultHandler interface {
	// GetData gets custom data from the given proto.Message
	GetData(proto.Message) ([]byte, error)

	// GetDigest gets a custom digest from the given proto.Message
	GetDigest(proto.Message) ([]byte, error)
}
