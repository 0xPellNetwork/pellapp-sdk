package types

import (
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
)

func NewValidatedResponse(validatedData *avsitypes.DVSResponse) *avsitypes.DVSResponse {
	var nonSignerStakeIndices = make([]*avsitypes.NonSignerStakeIndice, 0)
	for _, nonSignerStakeIndice := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices = append(nonSignerStakeIndices, &avsitypes.NonSignerStakeIndice{
			NonSignerStakeIndice: nonSignerStakeIndice.NonSignerStakeIndice,
		})
	}

	return &avsitypes.DVSResponse{
		Data:                        validatedData.Data,
		Error:                       validatedData.Error,
		Hash:                        validatedData.Hash,
		NonSignersPubkeysG1:         validatedData.NonSignersPubkeysG1,
		GroupApksG1:                 validatedData.GroupApksG1,
		SignersApkG2:                validatedData.SignersApkG2,
		SignersAggSigG1:             validatedData.SignersAggSigG1,
		NonSignerGroupBitmapIndices: validatedData.NonSignerGroupBitmapIndices,
		GroupApkIndices:             validatedData.GroupApkIndices,
		TotalStakeIndices:           validatedData.TotalStakeIndices,
		NonSignerStakeIndices:       nonSignerStakeIndices,
	}
}
