package keeper

import (
	"github.com/MVPWorkshop/legaler-bc/x/election/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Get single validator elect
func (k Keeper) getValidatorElect(ctx sdk.Context, operatorAddr sdk.ValAddress) (valElect types.ValidatorElect, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(getValidatorElectKey(operatorAddr))
	if value == nil {
		return valElect, false
	}

	valElect = types.MustUnmarshalValidatorElect(k.cdc, value)
	return valElect, true
}

// Store single validator elect
func (k Keeper) setValidatorElect(ctx sdk.Context, valElect types.ValidatorElect) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalValidatorElect(k.cdc, valElect)
	store.Set(getValidatorElectKey(valElect.OperatorAddr), bz)
}
