package keeper

import (
	"github.com/MVPWorkshop/legaler-bc/x/election/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Get single validator candidate
func (k Keeper) getValidatorCandidate(ctx sdk.Context, operatorAddr sdk.ValAddress) (valCandidate types.ValidatorCandidate, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(getValidatorCandidateKey(operatorAddr))
	if value == nil {
		return valCandidate, false
	}

	valCandidate = types.MustUnmarshalValidatorCandidate(k.cdc, value)
	return valCandidate, true
}

// Store single validator candidate
func (k Keeper) setValidatorCandidate(ctx sdk.Context, valCandidate types.ValidatorCandidate) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalValidatorCandidate(k.cdc, valCandidate)
	store.Set(getValidatorCandidateKey(valCandidate.OperatorAddr), bz)
}
