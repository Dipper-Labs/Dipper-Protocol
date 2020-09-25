package bank_test

import (
	"os"
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

const (
	dbPath = "/tmp/dipbenchmarktest"
)

func Benchmark_handleMsgSend(b *testing.B) {
	os.RemoveAll(dbPath)
	defer os.RemoveAll(dbPath)

	cdc := app.MakeLatestCodec()

	// setup params
	paramsKeeper := params.NewKeeper(cdc, protocol.Keys[params.StoreKey], protocol.TKeys[params.TStoreKey])
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := paramsKeeper.Subspace(bank.DefaultParamspace)

	// setup keepers
	accountKeeper := auth.NewAccountKeeper(cdc, protocol.Keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, bankSubspace, nil)

	// prepare db
	db, err := sdk.NewLevelDB("application", dbPath)
	require.Nil(b, err)
	ms := store.NewCommitMultiStore(db)
	ms.SetPruning(store.PruneSyncable)
	ms.MountStoreWithDB(protocol.Keys[params.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.Keys[auth.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.TKeys[params.TStoreKey], sdk.StoreTypeTransient, nil)
	ms.LoadLatestVersion()

	// setup context
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)

	// enable bank sending
	bankKeeper.SetSendEnabled(ctx, true)

	// generate account and setup accountKeeper
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

	// prepare send msgs
	sendAmount := sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(sdk.NativeTokenFraction*100)))
	var msgs []bank.MsgSend
	for i := 0; i < b.N; i++ {
		msgSend := bank.NewMsgSend(addrs[i], addrs[i+1], sendAmount)
		msgs = append(msgs, msgSend)
	}

	// reset benchmark timer
	b.ResetTimer()

	// ban sending benchmark
	for i := 0; i < b.N; i++ {
		bank.HandleMsgSend(ctx, bankKeeper, msgs[i])
	}

	ms.Commit()
}
