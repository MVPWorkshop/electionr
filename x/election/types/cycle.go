package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxCycles = 12

type Cycle struct {
	Num             sdk.Int          // Cycle number
	validatorElects []ValidatorElect // Both validator and standby validator elects
	NumVotes        sdk.Int          // Number of votes for this cycle election
}

// Unmarshal election cycle from a store value
func unmarshalCycle(cdc *codec.Codec, value []byte) (cycle Cycle, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &cycle)
	return cycle, err
}

// Unmarshal election cycle from a store value or panic
func MustUnmarshalCycle(cdc *codec.Codec, value []byte) Cycle {
	cycle, err := unmarshalCycle(cdc, value)
	if err != nil {
		panic(err)
	}
	return cycle
}

func MustMarshalCycle(cdc *codec.Codec, cycle Cycle) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(cycle)
}
