package legaler

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// This is the place to define which queries against application state users will be able to make.

// query endpoints supported by the governance Querier
const (
	QueryResolve = "resolve"
	QueryWhoIs   = "whois"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		// Note: because there isn't an interface similar to Msg for queries,
		// you need to manually define switch statement cases
		// (they can't be pulled off of the query .Route() function)
		switch path[0] {
		case QueryResolve:
			return queryResolve(ctx, path[1:], req, keeper)
		case QueryWhoIs:
			return queryWhoIs(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown legaler query endpoint")
		}
	}
}

// This takes a name and returns the value that is stored by the legaler. This is similar to a DNS query.
func queryResolve(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	name := path[0]
	value := keeper.ResolveName(ctx, name)
	if value == "" {
		return []byte{}, sdk.ErrUnknownRequest("could not resolve name")
	}
	return []byte(value), nil
}

// WhoIs represents a name -> value lookup
type WhoIs struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
}

// This takes a name and returns the value and owner of the name.
func queryWhoIs(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	name := path[0]
	whoIs := WhoIs{}

	whoIs.Value = keeper.ResolveName(ctx, name)
	whoIs.Owner = keeper.GetOwner(ctx, name)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, whoIs)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}
