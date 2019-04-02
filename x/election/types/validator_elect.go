package types

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxValidatorElectsPerCycle = 7

type ValidatorElect struct {
	OperatorAddr sdk.ValAddress `json:"operator_addr"` // Address of the validator's operator
	ConsPubKey   crypto.PubKey  `json:"cons_pub_key"`  // Consensus public key of the validator
	Place        sdk.Int        `json:"place"`         // Place that he achieved in the PoD "race"
	Left         bool           `json:"left"`          // If elect decided to unbond within protection period
}

func NewValidatorElect(operAddr sdk.ValAddress, consPubKey crypto.PubKey, place sdk.Int) ValidatorElect {
	return ValidatorElect{
		OperatorAddr: operAddr,
		ConsPubKey:   consPubKey,
		Place:        place,
		Left:         false,
	}
}

func (v ValidatorElect) GetOperatorAddress() sdk.ValAddress {
	return v.OperatorAddr
}

// Check if validator elect has left protection intentionally
func (v ValidatorElect) HasLeft() bool {
	return v.Left
}

// If validator elect decides to leave within protection period
func (v *ValidatorElect) LeaveProtection() {
	v.Left = true
}

// Since public key cannot be put in JSON we need to convert it to/from string
type ValidatorElectJSON struct {
	OperatorAddr sdk.ValAddress `json:"operator_addr"` // Address of the validator's operator
	ConsPubKey   string         `json:"cons_pub_key"`  // Consensus public key of the validator
	Place        sdk.Int        `json:"place"`         // Place that he achieved in the PoD "race"
}

func NewValidatorElectJSON(operAddr sdk.ValAddress, consPubKey string, place sdk.Int) ValidatorElectJSON {
	return ValidatorElectJSON{
		OperatorAddr: operAddr,
		ConsPubKey:   consPubKey,
		Place:        place,
	}
}
