package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// expected coin keeper
type DistributionKeeper interface {
	GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins
	GetValidatorOutstandingRewardsCoins(ctx sdk.Context, val sdk.ValAddress) sdk.DecCoins
}

// expected fee collection keeper
type FeeCollectionKeeper interface {
	GetCollectedFees(ctx sdk.Context) sdk.Coins
}

// expected bank keeper
type BankKeeper interface {
	DelegateCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
	UndelegateCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Tags, sdk.Error)
}

// Helpful interfaces for expected election keeper

type ValidatorElect interface {
	GetOperatorAddress() sdk.ValAddress
	HasLeft() bool
	LeaveProtection()
}

type Cycle interface {
	GetPrimaryKey() []byte
	GetValidatorElects() []ValidatorElect
	UpdateValidatorElects(elects []ValidatorElect)
	GetTimeEnded() time.Time
}

// expected election keeper
type ElectionKeeper interface {
	GetAllFinalizedCycles(ctx sdk.Context) (cycles []Cycle)
	SetCycle(ctx sdk.Context, cycle Cycle)
}
