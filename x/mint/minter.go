package mint

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Minter represents the minting state
type Minter struct {
	BlockProvisions sdk.Int `json:"block_provisions"` // maximum block provisions
	// TODO: Calculate annual provisions, etc.
}

// Create a new minter object
func NewMinter(blockProvisions sdk.Int) Minter {
	return Minter{
		BlockProvisions: blockProvisions,
	}
}

// Default initial minter object for a new chain
func DefaultInitialMinter() Minter {
	return NewMinter(
		sdk.NewInt(10), // TODO: Allow this or must be in params?
	)
}

func validateMinter(minter Minter) error {
	if minter.BlockProvisions.LT(sdk.ZeroInt()) {
		return fmt.Errorf("mint parameter BlockProvisions should be positive, is %s",
			minter.BlockProvisions.String())
	}
	return nil
}

// Get the provisions for a block based on the annual provisions rate
func (m Minter) BlockProvision(params Params) sdk.Coin {
	return sdk.NewCoin(params.MintDenom, m.BlockProvisions)
}
