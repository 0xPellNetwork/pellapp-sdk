package pelldvs

import (
	"fmt"

	dvstypes "github.com/0xPellNetwork/pellapp-sdk/pelldvs/types"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

func GetDvsRequestValidatedData(ctx sdktypes.Context) (*dvstypes.RequestPostRequestValidatedData, error) {
	validatedData := ctx.ValidatedResponse()
	if validatedData == nil {
		return nil, fmt.Errorf("not DvsRequestData found")
	}
	return validatedData, nil
}
