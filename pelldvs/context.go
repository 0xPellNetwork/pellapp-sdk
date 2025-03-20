package pelldvs

import (
	"fmt"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

func GetDvsRequestValidatedData(ctx sdktypes.Context) (*avsitypes.DVSResponse, error) {
	validatedData := ctx.ValidatedResponse()
	if validatedData == nil {
		return nil, fmt.Errorf("not DvsRequestData found")
	}

	return validatedData, nil
}
