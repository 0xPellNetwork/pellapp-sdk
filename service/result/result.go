package result

import (
	"github.com/cosmos/gogoproto/proto"

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// Result extends the sdk.Result structure by adding custom data and digest fields.
// It is used to carry additional information when processing DVS results.
type Result struct {
	*sdktypes.Result        // Embedded original SDK result
	CustomData       []byte // Custom data for storing application-specific result information
	CustomDigest     []byte // Custom digest for verification or indexing purposes
}

// ResultMsgExtractor defines an interface for handling custom result data.
// Types implementing this interface can generate custom data and digests from a proto.Message.
type ResultMsgExtractor interface {
	// GetData extracts custom data from the given proto.Message
	GetData(proto.Message) ([]byte, error)

	// GetDigest generates a custom digest from the given proto.Message
	GetDigest(proto.Message) ([]byte, error)
}
