package types

import "github.com/cosmos/cosmos-sdk/codec"

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgInsertValidatorElects{}, "cosmos-sdk/MsgInsertValidatorElects", nil)
}

// Generic sealed codec to be used throughout sdk
var MsgCdc *codec.Codec
