package legaler

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "legaler" type messages.
// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
// At the moment, there is only one Msg/Handler.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized legaler Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSetName
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg MsgSetName) sdk.Result {
	// Set the name to the value specified in the msg.
	keeper.SetName(ctx, msg.Name, msg.Value)
	// Set the owner of this name.
	keeper.SetOwner(ctx, msg.Name, msg.Owner)
	return sdk.Result{}
}
