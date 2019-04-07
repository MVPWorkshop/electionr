package election

import (
	"fmt"
	"strconv"

	"github.com/MVPWorkshop/electionr/x/election/keeper"
	"github.com/MVPWorkshop/electionr/x/staking"
	sk "github.com/MVPWorkshop/electionr/x/staking/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	finished := sk.IsElectionFinished(ctx)
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
		if cycle.HasMajorityVote {
			return ErrCycleElectionHasMajority(k.GetCodespace()).Result()
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
	if hasTwoThirdsMajority(k.GetLastBondedValidators(ctx), cycle.ConsPubKeysVoted, k.GetTotalPower(ctx).Int64()) {
		cycle.HasMajorityVote = true

		// Save latest block time as election time
		latestBlock := ctx.BlockHeader()
		cycle.TimeProtectionStarted = latestBlock.GetTime()

		// Increment number of max validators
		err := k.IncMaxValidatorsNum(ctx, uint16(len(msg.ElectedValidators)))
		if err != nil {
			return err.Result()
		}

		// Add enough tokens to each elect (from this cycle) for him to obtain enough starting power
		err = k.AddInitialCoinsToElects(ctx, cycle.ValidatorElects)

		// Log that this cycle has gained a majority vote
		logger := ctx.Logger().With("module", "x/"+ModuleName)
		logger.Info(fmt.Sprintf("Cycle number %s has gained a majority vote.", cycle.Num.String()))
		logger.Info(fmt.Sprintf(
			"Validator elects from this cycle can now become validators, and will be protected for the next %d days.",
			staking.ProtectionPeriod,
		))
	}

	// Save cycle in state
	k.SetCycle(ctx, &cycle)

	tags := sdk.NewTags(
		CycleNum, cycle.Num.String(),
		NumVotes, cycle.NumVotes.String(),
		HasMajorityVote, strconv.FormatBool(cycle.HasMajorityVote),
	)

	return sdk.Result{
		Tags: tags,
	}
}
