package types

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValidatorCandidate struct {
	OperatorAddr sdk.ValAddress // Address of the validator's operator
	ConsPubKey   crypto.PubKey  // Consensus public key of the validator
	Standby      bool           // whether the validator is standby or not
}

// Unmarshal validator candidate from a store value
func unmarshalValidatorCandidate(cdc *codec.Codec, value []byte) (valCandidate ValidatorCandidate, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &valCandidate)
	return valCandidate, err
}

// Unmarshal validator candidate from a store value or panic
func MustUnmarshalValidatorCandidate(cdc *codec.Codec, value []byte) ValidatorCandidate {
	valCandidate, err := unmarshalValidatorCandidate(cdc, value)
	if err != nil {
		panic(err)
	}
	return valCandidate
}

func MustMarshalValidatorCandidate(cdc *codec.Codec, valCandidate ValidatorCandidate) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(valCandidate)
}
