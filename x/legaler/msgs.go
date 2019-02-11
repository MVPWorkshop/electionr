package legaler

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgSetName defines a SetName message
type MsgSetName struct {
	Name  string
	Value string
	Owner sdk.AccAddress
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}

// Return the message type.
// Must be alphanumeric or empty.
// Type should return the name of the module
func (msg MsgSetName) Route() string {
	return "legaler"
}

// Returns a human-readable string for the message,
// intended for utilization within tags
// Name should return the action
func (msg MsgSetName) Type() string {
	return "set_name"
}

// The above functions are used by the SDK to route Msgs to the proper module for handling.
// They also add human readable names to database tags used for indexing.

// ValidateBasic does a simple validation check that
// doesn't require access to any other information.
// ValidateBasic Implements Msg.
// ValidateBasic is used to provide some basic stateless checks on the validity of the Msg.
// In this case, check that none of the attributes are empty.
func (msg MsgSetName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Value cannot be empty")
	}
	return nil
}

// Get the canonical byte representation of the Msg.
// GetSignBytes Implements Msg.
// GetSignBytes defines how the Msg gets encoded for signing.
// In most cases this means marshal to sorted JSON.
// The output should not be modified.
func (msg MsgSetName) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Signers returns the addrs of signers that must sign.
// CONTRACT: All signatures must be present to be valid.
// CONTRACT: Returns addrs in some deterministic order.
// GetSigners Implements Msg.
// GetSigners defines whose signature is required on a Tx in order for it to be valid.
// In this case, for example, the MsgSetName requires that the Owner signs the transaction
// when trying to reset what the name points to.
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
