package mint

import sdk "github.com/cosmos/cosmos-sdk/types"

// Inflate every block with fixed amount
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// Fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// Mint coins, add to collected fees, update supply
	mintedCoin := minter.BlockProvision(params)
	k.fck.AddCollectedFees(ctx, sdk.Coins{mintedCoin})
	k.sk.InflateSupply(ctx, mintedCoin.Amount)
}