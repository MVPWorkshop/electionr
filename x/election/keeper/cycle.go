package keeper

import (
	"github.com/MVPWorkshop/legaler-bc/x/election/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Get single election cycle
func (k Keeper) GetCycle(ctx sdk.Context, primaryKey types.Hash) (cycle types.Cycle, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(getCycleKey(primaryKey))
	if value == nil {
		return cycle, false
	}

	cycle = types.MustUnmarshalCycle(k.cdc, value)
	return cycle, true
}

// Store single election cycle
func (k Keeper) SetCycle(ctx sdk.Context, cycle types.Cycle) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalCycle(k.cdc, cycle)
	store.Set(getCycleKey(cycle.PrimaryKey), bz)
}
