package election

import (
	"github.com/MVPWorkshop/electionr/x/election/keeper"
	"github.com/MVPWorkshop/electionr/x/election/querier"
	"github.com/MVPWorkshop/electionr/x/election/tags"
	"github.com/MVPWorkshop/electionr/x/election/types"
)

const (
	ModuleName   = types.ModuleName
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute

	DefaultCodespace = types.DefaultCodespace

	MaxCycles                  = types.MaxCycles
	MaxValidatorElectsPerCycle = types.MaxValidatorElectsPerCycle
)

type (
	Keeper = keeper.Keeper
	Cycle = types.Cycle
	ValidatorElect = types.ValidatorElect
	ValidatorElectJSON = types.ValidatorElectJSON

	MsgInsertValidatorElects = types.MsgInsertValidatorElects
)

var (
	NewKeeper                   = keeper.NewKeeper
	NewCycle                    = types.NewCycle
	NewMsgInsertValidatorElects = types.NewMsgInsertValidatorElects
	NewQuerier                  = querier.NewQuerier
	NewValidatorElect           = types.NewValidatorElect

	RegisterCodec = types.RegisterCodec

	ErrValidatorNotBonded       = types.ErrValidatorNotBonded
	ErrValidatorAlreadyVoted    = types.ErrValidatorAlreadyVoted
	ErrElectionYearFinished     = types.ErrElectionYearFinished
	ErrCycleElectionHasMajority = types.ErrCycleElectionHasMajority
	ErrCycleNumInvalid          = types.ErrCycleNumInvalid
	ErrCycleNumberOutOfBounds   = types.ErrCycleNumberOutOfBounds
	ErrCycleMajorityVote        = types.ErrCycleMajorityVote

	CycleNum        = tags.CycleNum
	NumVotes        = tags.NumVotes
	HasMajorityVote = tags.HasMajorityVote
)
