package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	validatorElectsKey = []byte{0x61} // prefix for each key to a validator elect
	cycleKey           = []byte{0x62} // prefix for each key to an election cycle
)

// Get the key for the validator elect with address
func getValidatorElectKey(operatorAddr sdk.ValAddress) []byte {
	return append(validatorElectsKey, operatorAddr.Bytes()...)
}

// Get the key for the cycle election with cycle number
func getCycleKey(cycleNumber sdk.Int) []byte {
	return append(cycleKey, cycleNumber.BigInt().Bytes()...)
}
