package keeper

import (
	"fmt"
	"github.com/Dipper-Protocol/x/vm/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Protocol/types"
	"github.com/Dipper-Protocol/x/auth"
	"github.com/Dipper-Protocol/x/params"
)

type Keeper struct {
	Cdc        *codec.Codec
	paramstore params.Subspace
	StateDB    *types.CommitStateDB
}

func NewKeeper(cdc *codec.Codec, storeKey, codeKey, logKey, storageDebugKey sdk.StoreKey, paramstore params.Subspace, ak auth.AccountKeeper) Keeper {
	return Keeper{
		Cdc:        cdc,
		paramstore: paramstore.WithKeyTable(ParamKeyTable()),
		StateDB:    types.NewCommitStateDB(ak, storeKey, codeKey, logKey, storageDebugKey),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

func (k Keeper) GetState(ctx sdk.Context, addr sdk.AccAddress, hash sdk.Hash) sdk.Hash {
	return k.StateDB.WithContext(ctx).GetState(addr, hash)
}

func (k *Keeper) GetCode(ctx sdk.Context, addr sdk.AccAddress) []byte {
	return k.StateDB.WithContext(ctx).GetCode(addr)
}

func (k *Keeper) GetLogs(ctx sdk.Context, hash sdk.Hash) []*types.Log {
	return k.StateDB.WithContext(ctx).GetLogs(hash)
}
