package mint

import sdk "github.com/cosmos/cosmos-sdk/types"

// Expected staking keeper
type StakingKeeper interface {
	// TODO: We'll need validators num here
	InflateSupply(ctx sdk.Context, newTokens sdk.Int)
}

// Expected fee collection keeper interface
type FeeCollectionKeeper interface {
	AddCollectedFees(sdk.Context, sdk.Coins) sdk.Coins
}
