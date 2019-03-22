package rest

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/MVPWorkshop/legaler-bc/x/election"
	"github.com/MVPWorkshop/legaler-bc/x/election/querier"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get info about certain election cycle
	r.HandleFunc(
		"/election/cycle/{cycleNum}",
		cycleInfoHandlerFn(cliCtx, cdc),
	).Methods("GET")
}

func cycleInfoHandlerFn(cliCtx context.CLIContext, cdc *codec.Codec) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// Get parameters (cycle number)
		vars := mux.Vars(request)
		cycleNum := vars["cycleNum"]
		cycleNumInt, ok := sdk.NewIntFromString(cycleNum)
		if !ok {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, election.ErrCycleNumInvalid(election.DefaultCodespace).Error())
			return
		}
		// Check if cycle number is outside of bounds
		if cycleNumInt.LTE(sdk.ZeroInt()) || cycleNumInt.GT(sdk.NewInt(election.MaxCycles)) {
			rest.WriteErrorResponse(
				writer,
				http.StatusBadRequest,
				election.ErrCycleNumberOutOfBounds(election.DefaultCodespace, election.MaxCycles).Error(),
			)
			return
		}

		// Prepare the data
		params := querier.NewQueryCyclesParams(cycleNumInt)
		bz, err := cdc.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		// Fetch the data
		res, err := cliCtx.QueryWithData("custom/election/cycles", bz)
		if err != nil {
			rest.WriteErrorResponse(writer, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(writer, cdc, res, cliCtx.Indent)
	}
}
