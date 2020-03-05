package types

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Protocol/types"
	authexported "github.com/Dipper-Protocol/x/auth/exported"
)

// StakingKeeper defines the expected staking keeper (noalias)
type StakingKeeper interface {
	ApplyAndReturnValidatorSetUpdates(sdk.Context) (updates []abci.ValidatorUpdate)
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, authexported.Account) authexported.Account
	SetAccount(sdk.Context, authexported.Account)
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}

// GenesisAccountsIterator defines the expected iterating genesis accounts object (noalias)
type GenesisAccountsIterator interface {
	IterateGenesisAccounts(
		cdc *codec.Codec,
		appGenesis map[string]json.RawMessage,
		iterateFn func(authexported.Account) (stop bool),
	)
}
