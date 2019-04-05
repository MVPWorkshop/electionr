package types

import (
	"github.com/MVPWorkshop/electionr/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator staking.Validator, found bool)
	// Get the group of the bonded validators
	GetLastValidators(ctx sdk.Context) (validators []staking.Validator)

	GetParams(ctx sdk.Context) staking.Params
	SetParams(ctx sdk.Context, params staking.Params)
	BondDenom(ctx sdk.Context) (res string)

	// Increase non bonded tokens
	InflateSupply(ctx sdk.Context, newTokens sdk.Int)
	GetLastTotalPower(ctx sdk.Context) (power sdk.Int)
}

type BankKeeper interface {
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Tags, sdk.Error)
}
