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

// Get all cycles that have this cycle number
func (k Keeper) GetCyclesByCycleNum(ctx sdk.Context, cycleNum sdk.Int) (cycles []types.Cycle) {
	store := ctx.KVStore(k.storeKey)

	// Get all cycles
	// TODO: This can be optimized by saving map[cycleNum]cycle instead of just cycle. Then we wouldn't need to iterate through all cycles
	iterator := sdk.KVStorePrefixIterator(store, cycleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		cycle := types.MustUnmarshalCycle(k.cdc, iterator.Value())
		if cycle.Num.Equal(cycleNum) {
			cycles = append(cycles, cycle)
		}
	}

	return cycles
}

// Get finalized cycle
func (k Keeper) GetFinalizedCycle(ctx sdk.Context, cycleNum sdk.Int) (cycle types.Cycle, found bool) {
	store := ctx.KVStore(k.storeKey)

	// Get all cycles
	iterator := sdk.KVStorePrefixIterator(store, cycleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		cycle = types.MustUnmarshalCycle(k.cdc, iterator.Value())
		if cycle.HasEnded {
			return cycle, true
		}
	}

	return cycle, false
}
