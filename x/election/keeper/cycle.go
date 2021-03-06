package keeper

import (
	"github.com/MVPWorkshop/electionr/x/election/types"
	"github.com/MVPWorkshop/electionr/x/staking"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Starting power for newly elected validators
const electStartingPower = 1

// Get single election cycle
func (k Keeper) GetCycle(ctx sdk.Context, primaryKey []byte) (cycle types.Cycle, found bool) {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(getCycleKey(primaryKey))
	if value == nil {
		return cycle, false
	}

	cycle = types.MustUnmarshalCycle(k.cdc, value)
	return cycle, true
}

// Store single election cycle
func (k Keeper) SetCycle(ctx sdk.Context, cycle staking.Cycle) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalCycle(k.cdc, cycle)
	store.Set(getCycleKey(cycle.GetPrimaryKey()), bz)
}

// Get all cycles that have this cycle number
func (k Keeper) GetCyclesByCycleNum(ctx sdk.Context, cycleNum sdk.Int) (cycles []types.Cycle) {
	store := ctx.KVStore(k.storeKey)

	// Get all cycles
	// TODO: This can be optimized by saving map[cycleNum]cycle instead of just cycle
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

// Get cycle that has gained majority vote
func (k Keeper) GetVotedCycle(ctx sdk.Context, cycleNum sdk.Int) (cycle types.Cycle, found bool) {
	store := ctx.KVStore(k.storeKey)

	// Get all cycles
	iterator := sdk.KVStorePrefixIterator(store, cycleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		cycle = types.MustUnmarshalCycle(k.cdc, iterator.Value())
		if cycle.HasMajorityVote {
			return cycle, true
		}
	}

	return cycle, false
}

// Get cycles that have gained majority vote
func (k Keeper) GetAllVotedCycles(ctx sdk.Context) (cycles []staking.Cycle) {
	store := ctx.KVStore(k.storeKey)

	// Get all cycles
	iterator := sdk.KVStorePrefixIterator(store, cycleKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		cycle := types.MustUnmarshalCycle(k.cdc, iterator.Value())
		if cycle.HasMajorityVote {
			cycles = append(cycles, &cycle)
		}
	}

	return cycles
}

// Adds coins to validator elects (and inflate supply) by saving them in state
func (k Keeper) AddInitialCoinsToElects(ctx sdk.Context, elects []types.ValidatorElect) {
	// Coins to add in order to gain initial power
	amount := sdk.TokensFromTendermintPower(int64(electStartingPower))
	coins := sdk.Coins{
		sdk.NewCoin(
			k.stakingKeeper.BondDenom(ctx),
			amount,
		),
	}
	totalAmount := sdk.ZeroInt()
	for _, elect := range elects {
		_, _, err := k.bankKeeper.AddCoins(ctx, sdk.AccAddress(elect.OperatorAddr), coins)
		if err != nil {
			// Shouldn't happen
			panic(err.Error())
		}
		totalAmount = totalAmount.Add(amount)
	}

	// Inflate coin supply
	k.stakingKeeper.InflateSupply(ctx, totalAmount)
}
