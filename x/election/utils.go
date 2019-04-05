package election

import (
	"crypto/sha256"

	"github.com/MVPWorkshop/electionr/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// Calculates cycle primary key by concatenating cycle number with validator elects
func calculatePrimaryKey(cycleNum sdk.Int, elects []ValidatorElect) []byte {
	// Convert cycle number to byte slice
	data := []byte(cycleNum.String())
	// Append consensus public keys and operator addresses of the elects
	for _, elect := range elects {
		data = append(data, elect.ConsPubKey.Bytes()...)
		data = append(data, elect.OperatorAddr.Bytes()...)
	}
	// Calculate and return SHA256 checksum of the data
	pk := sha256.Sum256(data)
	return pk[:]
}

// Returns true if majority of currently active, bonded validators have voted for this cycle (their power is checked)
func hasTwoThirdsMajority(validators []staking.Validator, consPubKeysVoted []crypto.PubKey, totalPower int64) bool {
	var votersPower int64

	// Iterate through active (bonded) validators from latest block
	for _, validator := range validators {
		// Check that validator that voted for this cycle is still active
		for _, consPubKey := range consPubKeysVoted {
			// Validator should be bonded and not jailed
			if consPubKey.Equals(validator.GetConsPubKey()) && validator.GetStatus().Equal(sdk.Bonded) && !validator.GetJailed() {
				votersPower += validator.GetTendermintPower()
			}
		}
	}

	quorum := totalPower*2/3 + 1
	if votersPower >= quorum {
		return true
	}
	return false
}
