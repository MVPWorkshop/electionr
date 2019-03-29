package types

import (
	"github.com/tendermint/tendermint/crypto"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxCycles = 12

type Hash [32]byte

type Cycle struct {
	PrimaryKey      Hash             `json:"primary_key"`      // Hash(Num + ValidatorElects) - used as a primary key
	Num             sdk.Int          `json:"cycle_num"`        // Cycle number
	ValidatorElects []ValidatorElect `json:"validator_elects"` // Both validator and standby validator elects
	NumVotes        sdk.Int          `json:"num_votes"`        // Number of votes for this cycle election
	// Note: We cannot use []sdk.Validator here instead of []crypto.PubKey
	// because it cannot be marshaled, since it isn't registered with amino
	ConsPubKeysVoted []crypto.PubKey `json:"cons_pub_keys_voted"` // Consensus public key of validators that voted for this cycle
	HasEnded         bool            `json:"has_ended"`           // Whether the cycle has gained majority vote or not
	TimeEnded        time.Time       `json:"time_ended"`          // Block time that represents the moment of gaining majority vote
}

func NewCycle(pk Hash, num sdk.Int, valElects []ValidatorElect, initiatorPubKey crypto.PubKey) Cycle {
	return Cycle{
		PrimaryKey:       pk,
		Num:              num,
		ValidatorElects:  valElects,
		NumVotes:         sdk.OneInt(),
		ConsPubKeysVoted: []crypto.PubKey{initiatorPubKey},
		HasEnded:         false,
	}
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
