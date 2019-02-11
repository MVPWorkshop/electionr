package legaler

import "github.com/cosmos/cosmos-sdk/codec"

// Any interface you create and any struct that implements an interface
// needs to be declared in the RegisterCodec function.
// RegisterCodec registers concrete types on wire codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSetName{}, "legaler/SetName", nil)
}
