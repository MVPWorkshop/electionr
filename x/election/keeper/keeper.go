package keeper

import (
	"github.com/MVPWorkshop/electionr/x/election/types"
	"github.com/MVPWorkshop/electionr/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey sdk.StoreKey

	stakingKeeper types.StakingKeeper
	bankKeeper    types.BankKeeper
	cycles        map[sdk.Int]types.Cycle

	cdc       *codec.Codec // Codec for binary encoding/decoding
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk types.StakingKeeper, bk types.BankKeeper, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:      key,
		stakingKeeper: sk,
		bankKeeper:    bk,
		cycles:        make(map[sdk.Int]types.Cycle, types.MaxCycles),
		cdc:           cdc,
		codespace:     codespace,
	}
	return keeper
}

func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (staking.Validator, bool) {
	return k.stakingKeeper.GetValidator(ctx, addr)
}

// Get the group of the bonded validators
func (k Keeper) GetLastBondedValidators(ctx sdk.Context) []staking.Validator {
	return k.stakingKeeper.GetLastValidators(ctx)
}

// Get codespace
func (k Keeper) GetCodespace() sdk.CodespaceType {
	return k.codespace
}

func (k Keeper) GetTotalPower(ctx sdk.Context) sdk.Int {
	return k.stakingKeeper.GetLastTotalPower(ctx)
}

// Increments maximum number of validators by newValidatorNum
func (k Keeper) IncMaxValidatorsNum(ctx sdk.Context, newValidatorNum uint16) sdk.Error {
	// Check if new validators number is lower than max validators per cycle
	if newValidatorNum > types.MaxValidatorElectsPerCycle {
		return types.ErrValidatorElectsOutOfBounds(k.codespace, types.MaxValidatorElectsPerCycle)
	}
	params := k.stakingKeeper.GetParams(ctx)
	params.MaxValidators += newValidatorNum
	k.stakingKeeper.SetParams(ctx, params)

	return nil
}
