package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidInput     = 101
	CodeInvalidValidator = 102
	CodeInvalidTime      = 103
	CodeInvalidCycle     = 104
)

func ErrNilValidatorAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "Initiator operator address is required")
}

func ErrCycleNumberOutOfBounds(codespace sdk.CodespaceType, maxCycles int) sdk.Error {
	return sdk.NewError(
		codespace,
		CodeInvalidInput,
		fmt.Sprintf("cycle number must be a positive integer lower than %d", maxCycles),
	)
}

func ErrValidatorElectsOutOfBounds(codespace sdk.CodespaceType, maxElects int) sdk.Error {
	return sdk.NewError(
		codespace,
		CodeInvalidInput,
		fmt.Sprintf("validator elects shouldn't be empty or larger than %d", maxElects),
	)
}

func ErrNilValidatorElectPubKey(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "elect public key shouldn't be empty")
}

func ErrNilValidatorElectAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "elect operator address shouldn't be empty")
}

func ErrElectPlaceOutOfBounds(codespace sdk.CodespaceType, lastPlace int) sdk.Error {
	return sdk.NewError(
		codespace,
		CodeInvalidInput,
		fmt.Sprintf("cycle number must be a positive integer lower than %d", lastPlace),
	)
}

func ErrValidatorNotBonded(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator should be bonded")
}

func ErrValidatorAlreadyVoted(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "the initiator has already voted for this request")
}

func ErrElectionYearFinished(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTime, "election year has passed")
}

func ErrCycleElectionHasMajority(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTime, "election for this cycle has already gained a majority vote")
}

func ErrCycleNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "cycle of that number does not exist")
}

func ErrCycleNumInvalid(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "cycle number must be a positive integer")
}

func ErrCycleFinalized(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidCycle, "cycle with this number has already been finalized")
}
