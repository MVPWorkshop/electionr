package main

import (
	"encoding/json"
	"github.com/MVPWorkshop/legaler-bc"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
var DefaultNodeHome = os.ExpandEnv("$HOME/.led")

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "legalerd",
		Short:             "legaler App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	server.AddCommands(ctx, cdc, rootCmd, newApp, appExporter())

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "LE", DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		// handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return app.NewLegalerApp(logger, db)
}

func appExporter() server.AppExporter {
	return func(logger log.Logger, db dbm.DB, _ io.Writer, _ int64, _ bool, _ []string) (json.RawMessage, []tmtypes.GenesisValidator, error) {
		dapp := app.NewLegalerApp(logger, db)
		return dapp.ExportAppStateAndValidators()
	}
}
