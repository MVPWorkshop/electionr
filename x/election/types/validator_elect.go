package types

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxValidatorElectsPerCycle = 7

type ValidatorElect struct {
	OperatorAddr sdk.ValAddress // Address of the validator's operator
	ConsPubKey   crypto.PubKey  // Consensus public key of the validator
	Standby      bool           // whether the validator is standby or not
}

// Unmarshal validator elect from a store value
func unmarshalValidatorElect(cdc *codec.Codec, value []byte) (valElect ValidatorElect, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &valElect)
	return valElect, err
}

// Unmarshal validator elect from a store value or panic
func MustUnmarshalValidatorElect(cdc *codec.Codec, value []byte) ValidatorElect {
	valElect, err := unmarshalValidatorElect(cdc, value)
	if err != nil {
		panic(err)
	}
	return valElect
}

func MustMarshalValidatorElect(cdc *codec.Codec, valElect ValidatorElect) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(valElect)
}
