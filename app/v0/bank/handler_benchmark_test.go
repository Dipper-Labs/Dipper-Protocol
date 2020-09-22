package bank_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/Dipper-Labs/Dipper-Protocol/app"
	"github.com/Dipper-Labs/Dipper-Protocol/app/protocol"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/bank"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/params"
	"github.com/Dipper-Labs/Dipper-Protocol/store"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

func Benchmark_handleMsgSend(b *testing.B) {
	cdc := app.MakeLatestCodec()

	paramsKeeper := params.NewKeeper(cdc, protocol.Keys[params.StoreKey], protocol.TKeys[params.TStoreKey])
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := paramsKeeper.Subspace(bank.DefaultParamspace)

	accountKeeper := auth.NewAccountKeeper(cdc, protocol.Keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, bankSubspace, nil)

	dataDir := filepath.Join("/tmp", "testdata")
	defer os.RemoveAll(dataDir)
	db, err := sdk.NewLevelDB("application", dataDir)
	require.Nil(b, err)
	ms := store.NewCommitMultiStore(db)
	ms.SetPruning(store.PruneSyncable)
	ms.MountStoreWithDB(protocol.Keys[params.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.Keys[auth.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.TKeys[params.TStoreKey], sdk.StoreTypeTransient, nil)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)
	bankKeeper.SetSendEnabled(ctx, true)

	var addrs []sdk.AccAddress
	coins := sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(sdk.NativeTokenFraction*1000)))
	for i := 0; i < b.N+1; i++ {
		privateKey := ed25519.GenPrivKey()
		publicKey := privateKey.PubKey()
		address := publicKey.Address()
		addr, err := sdk.AccAddressFromHex(address.String())
		if err != nil {
			b.Fatal(err)
		}

		addrs = append(addrs, addr)
		acc := accountKeeper.NewAccountWithAddress(ctx, addr)
		acc.SetCoins(coins)
		accountKeeper.SetAccount(ctx, acc)
	}

	sendAmount := sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(sdk.NativeTokenFraction*100)))
	var msgs []bank.MsgSend
	for i := 0; i < b.N; i++ {
		msgSend := bank.NewMsgSend(addrs[i], addrs[i+1], sendAmount)
		msgs = append(msgs, msgSend)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bank.HandleMsgSend(ctx, bankKeeper, msgs[i])
	}

	ms.Commit()
}
