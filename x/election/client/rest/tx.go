package rest

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/MVPWorkshop/electionr/x/election"
	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/election/validator_elects",
		postInsertValidatorElectsHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type msgInsertValidatorElectsInput struct {
	BaseReq           rest.BaseReq                  `json:"base_req"`
	ElectedValidators []election.ValidatorElectJSON `json:"elected_validators"`
	CycleNum          sdk.Int                       `json:"cycle_number"`
}

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

		elects := make([]election.ValidatorElect, 0, election.MaxValidatorElectsPerCycle)
		for _, elect := range req.ElectedValidators {
			consPubKey, err := sdk.GetConsPubKeyBech32(elect.ConsPubKey)
			if err != nil {
				return
			}
			valElect := election.NewValidatorElect(elect.OperatorAddr, consPubKey, elect.Place)
			elects = append(elects, valElect)
		}
		msg := election.NewMsgInsertValidatorElects(
			elects,
			sdk.ValAddress(fromAddress),
			req.CycleNum,
		)

		// Perform basic request validation
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		// Upgrade context's account name and address
		cliCtx = cliCtx.WithFromName(fromName).WithFromAddress(fromAddress)

		// Return message for signing
		clientrest.WriteGenerateStdTxResponse(writer, cdc, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
