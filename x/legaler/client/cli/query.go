package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// Define cobra.Commands for each of your modules Queriers (resolve, and whoIs)

// GetCmdResolveName queries information about a name
func GetCmdResolveName(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "resolve [name]",
		Short: "resolve name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// CLIContext carries data about user input and application configuration that are needed for CLI interactions
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			// The first part of the path is used to differentiate the types of queries possible to SDK applications: custom is for Queriers
			// The second piece (legaler) is the name of the module to route the query to.
			// Third part is the specific querier in the module that will be called.
			// In this example the fourth piece is the query. This works because the query parameter is a simple string.
			// To enable more complex query inputs you need to use the second argument of the .QueryWithData() function to pass in data.
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/resolve/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not resolve name - %s \n", string(name))
				return nil
			}
			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdWhoIs queries information about a domain
func GetCmdWhoIs(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "whois [name]",
		Short: "Query whois info of name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("could not resolve whois - %s \n", string(name))
				return nil
			}
			fmt.Println(string(res))
			return nil
		},
	}
}
