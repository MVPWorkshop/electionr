package keeper

import (
	"github.com/MVPWorkshop/legaler-bc/x/election/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Get single election cycle
func (k Keeper) getCycle(ctx sdk.Context, cycleNum sdk.Int) (cycle types.Cycle, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(getCycleKey(cycleNum))
	if value == nil {
		return cycle, false
	}

	cycle = types.MustUnmarshalCycle(k.cdc, value)
	return cycle, true
}

// Store single election cycle
func (k Keeper) setCycle(ctx sdk.Context, cycle types.Cycle) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalCycle(k.cdc, cycle)
	store.Set(getCycleKey(cycle.Num), bz)
}
