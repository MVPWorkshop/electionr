package querier

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/MVPWorkshop/electionr/x/election/keeper"
	"github.com/MVPWorkshop/electionr/x/election/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Query endpoints supported by the election Querier
const (
	QueryCycles = "cycles"
)

// Creates a module level router for state queries
func NewQuerier(k keeper.Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryCycles:
			return queryCycles(ctx, cdc, req, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown election query endpoint")
		}
	}
}

// Defines the parameters for the following queries:
// - 'custom/election/cycles
type QueryCyclesParams struct {
	CycleNum sdk.Int
}

func NewQueryCyclesParams(cycleNum sdk.Int) QueryCyclesParams {
	return QueryCyclesParams{
		CycleNum: cycleNum,
	}
}

func queryCycles(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keeper.Keeper) (res []byte, err sdk.Error) {
	var params QueryCyclesParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, types.ErrCycleNotFound(types.DefaultCodespace)
	}

	// Get cycles from state
	cycles := k.GetCyclesByCycleNum(ctx, params.CycleNum)

	res, errRes = codec.MarshalJSONIndent(cdc, cycles)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}
