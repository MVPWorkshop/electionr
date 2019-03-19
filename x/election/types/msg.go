package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Ensure Msg interface compliance at compile time
// (It implements Msg interface)
var _ sdk.Msg = &MsgInsertValidatorElects{}

//______________________________________________________________________

// MsgInsertValidatorElects - struct for bonding transactions
type MsgInsertValidatorElects struct {
	ElectedValidators []ValidatorElect `json:"elected_validators"`
	InitiatorAddr     sdk.ValAddress   `json:"initiator_address"` // Operator address of the validator that initiates this transaction
	CycleNum          sdk.Int          `json:"cycle_number"`
}

// Since public key cannot be put in JSON we need to convert it to/from string
type msgInsertValidatorElectsJSON struct {
	ElectedValidators []ValidatorElect `json:"elected_validators"`
	InitiatorAddr     sdk.ValAddress   `json:"initiator_address"` // Operator address of the validator that initiates this transaction
	CycleNum          sdk.Int          `json:"cycle_number"`
}

func NewMsgInsertValidatorElects(elects []ValidatorElect, initiatorAddr sdk.ValAddress, cycleNum sdk.Int) MsgInsertValidatorElects {
	return MsgInsertValidatorElects{
		ElectedValidators: elects,
		InitiatorAddr:     initiatorAddr,
		CycleNum:          cycleNum,
	}
}

// Return the message type
func (msg MsgInsertValidatorElects) Route() string {
	return RouterKey
}

// Returns a human-readable string for the message, intended for utilization within tags
func (msg MsgInsertValidatorElects) Type() string {
	return "insert_validator_elects"
}

// Simple validation check that doesn't require access to any other information
func (msg MsgInsertValidatorElects) ValidateBasic() sdk.Error {
	// TODO: Check if contract address is empty
	if msg.InitiatorAddr.Empty() {
		return ErrNilValidatorAddress(DefaultCodespace)
	}
	if msg.CycleNum.LTE(sdk.ZeroInt()) || msg.CycleNum.GT(sdk.NewInt(MaxCycles)) {
		return ErrCycleNumberOutOfBounds(DefaultCodespace, MaxCycles)
	}
	if len(msg.ElectedValidators) == 0 || len(msg.ElectedValidators) > MaxValidatorElectsPerCycle {
		return ErrValidatorElectsOutOfBounds(DefaultCodespace, MaxValidatorElectsPerCycle)
	}
	// Check whether any public key or operator address is empty
	for _, elect := range msg.ElectedValidators {
		if len(elect.ConsPubKey.Bytes()) == 0 {
			return ErrNilValidatorElectPubKey(DefaultCodespace)
		}
		if elect.OperatorAddr.Empty() {
			return ErrNilValidatorElectAddress(DefaultCodespace)
		}
		// Check if any of their places is out of bounds
		if elect.Place.LTE(sdk.ZeroInt()) || elect.Place.GT(sdk.NewInt(MaxValidatorElectsPerCycle)) {
			return ErrElectPlaceOutOfBounds(DefaultCodespace, MaxValidatorElectsPerCycle)
		}
	}

	return nil
}

// GetSignBytes returns sorted message bytes to sign over
func (msg MsgInsertValidatorElects) GetSignBytes() []byte {
	bz := msgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign over GetSignBytes()
// Only validator that creates this message needs to sign
func (msg MsgInsertValidatorElects) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.InitiatorAddr)}
}

// MarshalJSON implements the json.Marshaller interface to provide custom JSON serialization of the MsgInsertValidatorElects type
func (msg MsgInsertValidatorElects) MarshalJSON() ([]byte, error) {
	return json.Marshal(msgInsertValidatorElectsJSON{
		ElectedValidators: msg.ElectedValidators,
		InitiatorAddr:     msg.InitiatorAddr,
		CycleNum:          msg.CycleNum,
	})
}

// UnmarshalJSON implements the json.Unmarshaller interface to provide custom JSON deserialization of the MsgInsertValidatorElects type
func (msg *MsgInsertValidatorElects) UnmarshalJSON(bz []byte) error {
	var msgJSON msgInsertValidatorElectsJSON
	if err := json.Unmarshal(bz, &msgJSON); err != nil {
		return err
	}

	msg.ElectedValidators = msgJSON.ElectedValidators
	msg.InitiatorAddr = msgJSON.InitiatorAddr
	msg.CycleNum = msgJSON.CycleNum

	return nil
}
