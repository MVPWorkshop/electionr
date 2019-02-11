package rest

import (
	"fmt"
	"github.com/MVPWorkshop/legaler-bc/x/legaler"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const restName = "name"

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), setNameHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}", storeName, restName), resolveNameHandler(cdc, cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/whois", storeName, restName), whoIsNameHandler(cdc, cliCtx, storeName)).Methods("GET")
}

// Define query handlers

func resolveNameHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/resolve/%s", storeName, paramType), nil)
		if err != nil {
			utils.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
			return
		}

		utils.PostProcessResponse(writer, cdc, res, cliCtx.Indent)
	}
}

func whoIsNameHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", storeName, paramType), nil)
		if err != nil {
			utils.WriteErrorResponse(writer, http.StatusNotFound, err.Error())
			return
		}

		utils.PostProcessResponse(writer, cdc, res, cliCtx.Indent)
	}
}

type setNameReq struct {
	BaseReq utils.BaseReq	`json:"base_req"`
	Name string				`json:"name"`
	Value string			`json:"value"`
	Owner string			`json:"owner"`
}

func setNameHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req setNameReq

		err := utils.ReadRESTReq(writer, request, cdc, &req)
		if err != nil {
			utils.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(writer) {
			return
		}
		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			utils.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}

		// Create the message
		msg := legaler.NewMsgSetName(req.Name, req.Value, addr)
		err = msg.ValidateBasic()
		if err != nil {
			utils.WriteErrorResponse(writer, http.StatusBadRequest, err.Error())
			return
		}
		utils.CompleteAndBroadcastTxREST(writer, request, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}