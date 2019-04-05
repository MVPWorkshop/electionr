package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	daysInYear = 365
	hoursInDay = 24
)

// Check whether election process has finished
// Returns error in case Tendermint status fetching fails
func IsElectionFinished(ctx sdk.Context) bool {
	// Get first block
	firstBlockContext := ctx.WithBlockHeight(int64(1))
	firstBlock := firstBlockContext.BlockHeader()
	// Get latest block
	latestBlock := ctx.BlockHeader()
	// Check if election year has passed
	timePassed := latestBlock.GetTime().Sub(firstBlock.GetTime())
	if timePassed.Hours()/hoursInDay > daysInYear {
		return true
	}
	return false
}
