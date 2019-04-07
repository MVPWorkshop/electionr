package keeper

import (
	"bytes"
	"fmt"
	"sort"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/MVPWorkshop/electionr/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Protection period for newly elected validators.
// They are protected since it would be unfair
// if someone could push them out immediately
// after they have passed proof of determination.
const ProtectionPeriod = 30

// Apply and return accumulated updates to the bonded validator set. Also,
// * Updates the active valset as keyed by LastValidatorPowerKey.
// * Updates the total power as keyed by LastTotalPowerKey.
// * Updates validator status' according to updated powers.
// * Updates the fee pool bonded vs not-bonded tokens.
// * Updates relevant indices.
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
//
// CONTRACT: Only validators with non-zero power or zero-power that were bonded
// at the previous block height or were removed from the validator set entirely
// are returned to Tendermint.
func (k Keeper) ApplyAndReturnValidatorSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate) {

	store := ctx.KVStore(k.storeKey)
	maxValidators := k.GetParams(ctx).MaxValidators
	totalPower := sdk.ZeroInt()

	// Retrieve the last validator set.
	// The persistent set is updated later in this function.
	// (see LastValidatorPowerKey).
	last := k.getLastValidatorsByAddr(ctx)

	// Validators that should be protected
	protectedValidators := make([]types.Validator, 0)

	// Check if election process is underway
	//if !election.IsElectionFinished() {
	if !IsElectionFinished(ctx) {

		// Get latest block
		latestBlock := ctx.BlockHeader()

		// Determine all cycles whose validator elects should be protected
		finalizedCycles := k.electionKeeper.GetAllVotedCycles(ctx)
		for _, cycle := range finalizedCycles {

			// Check if protection period for this cycle is still underway
			timePassed := latestBlock.GetTime().Sub(cycle.GetTimeProtectionStarted())
			if timePassed.Hours()/hoursInDay <= ProtectionPeriod {

				// Did this cycle's state change?
				hasUpdated := false
				elects := cycle.GetValidatorElects()
				// Since protection period hasn't passed, add candidates to protected validators
				for i, elect := range elects {

					// Try to get this validator, and add it to protected validators
					// This will show us whether elect has set up his node
					val, found := k.GetValidator(ctx, elect.GetOperatorAddress())
					// Check whether validator elect wants to leave
					if val.GetStatus() == sdk.Unbonding {
						// Update original object, and not just copied value
						elects[i].LeaveProtection()
						hasUpdated = true
					} else if found && !val.GetJailed() { // Validator elect shouldn't be jailed
						protectedValidators = append(protectedValidators, val)
					}
				}

				// Update cycle state if anything changed
				if hasUpdated {
					cycle.UpdateValidatorElects(elects)
					k.electionKeeper.SetCycle(ctx, cycle)
				}
			}
		}

	}

	// We should push out enough validators with lowest power to make room for protected validators
	protectIndex := int(maxValidators) - len(protectedValidators)
	protectCount := 0

	// Iterate over validators, highest power to lowest.
	iterator := sdk.KVStoreReversePrefixIterator(store, ValidatorsByPowerIndexKey)
	defer iterator.Close()
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {

		// fetch the validator
		valAddr := sdk.ValAddress(iterator.Value())
		validator := k.mustGetValidator(ctx, valAddr)

		if validator.Jailed {
			panic("should never retrieve a jailed validator from the power store")
		}

		// if we get to a zero-power validator (which we don't bond),
		// there are no more possible bonded validators
		if validator.PotentialTendermintPower() == 0 {
			break
		}

		// Check if election process is underway
		if !IsElectionFinished(ctx) {
			// Check if this validator is protected even if he has enough power
			temp := make([]types.Validator, 0)
			for _, protected := range protectedValidators {
				// If he is, remove him from protected slice
				if !validator.OperatorAddress.Equals(protected.OperatorAddress) {
					temp = append(temp, protected)
				} else {
					// Decrease spots which need to be freed for protected validators
					protectIndex++
					break
				}
			}
			protectedValidators = temp

			if protectIndex == count {
				// Substitute validator with lowest power with protected validator
				validator = protectedValidators[protectCount]
				valAddr = validator.OperatorAddress
				// Increase counts to protect next validator elect
				protectIndex++
				protectCount++
			}
		}

		// apply the appropriate state change if necessary
		switch validator.Status {
		case sdk.Unbonded:
			validator = k.unbondedToBonded(ctx, validator)
		case sdk.Unbonding:
			validator = k.unbondingToBonded(ctx, validator)
		case sdk.Bonded:
			// no state change
		default:
			panic("unexpected validator status")
		}

		// fetch the old power bytes
		var valAddrBytes [sdk.AddrLen]byte
		copy(valAddrBytes[:], valAddr[:])
		oldPowerBytes, found := last[valAddrBytes]

		// calculate the new power bytes
		newPower := validator.TendermintPower()
		newPowerBytes := k.cdc.MustMarshalBinaryLengthPrefixed(newPower)

		// update the validator set if power has changed
		if !found || !bytes.Equal(oldPowerBytes, newPowerBytes) {
			updates = append(updates, validator.ABCIValidatorUpdate())

			// set validator power on lookup index
			k.SetLastValidatorPower(ctx, valAddr, newPower)
		}

		// validator still in the validator set, so delete from the copy
		delete(last, valAddrBytes)

		// keep count
		count++
		totalPower = totalPower.Add(sdk.NewInt(newPower))
	}

	// sort the no-longer-bonded validators
	noLongerBonded := sortNoLongerBonded(last)

	// iterate through the sorted no-longer-bonded validators
	for _, valAddrBytes := range noLongerBonded {

		// fetch the validator
		validator := k.mustGetValidator(ctx, sdk.ValAddress(valAddrBytes))

		// bonded to unbonding
		k.bondedToUnbonding(ctx, validator)

		// delete from the bonded validator index
		k.DeleteLastValidatorPower(ctx, sdk.ValAddress(valAddrBytes))

		// update the validator set
		updates = append(updates, validator.ABCIValidatorUpdateZero())
	}

	// set total power on lookup index if there are any updates
	if len(updates) > 0 {
		k.SetLastTotalPower(ctx, totalPower)
	}

	return updates
}

// Validator state transitions

func (k Keeper) bondedToUnbonding(ctx sdk.Context, validator types.Validator) types.Validator {
	if validator.Status != sdk.Bonded {
		panic(fmt.Sprintf("bad state transition bondedToUnbonding, validator: %v\n", validator))
	}
	return k.beginUnbondingValidator(ctx, validator)
}

func (k Keeper) unbondingToBonded(ctx sdk.Context, validator types.Validator) types.Validator {
	if validator.Status != sdk.Unbonding {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}
	return k.bondValidator(ctx, validator)
}

func (k Keeper) unbondedToBonded(ctx sdk.Context, validator types.Validator) types.Validator {
	if validator.Status != sdk.Unbonded {
		panic(fmt.Sprintf("bad state transition unbondedToBonded, validator: %v\n", validator))
	}
	return k.bondValidator(ctx, validator)
}

// switches a validator from unbonding state to unbonded state
func (k Keeper) unbondingToUnbonded(ctx sdk.Context, validator types.Validator) types.Validator {
	if validator.Status != sdk.Unbonding {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, validator: %v\n", validator))
	}
	return k.completeUnbondingValidator(ctx, validator)
}

// send a validator to jail
func (k Keeper) jailValidator(ctx sdk.Context, validator types.Validator) {
	if validator.Jailed {
		panic(fmt.Sprintf("cannot jail already jailed validator, validator: %v\n", validator))
	}

	validator.Jailed = true
	k.SetValidator(ctx, validator)
	k.DeleteValidatorByPowerIndex(ctx, validator)
}

// remove a validator from jail
func (k Keeper) unjailValidator(ctx sdk.Context, validator types.Validator) {
	if !validator.Jailed {
		panic(fmt.Sprintf("cannot unjail already unjailed validator, validator: %v\n", validator))
	}

	validator.Jailed = false
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)
}

// perform all the store operations for when a validator status becomes bonded
func (k Keeper) bondValidator(ctx sdk.Context, validator types.Validator) types.Validator {

	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	// set the status
	pool := k.GetPool(ctx)
	validator, pool = validator.UpdateStatus(pool, sdk.Bonded)
	k.SetPool(ctx, pool)

	// save the now bonded validator record to the two referenced stores
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)

	// delete from queue if present
	k.DeleteValidatorQueue(ctx, validator)

	// trigger hook
	k.AfterValidatorBonded(ctx, validator.ConsAddress(), validator.OperatorAddress)

	return validator
}

// perform all the store operations for when a validator begins unbonding
func (k Keeper) beginUnbondingValidator(ctx sdk.Context, validator types.Validator) types.Validator {

	params := k.GetParams(ctx)

	// delete the validator by power index, as the key will change
	k.DeleteValidatorByPowerIndex(ctx, validator)

	// sanity check
	if validator.Status != sdk.Bonded {
		panic(fmt.Sprintf("should not already be unbonded or unbonding, validator: %v\n", validator))
	}

	// set the status
	pool := k.GetPool(ctx)
	validator, pool = validator.UpdateStatus(pool, sdk.Unbonding)
	k.SetPool(ctx, pool)

	// set the unbonding completion time and completion height appropriately
	validator.UnbondingCompletionTime = ctx.BlockHeader().Time.Add(params.UnbondingTime)
	validator.UnbondingHeight = ctx.BlockHeader().Height

	// save the now unbonded validator record and power index
	k.SetValidator(ctx, validator)
	k.SetValidatorByPowerIndex(ctx, validator)

	// Adds to unbonding validator queue
	k.InsertValidatorQueue(ctx, validator)

	// trigger hook
	k.AfterValidatorBeginUnbonding(ctx, validator.ConsAddress(), validator.OperatorAddress)

	return validator
}

// perform all the store operations for when a validator status becomes unbonded
func (k Keeper) completeUnbondingValidator(ctx sdk.Context, validator types.Validator) types.Validator {
	pool := k.GetPool(ctx)
	validator, pool = validator.UpdateStatus(pool, sdk.Unbonded)
	k.SetPool(ctx, pool)
	k.SetValidator(ctx, validator)
	return validator
}

// map of operator addresses to serialized power
type validatorsByAddr map[[sdk.AddrLen]byte][]byte

// get the last validator set
func (k Keeper) getLastValidatorsByAddr(ctx sdk.Context) validatorsByAddr {
	last := make(validatorsByAddr)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, LastValidatorPowerKey)
	defer iterator.Close()
	// iterate over the last validator set index
	for ; iterator.Valid(); iterator.Next() {
		var valAddr [sdk.AddrLen]byte
		// extract the validator address from the key (prefix is 1-byte)
		copy(valAddr[:], iterator.Key()[1:])
		// power bytes is just the value
		powerBytes := iterator.Value()
		last[valAddr] = make([]byte, len(powerBytes))
		copy(last[valAddr][:], powerBytes[:])
	}
	return last
}

// given a map of remaining validators to previous bonded power
// returns the list of validators to be unbonded, sorted by operator address
func sortNoLongerBonded(last validatorsByAddr) [][]byte {
	// sort the map keys for determinism
	noLongerBonded := make([][]byte, len(last))
	index := 0
	for valAddrBytes := range last {
		valAddr := make([]byte, sdk.AddrLen)
		copy(valAddr[:], valAddrBytes[:])
		noLongerBonded[index] = valAddr
		index++
	}
	// sorted by address - order doesn't matter
	sort.SliceStable(noLongerBonded, func(i, j int) bool {
		// -1 means strictly less than
		return bytes.Compare(noLongerBonded[i], noLongerBonded[j]) == -1
	})
	return noLongerBonded
}
