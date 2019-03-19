package election

import (
	"github.com/MVPWorkshop/legaler-bc/x/election/keeper"
	"github.com/MVPWorkshop/legaler-bc/x/election/tags"
	"github.com/MVPWorkshop/legaler-bc/x/election/types"
)

const (
	RouterKey = types.RouterKey
	StoreKey  = types.StoreKey

	DefaultCodespace = types.DefaultCodespace
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

	RegisterCodec = types.RegisterCodec

	ErrValidatorNotBonded    = types.ErrValidatorNotBonded
	ErrValidatorAlreadyVoted = types.ErrValidatorAlreadyVoted
	ErrGetBlock              = types.ErrGetBlock
	ErrElectionYearFinished  = types.ErrElectionYearFinished
	ErrCycleElectionHasEnded = types.ErrCycleElectionHasEnded

	CycleNum   = tags.CycleNum
	NumVotes   = tags.NumVotes
	IsFinished = tags.IsFinished
)
