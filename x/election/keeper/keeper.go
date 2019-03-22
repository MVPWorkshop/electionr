package keeper

import (
	"github.com/MVPWorkshop/legaler-bc/x/election/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey sdk.StoreKey

	validatorSet sdk.ValidatorSet
	cycles       map[sdk.Int]types.Cycle

	cdc       *codec.Codec // Codec for binary encoding/decoding
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, vs sdk.ValidatorSet, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:     key,
		validatorSet: vs,
		cycles:       make(map[sdk.Int]types.Cycle, types.MaxCycles),
		cdc:          cdc,
		codespace:    codespace,
	}
	return keeper
}

// Get validator set
func (k Keeper) GetValidatorSet() sdk.ValidatorSet {
	return k.validatorSet
}

// Get codespace
func (k Keeper) GetCodespace() sdk.CodespaceType {
	return k.codespace
}
