package election

import (
	"crypto/sha256"

	"github.com/MVPWorkshop/legaler-bc/x/staking"
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

// Returns true if more than 2/3 of currently active, bonded validators have voted for this cycle
func hasTwoThirdsMajority(validators []staking.Validator, consPubKeysVoted []crypto.PubKey) bool {
	activeValidatorsNum := 0
	votersStillActive := 0

	// Iterate through active (bonded) validators from latest block
	for _, validator := range validators {
		// Check that validator that voted for this cycle is still active
		for _, consPubKey := range consPubKeysVoted {
			// Validator should be bonded and not jailed
			if consPubKey.Equals(validator.GetConsPubKey()) && validator.GetStatus().Equal(sdk.Bonded) && !validator.GetJailed() {
				votersStillActive++
				break
			}
		}
		// Increment active validators number
		activeValidatorsNum++
	}

	quorum := activeValidatorsNum*2/3 + 1
	if votersStillActive >= quorum {
		return true
	}
	return false
}
