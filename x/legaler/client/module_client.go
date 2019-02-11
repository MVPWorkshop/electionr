package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"

	legalercmd "github.com/MVPWorkshop/legaler-bc/x/legaler/client/cli"
)

// This abstraction allows clients to import the client functionality from your module in a standard way.
// There is an open issue to add rest functionality to this interface as well

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:   "legaler",
		Short: "Querying commands for the legaler module",
	}
	govQueryCmd.AddCommand(client.GetCommands(
		legalercmd.GetCmdResolveName(mc.storeKey, mc.cdc),
		legalercmd.GetCmdWhoIs(mc.storeKey, mc.cdc),
	)...)
	return govQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:   "legaler",
		Short: "Legaler transactions subcommands",
	}
	govTxCmd.AddCommand(client.PostCommands(
		legalercmd.GetCmdSetName(mc.cdc),
	)...)
	return govTxCmd
}
