package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidInput          = 101
	CodeInvalidValidatorElect = 102
)

func ErrNilValidatorAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "validator operator address is nil")
}

func ErrCycleNumberOutOfBounds(codespace sdk.CodespaceType, maxCycles int) sdk.Error {
	return sdk.NewError(
		codespace,
		CodeInvalidValidatorElect,
		fmt.Sprintf("cycle number must be a positive integer lower than %d", maxCycles),
	)
}

func ErrPublicKeysOutOfBounds(codespace sdk.CodespaceType, maxPubKeys int) sdk.Error {
	return sdk.NewError(
		codespace,
		CodeInvalidValidatorElect,
		fmt.Sprintf("validator elect public keys shouldn't be empty or larger than %d", maxPubKeys),
	)
}

func ErrValidatorAddressesOutOfBounds(codespace sdk.CodespaceType, maxValidators int) sdk.Error {
	return sdk.NewError(
		codespace,
		CodeInvalidValidatorElect,
		fmt.Sprintf("validator elect operator addresses shouldn't be empty or larger than %d", maxValidators),
	)
}

func ErrNilValidatorElectorAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "validator elect operator address is nil")
}

func ErrPubKeysValAddressesMissmatch(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(
		codespace,
		CodeInvalidValidatorElect,
		"Number of public keys should be equal to number of validator operator addresses",
	)
}
