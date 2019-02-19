package mint

import (
	"fmt"

	stakeTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// mint parameters
type Params struct {
	MintDenom string `json:"mint_denom"` // type of coin to mint
}

func NewParams(mintDenom string) Params {
	return Params{
		MintDenom: mintDenom,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom: stakeTypes.DefaultBondDenom,
	}
}

func validateParams(params Params) error {
	if params.MintDenom == "" {
		return fmt.Errorf("mint parameter MintDenom can't be an empty string")
	}
	return nil
}
