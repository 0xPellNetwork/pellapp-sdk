package pelldvs

import (
	"fmt"

	dvstypes "github.com/pelldvs/pellapp-sdk/pelldvs/types"
	sdktypes "github.com/pelldvs/pellapp-sdk/types"
)

func GetDvsRequestValidatedData(ctx sdktypes.Context) (*dvstypes.RequestPostRequestValidatedData, error) {
	validatedData := ctx.ValidatedResponse()
	if validatedData == nil {
		return nil, fmt.Errorf("not DvsRequestData found")
	}
	return validatedData, nil
}
