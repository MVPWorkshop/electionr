package types

import (
	"time"

	"github.com/tendermint/tendermint/crypto"

	"github.com/MVPWorkshop/electionr/x/staking"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const MaxCycles = 12

type Cycle struct {
	PrimaryKey      []byte           `json:"primary_key"`      // Hash(Num + ValidatorElects) - used as a primary key
	Num             sdk.Int          `json:"cycle_num"`        // Cycle number
	ValidatorElects []ValidatorElect `json:"validator_elects"` // Elects who have passed Proof of Determination
	NumVotes        sdk.Int          `json:"num_votes"`        // Number of votes for this cycle election
	// Note: We cannot use []sdk.Validator here instead of []crypto.PubKey
	// because it cannot be marshaled, since it isn't registered with amino
	ConsPubKeysVoted []crypto.PubKey `json:"cons_pub_keys_voted"` // Consensus public key of validators that voted for this cycle
	HasEnded         bool            `json:"has_ended"`           // Whether the cycle has gained majority vote or not
	TimeEnded        time.Time       `json:"time_ended"`          // Block time that represents the moment of gaining majority vote
}

func NewCycle(pk []byte, num sdk.Int, valElects []ValidatorElect, initiatorPubKey crypto.PubKey) Cycle {
	return Cycle{
		PrimaryKey:       pk,
		Num:              num,
		ValidatorElects:  valElects,
		NumVotes:         sdk.OneInt(),
		ConsPubKeysVoted: []crypto.PubKey{initiatorPubKey},
		HasEnded:         false,
	}
}

func (c Cycle) GetPrimaryKey() []byte {
	return c.PrimaryKey
}

func (c Cycle) GetValidatorElects() []staking.ValidatorElect {
	elects := make([]staking.ValidatorElect, 0)
	for _, elect := range c.ValidatorElects {
		elects = append(elects, &elect)
	}
	return elects
}

func (c *Cycle) UpdateValidatorElects(elects []staking.ValidatorElect) {
	c.ValidatorElects = make([]ValidatorElect, len(elects))
	for _, elect := range elects {
		c.ValidatorElects = append(c.ValidatorElects, *elect.(*ValidatorElect))
	}
}

func (c Cycle) GetTimeEnded() time.Time {
	return c.TimeEnded
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

func MustMarshalCycle(cdc *codec.Codec, cycle staking.Cycle) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(cycle)
}

// For cycle display
type CycleJSON struct {
	PrimaryKey      []byte                      `json:"primary_key"`      // Hash(Num + ValidatorElects) - used as a primary key
	Num             sdk.Int                     `json:"cycle_num"`        // Cycle number
	ValidatorElects []ValidatorElectDisplayJSON `json:"validator_elects"` // Elects who have passed Proof of Determination
	NumVotes        sdk.Int                     `json:"num_votes"`        // Number of votes for this cycle election
	// Note: We cannot use []sdk.Validator here instead of []crypto.PubKey
	// because it cannot be marshaled, since it isn't registered with amino
	ConsPubKeysVoted []string  `json:"cons_pub_keys_voted"` // Consensus public keys (bech32) of validators that voted for this cycle
	HasEnded         bool      `json:"has_ended"`           // Whether the cycle has gained majority vote or not
	TimeEnded        time.Time `json:"time_ended"`          // Block time that represents the moment of gaining majority vote
}

func NewCycleJSON(pk []byte, num sdk.Int, valElects []ValidatorElectDisplayJSON,
	consPubKeysVoted []string, hasEnded bool, timeEnded time.Time) CycleJSON {
	return CycleJSON{
		PrimaryKey:       pk,
		Num:              num,
		ValidatorElects:  valElects,
		NumVotes:         sdk.OneInt(),
		ConsPubKeysVoted: consPubKeysVoted,
		HasEnded:         hasEnded,
		TimeEnded:        timeEnded,
	}
}
