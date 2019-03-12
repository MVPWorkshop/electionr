package rest

import (
	"github.com/MVPWorkshop/legaler-bc/x/election/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"github.com/tendermint/tendermint/crypto"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/election/validator_elects",
		postInsertValidatorElectsHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type msgInsertValidatorElectsInput struct {
	BaseReq      rest.BaseReq     `json:"base_req"`
	PubKeys      []string         `json:"pub_keys"`
	ValAddresses []sdk.ValAddress `json:"validator_addresses"` // Validator operator addresses
	CycleNum     sdk.Int          `json:"cycle_num"`
}

// TODO: Probaj samo sa nizom adresa bez pubkeys
func postInsertValidatorElectsHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req msgInsertValidatorElectsInput

		// Get parameters
		if !rest.ReadRESTReq(writer, request, cdc, &req) {
			return
		}

		// Trim string parameters of base request object
		req.BaseReq = req.BaseReq.Sanitize()
		// Perform base request validation
		if !req.BaseReq.ValidateBasic(writer) {
			return
		}

		// Derive account address from either bech32 address or Keybase name
		// We'll need it for validator's operator address
		fromAddress, fromName, err := context.GetFromFields(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		// Get validator consensus public keys
		pubKeys := make([]crypto.PubKey, types.MaxValidatorElectsPerCycle)
		for _, pubKey := range req.PubKeys {
			valPubKey, err := sdk.GetConsPubKeyBech32(pubKey)
			if err != nil {
				rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
				return
			}
			pubKeys = append(pubKeys, valPubKey)
		}

		msg := types.NewMsgInsertValidatorElects(pubKeys, req.ValAddresses, sdk.ValAddress(fromAddress), req.CycleNum)

		// Perform basic request validation
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		// Check whether request was to generate transaction without signing and broadcasting
		if req.BaseReq.GenerateOnly {
			rest.WriteGenerateStdTxResponse(writer, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
			return
		}

		// Upgrade context's account name and address
		cliCtx = cliCtx.WithFromName(fromName).WithFromAddress(fromAddress)

		// Broadcast message to other nodes
		rest.CompleteAndBroadcastTxREST(writer, request, cliCtx, req.BaseReq, []sdk.Msg{msg}, cdc)
	}
}
