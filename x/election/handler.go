package election

import (
	"crypto/sha256"
	"strconv"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/rpc/core"

	"github.com/MVPWorkshop/legaler-bc/x/election/keeper"
	"github.com/MVPWorkshop/legaler-bc/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	daysInYear = 365
	hoursInDay = 24
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// Note: Message has already passed basic validation
		switch msg := msg.(type) {
		case MsgInsertValidatorElects:
			return handleMsgInsertValidatorElects(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message parse in election module").Result()
		}
	}
}

// These functions assume everything has been authenticated, now we just perform action and save

func handleMsgInsertValidatorElects(ctx sdk.Context, msg MsgInsertValidatorElects, k keeper.Keeper) sdk.Result {
	// TODO: Check contract address
	// Check whether election process is over
	finished, err := isElectionFinished(k)
	if err != nil {
		return err.Result()
	}
	if finished {
		return ErrElectionYearFinished(k.GetCodespace()).Result()
	}

	// Get validator (initiator) from validator set
	initiator, found := k.GetValidator(ctx, msg.InitiatorAddr)
	if !found {
		// Initiator of this transaction is not validator
		return staking.ErrNoValidatorFound(k.GetCodespace()).Result()
	}
	// Initiator shouldn't be jailed
	if initiator.GetJailed() {
		return staking.ErrValidatorJailed(k.GetCodespace()).Result()
	}
	// Initiator must be bonded
	if !initiator.GetStatus().Equal(sdk.Bonded) {
		return ErrValidatorNotBonded(k.GetCodespace()).Result()
	}
	// Check whether there is already finished cycle with this number
	// Note: This means that validators who didn't get to vote for this cycle
	// but it has been elected nonetheless will get the error message, but their
	// vote wouldn't change anything anyway
	if _, found := k.GetFinalizedCycle(ctx, msg.CycleNum); found {
		return ErrCycleFinalized(k.GetCodespace()).Result()
	}

	// Calculate primary key
	primaryKey := calculatePrimaryKey(msg.CycleNum, msg.ElectedValidators)
	// Try to get cycle for this primary key
	cycle, found := k.GetCycle(ctx, primaryKey)
	if !found {
		// If there is no cycle, create one now
		cycle = NewCycle(primaryKey, msg.CycleNum, msg.ElectedValidators, initiator.GetConsPubKey())
	} else {
		// Check if election for this cycle has ended
		if cycle.HasEnded {
			return ErrCycleElectionHasEnded(k.GetCodespace()).Result()
		}
		// Check whether the initiator already voted for this request
		for _, consPubKey := range cycle.ConsPubKeysVoted {
			if consPubKey.Equals(initiator.GetConsPubKey()) {
				return ErrValidatorAlreadyVoted(k.GetCodespace()).Result()
			}
		}

		// Increase number of votes
		cycle.NumVotes = cycle.NumVotes.AddRaw(1)
		// Add initiator's consensus public key to voters
		cycle.ConsPubKeysVoted = append(cycle.ConsPubKeysVoted, initiator.GetConsPubKey())
	}

	// Check whether more than 2/3 of currently active, bonded validators have voted for this cycle
	if hasTwoThirdsMajority(ctx, k.GetLastBondedValidators(ctx), cycle.ConsPubKeysVoted) {
		cycle.HasEnded = true

		// Save latest block time as election time
		latestBlock, e := core.Block(nil)
		if e != nil {
			panic("empty blockchain (no blocks)")
		}
		cycle.TimeEnded = latestBlock.BlockMeta.Header.Time

		// Increment number of max validators
		err = k.IncMaxValidatorsNum(ctx, uint16(len(msg.ElectedValidators)))
		if err != nil {
			return err.Result()
		}
	}

	// Save cycle in state
	k.SetCycle(ctx, cycle)

	tags := sdk.NewTags(
		CycleNum, cycle.Num.String(),
		NumVotes, cycle.NumVotes.String(),
		IsFinished, strconv.FormatBool(cycle.HasEnded),
	)

	return sdk.Result{
		Tags: tags,
	}
}

// Calculates cycle primary key by concatenating cycle number with validator elects
func calculatePrimaryKey(cycleNum sdk.Int, elects []ValidatorElect) Hash {
	// Convert cycle number to byte slice
	data := []byte(cycleNum.String())
	// Append consensus public keys and operator addresses of the elects
	for _, elect := range elects {
		data = append(data, elect.ConsPubKey.Bytes()...)
		data = append(data, elect.OperatorAddr.Bytes()...)
	}
	// Calculate and return SHA256 checksum of the data
	return sha256.Sum256(data)
}

// Check whether election process has finished
// Returns error in case Tendermint status fetching fails
func isElectionFinished(k Keeper) (bool, sdk.Error) {
	// Get first block
	firstBlockNum := int64(1)
	firstBlock, err := core.Block(&firstBlockNum)
	if err != nil {
		return false, ErrGetBlock(k.GetCodespace())
	}
	// Get latest block
	latestBlock, err := core.Block(nil)
	if err != nil {
		return false, ErrGetBlock(k.GetCodespace())
	}
	// Check if election year has passed
	timePassed := latestBlock.BlockMeta.Header.Time.Sub(firstBlock.BlockMeta.Header.Time)
	x := timePassed.Hours()
	if x/hoursInDay > daysInYear {
		return true, nil
	}
	return false, nil
}

// Returns true if more than 2/3 of currently active, bonded validators have voted for this cycle
func hasTwoThirdsMajority(ctx sdk.Context, validators []staking.Validator, consPubKeysVoted []crypto.PubKey) bool {
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
