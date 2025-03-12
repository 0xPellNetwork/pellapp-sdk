package types

import (
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
)

func NewValidatedResponse(validatedData *avsitypes.DVSResponse) *RequestPostRequestValidatedData {
	var nonSignerStakeIndices []*NonSignerStakeIndice
	for _, nonSignerStakeIndice := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices = append(nonSignerStakeIndices, &NonSignerStakeIndice{
			NonSignerStakeIndice: nonSignerStakeIndice.NonSignerStakeIndice,
		})
	}

	return &RequestPostRequestValidatedData{
		Data:                         validatedData.Data,
		Error:                        validatedData.Error,
		Hash:                         validatedData.Hash,
		NonSignersPubkeysG1:          validatedData.NonSignersPubkeysG1,
		QuorumApksG1:                 validatedData.GroupApksG1,
		SignersApkG2:                 validatedData.SignersApkG2,
		SignersAggSigG1:              validatedData.SignersAggSigG1,
		NonSignerQuorumBitmapIndices: validatedData.NonSignerGroupBitmapIndices,
		QuorumApkIndices:             validatedData.GroupApkIndices,
		TotalStakeIndices:            validatedData.TotalStakeIndices,
		NonSignerStakeIndices:        nonSignerStakeIndices,
	}
}
