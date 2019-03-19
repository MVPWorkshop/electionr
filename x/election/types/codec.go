package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgInsertValidatorElects{}, "legaler/MsgInsertValidatorElects", nil)
}

// generic sealed codec to be used throughout sdk
var msgCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	// Register crypto types like crypto.PubKey
	codec.RegisterCrypto(cdc)
	// Seal the codec so nothing new can be registered
	msgCdc = cdc.Seal()
}
