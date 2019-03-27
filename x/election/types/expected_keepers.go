package types

import (
	"github.com/MVPWorkshop/legaler-bc/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator types.Validator, found bool)
	// Get the group of the bonded validators
	GetLastValidators(ctx sdk.Context) (validators []types.Validator)

	GetParams(ctx sdk.Context) types.Params
	SetParams(ctx sdk.Context, params types.Params)
}
