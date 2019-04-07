package querier

import (
	"fmt"
	"github.com/tendermint/tendermint/crypto"

	"github.com/MVPWorkshop/electionr/x/election/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Converts consensus public keys of elected validators to bech32 format (cosmosvalconspub)
func bech32ifyConsPubs(consPubKeys []crypto.PubKey) (consPubKeysBech32 []string) {
	for _, pubKey := range consPubKeys {
		pubKeyBech32, err := sdk.Bech32ifyConsPub(pubKey)
		if err != nil {
			panic(fmt.Sprintf("Invalid public key %s", pubKey))
		}
		consPubKeysBech32 = append(consPubKeysBech32, pubKeyBech32)
	}
	return consPubKeysBech32
}

// Converts consensus public keys of elected validators to bech32 format (cosmosvalconspub)
func convertValElectsConsPubKeys(elects []types.ValidatorElect) (electsJSON []types.ValidatorElectDisplayJSON) {
	for _, elect := range elects {
		// Convert consensus public key to bech32 format
		pubKey, err := sdk.Bech32ifyConsPub(elect.ConsPubKey)
		if err != nil {
			panic(fmt.Sprintf("Invalid public key of one of validator elects %s", elect.ConsPubKey))
		}
		electJSON := types.NewValidatorElectDisplayJSON(elect.OperatorAddr, pubKey, elect.Place, elect.Left)
		electsJSON = append(electsJSON, electJSON)
	}
	return electsJSON
}
