package election

import (
	"github.com/MVPWorkshop/legaler-bc/x/election/keeper"
	"github.com/MVPWorkshop/legaler-bc/x/election/querier"
	"github.com/MVPWorkshop/legaler-bc/x/election/tags"
	"github.com/MVPWorkshop/legaler-bc/x/election/types"
)

const (
	RouterKey    = types.RouterKey
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute

	DefaultCodespace = types.DefaultCodespace

	MaxCycles = types.MaxCycles
)

type (
	Hash = types.Hash

	Keeper = keeper.Keeper
	Cycle = types.Cycle
	ValidatorElect = types.ValidatorElect

	MsgInsertValidatorElects = types.MsgInsertValidatorElects
)

var (
	NewKeeper                   = keeper.NewKeeper
	NewCycle                    = types.NewCycle
	NewMsgInsertValidatorElects = types.NewMsgInsertValidatorElects
	NewQuerier                  = querier.NewQuerier

	RegisterCodec = types.RegisterCodec

	ErrValidatorNotBonded     = types.ErrValidatorNotBonded
	ErrValidatorAlreadyVoted  = types.ErrValidatorAlreadyVoted
	ErrGetBlock               = types.ErrGetBlock
	ErrElectionYearFinished   = types.ErrElectionYearFinished
	ErrCycleElectionHasEnded  = types.ErrCycleElectionHasEnded
	ErrCycleNumInvalid        = types.ErrCycleNumInvalid
	ErrCycleNumberOutOfBounds = types.ErrCycleNumberOutOfBounds

	CycleNum   = tags.CycleNum
	NumVotes   = tags.NumVotes
	IsFinished = tags.IsFinished
)
