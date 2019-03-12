package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// Ensure Msg interface compliance at compile time
// (It implements Msg interface)
var _ sdk.Msg = &MsgInsertValidatorElects{}

//______________________________________________________________________

// MsgInsertValidatorElects - struct for bonding transactions
type MsgInsertValidatorElects struct {
	PubKeys      []crypto.PubKey  `json:"pub_keys"`
	ValAddresses []sdk.ValAddress `json:"validator_addresses"` // Validator operator addresses
	ValAddress   sdk.ValAddress   `json:"validator_address"`   // Operator address of the validator that initiated this transaction
	CycleNum     sdk.Int          `json:"cycle_number"`
}

// Since public key cannot be put in JSON we need to convert it to/from string
type msgInsertValidatorElectsJSON struct {
	PubKeys      []string         `json:"pub_keys"`
	ValAddresses []sdk.ValAddress `json:"validator_addresses"` // Validator operator addresses
	ValAddress   sdk.ValAddress   `json:"validator_address"`   // Operator address of the validator that initiated this transaction
	CycleNum     sdk.Int          `json:"cycle_number"`
}

func NewMsgInsertValidatorElects(pubKeys []crypto.PubKey, valAddrs []sdk.ValAddress, valAddr sdk.ValAddress, cycleNum sdk.Int) MsgInsertValidatorElects {
	return MsgInsertValidatorElects{
		PubKeys:      pubKeys,
		ValAddresses: valAddrs,
		ValAddress:   valAddr,
		CycleNum:     cycleNum,
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
	if msg.ValAddress.Empty() {
		return ErrNilValidatorAddress(DefaultCodespace)
	}
	if msg.CycleNum.LTE(sdk.ZeroInt()) || msg.CycleNum.GT(sdk.NewInt(MaxCycles)) {
		return ErrCycleNumberOutOfBounds(DefaultCodespace, MaxCycles)
	}
	if len(msg.PubKeys) == 0 || len(msg.PubKeys) > MaxValidatorElectsPerCycle {
		return ErrPublicKeysOutOfBounds(DefaultCodespace, MaxValidatorElectsPerCycle)
	}
	if len(msg.ValAddresses) == 0 || len(msg.ValAddresses) > MaxValidatorElectsPerCycle {
		return ErrValidatorAddressesOutOfBounds(DefaultCodespace, MaxValidatorElectsPerCycle)
	}
	// Check whether any validator operator address is empty
	for _, address := range msg.ValAddresses {
		if address.Empty() {
			return ErrNilValidatorElectorAddress(DefaultCodespace)
		}
	}
	// Number of public keys should be equal to number of validator operator addresses
	if len(msg.PubKeys) != len(msg.ValAddresses) {
		return ErrPubKeysValAddressesMissmatch(DefaultCodespace)
	}

	return nil
}

// GetSignBytes returns sorted message bytes to sign over
func (msg MsgInsertValidatorElects) GetSignBytes() []byte {
	bz := MsgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners returns the addresses of signers that must sign over GetSignBytes()
// Only validator that creates this message needs to sign
func (msg MsgInsertValidatorElects) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValAddress)}
}

// MarshalJSON implements the json.Marshaler interface to provide custom JSON serialization of the MsgInsertValidatorElects type
func (msg MsgInsertValidatorElects) MarshalJSON() ([]byte, error) {
	pubKeys := make([]string, MaxValidatorElectsPerCycle)
	for _, pubKey := range msg.PubKeys {
		pubKeys = append(pubKeys, sdk.MustBech32ifyConsPub(pubKey))
	}
	return json.Marshal(msgInsertValidatorElectsJSON{
		PubKeys:      pubKeys,
		ValAddresses: msg.ValAddresses,
		ValAddress:   msg.ValAddress,
		CycleNum:     msg.CycleNum,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface to provide custom JSON deserialization of the MsgInsertValidatorElects type
func (msg *MsgInsertValidatorElects) UnmarshalJSON(bz []byte) error {
	var msgJSON msgInsertValidatorElectsJSON
	if err := json.Unmarshal(bz, &msgJSON); err != nil {
		return err
	}

	msg.PubKeys = make([]crypto.PubKey, MaxValidatorElectsPerCycle)
	for _, pubKey := range msgJSON.PubKeys {
		msg.PubKeys = append(msg.PubKeys, sdk.MustGetConsPubKeyBech32(pubKey))
	}
	msg.ValAddresses = msgJSON.ValAddresses
	msg.ValAddress = msgJSON.ValAddress
	msg.CycleNum = msgJSON.CycleNum

	return nil
}
