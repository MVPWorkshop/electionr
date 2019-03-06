package keeper

import (
	"github.com/MVPWorkshop/legaler-bc/x/election/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const maxCycles = 12

type Keeper struct {
	storeKey sdk.StoreKey

	cycles map[sdk.Int]types.Cycle
	currentCycle sdk.Int

	cdc *codec.Codec // Codec for binary encoding/decoding
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cycles: make(map[sdk.Int]types.Cycle, maxCycles),
		currentCycle: sdk.NewInt(0),
		cdc: cdc,
	}
	return keeper
}

// Increment current cycle number
// Panics if maximum number of cycles has already been reached
func (k *Keeper) incCurrentCycle() {
	if k.currentCycle.GTE(sdk.NewInt(maxCycles)) {
		panic("last cycle already reached")
	}
	k.currentCycle.Add(sdk.NewInt(1))
}
