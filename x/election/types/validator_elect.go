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
}
