package main

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	app "github.com/MVPWorkshop/legaler-bc"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "legalerd",
		Short:             "Gaia Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	rootCmd.AddCommand(gaiaInit.InitCmd(ctx, cdc))
	rootCmd.AddCommand(gaiaInit.CollectGenTxsCmd(ctx, cdc))
	rootCmd.AddCommand(gaiaInit.TestnetFilesCmd(ctx, cdc))
	rootCmd.AddCommand(gaiaInit.GenTxCmd(ctx, cdc))
	rootCmd.AddCommand(gaiaInit.AddGenesisAccountCmd(ctx, cdc))
	rootCmd.AddCommand(gaiaInit.ValidateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(client.NewCompletionCmd(rootCmd, true))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "GA", app.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewLegalerApp(
		logger, db, traceStore, true,
		baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	if height != -1 {
		legApp := app.NewLegalerApp(logger, db, traceStore, false)
		err := legApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return legApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}
	legApp := app.NewLegalerApp(logger, db, traceStore, true)
	return legApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
