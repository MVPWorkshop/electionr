package app

import (
	"fmt"
	"sort"
	"time"

	"github.com/MVPWorkshop/electionr/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
)

// State to Unmarshal
type GenesisState struct {
	Accounts     []GenesisAccount      `json:"accounts"`
	AuthData     auth.GenesisState     `json:"auth"`
	BankData     bank.GenesisState     `json:"bank"`
	StakingData  staking.GenesisState  `json:"staking"`
	MintData     mint.GenesisState     `json:"mint"`
	DistrData    distr.GenesisState    `json:"distr"`
	SlashingData slashing.GenesisState `json:"slashing"`
}

func NewGenesisState(accounts []GenesisAccount, authData auth.GenesisState,
	bankData bank.GenesisState,
	stakingData staking.GenesisState, mintData mint.GenesisState,
	distrData distr.GenesisState,
	slashingData slashing.GenesisState) GenesisState {

	return GenesisState{
		Accounts:     accounts,
		AuthData:     authData,
		BankData:     bankData,
		StakingData:  stakingData,
		MintData:     mintData,
		DistrData:    distrData,
		SlashingData: slashingData,
	}
}

// Sanitize sorts accounts and coin sets.
func (gs GenesisState) Sanitize() {
	sort.Slice(gs.Accounts, func(i, j int) bool {
		return gs.Accounts[i].AccountNumber < gs.Accounts[j].AccountNumber
	})

	for _, acc := range gs.Accounts {
		acc.Coins = acc.Coins.Sort()
	}
}

// GenesisAccount defines an account initialized at genesis.
type GenesisAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	Sequence      uint64         `json:"sequence_number"`
	AccountNumber uint64         `json:"account_number"`

	// vesting account fields
	OriginalVesting  sdk.Coins `json:"original_vesting"`  // total vesting coins upon initialization
	DelegatedFree    sdk.Coins `json:"delegated_free"`    // delegated vested coins at time of delegation
	DelegatedVesting sdk.Coins `json:"delegated_vesting"` // delegated vesting coins at time of delegation
	StartTime        int64     `json:"start_time"`        // vesting start time (UNIX Epoch time)
	EndTime          int64     `json:"end_time"`          // vesting end time (UNIX Epoch time)
}

func NewGenesisAccount(acc auth.Account) GenesisAccount {
	gacc := GenesisAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	vacc, ok := acc.(auth.VestingAccount)
	if ok {
		gacc.OriginalVesting = vacc.GetOriginalVesting()
		gacc.DelegatedFree = vacc.GetDelegatedFree()
		gacc.DelegatedVesting = vacc.GetDelegatedVesting()
		gacc.StartTime = vacc.GetStartTime()
		gacc.EndTime = vacc.GetEndTime()
	}

	return gacc
}

// convert GenesisAccount to auth.BaseAccount
func (ga *GenesisAccount) ToAccount() auth.Account {
	bacc := &auth.BaseAccount{
		Address:       ga.Address,
		Coins:         ga.Coins.Sort(),
		AccountNumber: ga.AccountNumber,
		Sequence:      ga.Sequence,
	}

	if !ga.OriginalVesting.IsZero() {
		baseVestingAcc := &auth.BaseVestingAccount{
			BaseAccount:      bacc,
			OriginalVesting:  ga.OriginalVesting,
			DelegatedFree:    ga.DelegatedFree,
			DelegatedVesting: ga.DelegatedVesting,
			EndTime:          ga.EndTime,
		}

		if ga.StartTime != 0 && ga.EndTime != 0 {
			return &auth.ContinuousVestingAccount{
				BaseVestingAccount: baseVestingAcc,
				StartTime:          ga.StartTime,
			}
		} else if ga.EndTime != 0 {
			return &auth.DelayedVestingAccount{
				BaseVestingAccount: baseVestingAcc,
			}
		} else {
			panic(fmt.Sprintf("invalid genesis vesting account: %+v", ga))
		}
	}

	return bacc
}

// NewDefaultGenesisState generates the default state for electionr.
func NewDefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts:     nil,
		AuthData:     auth.DefaultGenesisState(),
		BankData:     bank.DefaultGenesisState(),
		StakingData:  staking.DefaultGenesisState(),
		MintData:     mint.DefaultGenesisState(),
		DistrData:    distr.DefaultGenesisState(),
		SlashingData: slashing.DefaultGenesisState(),
	}
}

// ElectionrValidateGenesisState ensures that the genesis state obeys the expected invariants
// TODO: No validators are both bonded and jailed (#2088)
// TODO: Error if there is a duplicate validator (#1708)
// TODO: Ensure all state machine parameters are in genesis (#1704)
func ElectionrValidateGenesisState(genesisState GenesisState) error {
	if err := validateGenesisStateAccounts(genesisState.Accounts); err != nil {
		return err
	}

	if err := auth.ValidateGenesis(genesisState.AuthData); err != nil {
		return err
	}
	if err := bank.ValidateGenesis(genesisState.BankData); err != nil {
		return err
	}
	if err := staking.ValidateGenesis(genesisState.StakingData); err != nil {
		return err
	}
	if err := mint.ValidateGenesis(genesisState.MintData); err != nil {
		return err
	}
	if err := distr.ValidateGenesis(genesisState.DistrData); err != nil {
		return err
	}

	return slashing.ValidateGenesis(genesisState.SlashingData)
}

// validateGenesisStateAccounts performs validation of genesis accounts. It
// ensures that there are no duplicate accounts in the genesis state and any
// provided vesting accounts are valid.
func validateGenesisStateAccounts(accs []GenesisAccount) error {
	addrMap := make(map[string]bool, len(accs))
	for _, acc := range accs {
		addrStr := acc.Address.String()

		// disallow any duplicate accounts
		if _, ok := addrMap[addrStr]; ok {
			return fmt.Errorf("duplicate account found in genesis state; address: %s", addrStr)
		}

		// validate any vesting fields
		if !acc.OriginalVesting.IsZero() {
			if acc.EndTime == 0 {
				return fmt.Errorf("missing end time for vesting account; address: %s", addrStr)
			}

			if acc.StartTime >= acc.EndTime {
				return fmt.Errorf(
					"vesting start time must before end time; address: %s, start: %s, end: %s",
					addrStr,
					time.Unix(acc.StartTime, 0).UTC().Format(time.RFC3339),
					time.Unix(acc.EndTime, 0).UTC().Format(time.RFC3339),
				)
			}
		}

		addrMap[addrStr] = true
	}

	return nil
}
