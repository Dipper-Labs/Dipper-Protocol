package keeper

// DONTCOVER

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/Dipper-Protocol/codec"
	"github.com/Dipper-Protocol/store"
	sdk "github.com/Dipper-Protocol/types"
	"github.com/Dipper-Protocol/x/auth"
	"github.com/Dipper-Protocol/x/bank/internal/types"
	"github.com/Dipper-Protocol/x/params"
)

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context
	k   Keeper
	ak  auth.AccountKeeper
	pk  params.Keeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := codec.New()
	auth.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	authCapKey := sdk.NewKVStoreKey("authCapKey")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.AccAddress([]byte("moduleAcc")).String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	ak := auth.NewAccountKeeper(
		cdc, authCapKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	ak.SetParams(ctx, auth.DefaultParams())

	bankKeeper := NewBaseKeeper(ak, pk.Subspace(types.DefaultParamspace), types.DefaultCodespace, blacklistedAddrs)
	bankKeeper.SetSendEnabled(ctx, true)

	return testInput{cdc: cdc, ctx: ctx, k: bankKeeper, ak: ak, pk: pk}
}
