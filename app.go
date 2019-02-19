package app

import (
	"encoding/json"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	"sort"
)

const appName = "legaler"

type legalerApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyStaking       *sdk.KVStoreKey
	tkeyStaking      *sdk.TransientStoreKey
	keyMint          *sdk.KVStoreKey
	keyDistribution  *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tkeyParams       *sdk.TransientStoreKey

	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	paramsKeeper        params.Keeper
	bankKeeper          bank.Keeper
	mintKeeper          mint.Keeper
	distributionKeeper  distribution.Keeper
	stakingKeeper       staking.Keeper
}

func NewLegalerApp(logger log.Logger, db dbm.DB) *legalerApp {
	// First define the top level codec that will be shared by the different modules.
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	// Here you initialize your application with the store keys it requires
	var app = &legalerApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:          sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:       sdk.NewKVStoreKey(auth.StoreKey),
		keyFeeCollection: sdk.NewKVStoreKey(auth.FeeStoreKey),
		keyStaking:       sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      sdk.NewTransientStoreKey(staking.TStoreKey),
		keyMint:          sdk.NewKVStoreKey(mint.StoreKey),
		keyDistribution:  sdk.NewKVStoreKey(distribution.StoreKey),
		keyParams:        sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:       sdk.NewTransientStoreKey(params.TStoreKey),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	// The FeeCollectionKeeper collects transaction fees and renders them to the fee distribution module
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(cdc, app.keyFeeCollection)

	// The BankKeeper allows you to perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	stakeKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking,
		app.tkeyStaking,
		app.bankKeeper,
		app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)

	app.mintKeeper = mint.NewKeeper(
		app.cdc,
		app.keyMint,
		app.paramsKeeper.Subspace(mint.DefaultParamspace),
		&stakeKeeper,
		app.feeCollectionKeeper,
	)

	app.distributionKeeper = distribution.NewKeeper(
		app.cdc,
		app.keyDistribution,
		app.paramsKeeper.Subspace(distribution.DefaultParamspace),
		app.bankKeeper,
		&stakeKeeper,
		app.feeCollectionKeeper,
		distribution.DefaultCodespace,
	)

	app.stakingKeeper = stakeKeeper

	// The app.Router is the main transaction router where each module registers its routes
	app.Router().
		AddRoute(bank.RouterKey, bank.NewHandler(app.bankKeeper)).
		AddRoute(distribution.RouterKey, distribution.NewHandler(app.distributionKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper))

	// Register query routes
	app.QueryRouter().
		AddRoute(distribution.QuerierRoute, distribution.NewQuerier(app.distributionKeeper)).
		AddRoute(staking.QuerierRoute, staking.NewQuerier(app.stakingKeeper, app.cdc))

	// Initialize BaseApp
	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keyFeeCollection,
		app.keyStaking,
		app.tkeyStaking,
		app.keyMint,
		app.keyDistribution,
		app.keyParams,
		app.tkeyParams,
	)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.initChainer)
	//app.SetBeginBlocker(app.BeginBlocker)
	//app.SetEndBlocker(app.EndBlocker)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState struct {
	AuthData         auth.GenesisState         `json:"auth"`
	Accounts         []*auth.BaseAccount       `json:"accounts"`
	StakeData        staking.GenesisState      `json:"stake"`
	MintData         mint.GenesisState         `json:"mint"`
	DistributionData distribution.GenesisState `json:"distr"`
}

// The initChainer defines how accounts in genesis.json are mapped into the application state on initial chain start.
func (app *legalerApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	validators := app.initFromGenesisState(ctx, *genesisState)

	return abci.ResponseInitChain{Validators: validators}
}

func (app *legalerApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	mint.BeginBlocker(ctx, app.mintKeeper)
	distribution.BeginBlocker(ctx, req, app.distributionKeeper)
	return abci.ResponseBeginBlock{}
}

func (app *legalerApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	validatorUpdates, tags := staking.EndBlocker(ctx, app.stakingKeeper)

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}

// Inistialize store from a genesis state
func (app *legalerApp) initFromGenesisState(ctx sdk.Context, genesisState GenesisState) []abci.ValidatorUpdate {
	// Sort by account number to maintain consistency
	sort.Slice(genesisState.Accounts, func(i, j int) bool {
		return genesisState.Accounts[i].AccountNumber < genesisState.Accounts[j].AccountNumber
	})

	// Load the accounts (BaseAccount for now)
	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	// Load the initial stake information
	validators, err := staking.InitGenesis(ctx, app.stakingKeeper, genesisState.StakeData)
	if err != nil {
		panic(err) // TODO: Do this without panic
	}

	// Initialize module-specific stores
	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	mint.InitGenesis(ctx, app.mintKeeper, genesisState.MintData)
	distribution.InitGenesis(ctx, app.distributionKeeper, genesisState.DistributionData)

	return validators
}

// The ExportAppStateAndValidators function helps bootstrap the initial state for application.
func (app *legalerApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*auth.BaseAccount{}

	appendAccountsFn := func(acc auth.Account) bool {
		account := &auth.BaseAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}
		accounts = append(accounts, account)
		return false
	}

	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	genState := GenesisState{
		Accounts: accounts,
		AuthData: auth.DefaultGenesisState(),
	}

	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
