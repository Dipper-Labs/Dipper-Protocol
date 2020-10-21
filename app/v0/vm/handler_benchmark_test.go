package vm_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Dipper-Labs/Dipper-Protocol/app"
	"github.com/Dipper-Labs/Dipper-Protocol/app/protocol"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/params"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/vm"
	"github.com/Dipper-Labs/Dipper-Protocol/store"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

const (
	dbPath             = "/Users/sun/testdata"
	contractAbiPath    = "./benchmark_test_files/contract_abi"
	contractCodePath   = "./benchmark_test_files/contract_code"
	contractCallMethod = "ipalClaim"
	vmUpdateMethod     = "ipalUpdate"
	hashTest           = "hashTest"
)

func translateSize(size int) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%dByte", size)
	case size < 1024*1024:
		return fmt.Sprintf("%dK%dByte", size/1024, size%1024)
	default:
		return fmt.Sprintf("%dM%dK%dByte", size/1024/1024, (size/1024)%1024, size%1024)
	}
}

func prepareTest(b *testing.B, dbPath string) (k vm.Keeper, addrs []sdk.AccAddress, ms store.CommitMultiStore, ctx sdk.Context) {
	cdc := app.MakeLatestCodec()

	// params
	paramsKeeper := params.NewKeeper(cdc, protocol.Keys[params.StoreKey], protocol.TKeys[params.TStoreKey])
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	vmSubspace := paramsKeeper.Subspace(vm.DefaultParamspace)

	// keeper
	accountKeeper := auth.NewAccountKeeper(cdc, protocol.Keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	vmKeeper := vm.NewKeeper(cdc, protocol.Keys[vm.StoreKey], vmSubspace, accountKeeper)

	// new db
	db, err := sdk.NewLevelDB("application", dbPath)
	require.Nil(b, err)

	// new store
	ms = store.NewCommitMultiStore(db)
	ms.SetPruning(store.PruneSyncable)

	// mount stores
	ms.MountStoreWithDB(protocol.Keys[params.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.Keys[auth.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.Keys[vm.StoreKey], sdk.StoreTypeIAVL, nil)

	// mount tstores
	ms.MountStoreWithDB(protocol.TKeys[params.TStoreKey], sdk.StoreTypeTransient, nil)

	// load store
	ms.LoadLatestVersion()

	// new context
	logger := log.NewNopLogger()
	//logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	ctx = sdk.NewContext(ms, abci.Header{}, false, nil).WithLogger(logger).WithGasMeter(sdk.NewGasMeter(1000000000000000))

	// setup vm keeper params
	params := vm.DefaultParams()
	vmKeeper.SetParams(ctx, params)

	// generate accounts
	coins := sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(sdk.NativeTokenFraction*1000)))
	for i := 0; i < b.N; i++ {
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

	return vmKeeper, addrs, ms, ctx
}

func createContracts(b *testing.B, k vm.Keeper, addrs []sdk.AccAddress, ctx sdk.Context) (contractAddrs []sdk.AccAddress) {
	// load contract code
	code, err := vm.CodeFromFile(contractCodePath)
	require.Nil(b, err)

	// genenrate contract create msgs
	var msgs []vm.MsgContract
	amount := sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0))
	for _, addr := range addrs {
		msg := vm.NewMsgContract(addr, nil, code, amount)
		msgs = append(msgs, msg)
	}

	// create contract
	var r *sdk.Result
	for i := 0; i < b.N; i++ {
		r, err = vm.HandleMsgContract(ctx, msgs[i], k)
		require.Nil(b, err)

		for _, event := range r.Events {
			if event.Type == vm.EventTypeContractCreated {
				for _, attr := range event.Attributes {
					if string(attr.GetKey()) == vm.AttributeKeyAddress {
						contractAddr, err := sdk.AccAddressFromBech32(string(attr.GetValue()))
						require.Nil(b, err)
						contractAddrs = append(contractAddrs, contractAddr)
					}
				}
			}
		}
	}

	return
}

func mockSomeString(baseStr string, length int) string {
	str := baseStr
	for len(str) < length {
		str += baseStr
	}
	return str[:length]
}

func generateContractCallMsgs(b *testing.B, addrs []sdk.AccAddress, contractAddrs []sdk.AccAddress, method string, dataSize, bond int) (msgs []vm.MsgContract) {
	amount := sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(int64(bond)))

	for i := 0; i < b.N; i++ {
		data := mockSomeString(contractAddrs[i].String(), dataSize)
		argList := []string{contractAddrs[i].String(), data}
		payload, _, err := vm.GenPayload(contractAbiPath, method, argList)
		require.Nil(b, err)

		msg := vm.NewMsgContract(addrs[i], contractAddrs[i], payload, amount)
		msgs = append(msgs, msg)
	}

	return
}

func generateContractCallMsgsRepeatHash(b *testing.B, addrs []sdk.AccAddress, contractAddrs []sdk.AccAddress, method string, dataSize, repeatHash int) (msgs []vm.MsgContract) {
	amount := sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0))

	for i := 0; i < b.N; i++ {
		data := mockSomeString(contractAddrs[i].String(), dataSize)
		argList := []string{data, fmt.Sprintf("%d", repeatHash)}
		payload, _, err := vm.GenPayload(contractAbiPath, method, argList)
		require.Nil(b, err)

		msg := vm.NewMsgContract(addrs[i], contractAddrs[i], payload, amount)
		msgs = append(msgs, msg)
	}

	return
}

func Benchmark_handleMsgContract_Create(b *testing.B) {
	//prepare
	k, addrs, ms, ctx := prepareTest(b, dbPath)
	defer os.RemoveAll(dbPath)

	// load contract code
	code, err := vm.CodeFromFile(contractCodePath)
	require.Nil(b, err)

	// genenrate contract create msgs
	var msgs []vm.MsgContract
	amount := sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(0))
	for _, addr := range addrs {
		msg := vm.NewMsgContract(addr, nil, code, amount)
		msgs = append(msgs, msg)
	}

	// commit store
	ms.Commit()

	// reset timer
	b.ResetTimer()

	// benchmark test
	for i := 0; i < b.N; i++ {
		vm.HandleMsgContract(ctx, msgs[i], k)
	}

	ms.Commit()
}

func Benchmark_handleMsgContract_Call_VMUpdate(b *testing.B) {
	//prepare
	k, addrs, ms, ctx := prepareTest(b, dbPath)
	defer os.RemoveAll(dbPath)

	contractAddrs := createContracts(b, k, addrs, ctx)
	msgs := generateContractCallMsgs(b, addrs, contractAddrs, contractCallMethod, 1024, 1000000)

	for i := 0; i < b.N; i++ {
		vm.HandleMsgContract(ctx, msgs[i], k)
	}

	msgs = generateContractCallMsgs(b, addrs, contractAddrs, vmUpdateMethod, 1024*10, 1200000)

	ms.Commit()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		vm.HandleMsgContract(ctx, msgs[i], k)
	}

	ms.Commit()
}

func benchmarkContractCallFuncFactory(method string, memSize, bond int) func(b *testing.B) {
	return func(b *testing.B) {
		os.RemoveAll(dbPath)
		k, addrs, ms, ctx := prepareTest(b, dbPath)
		defer os.RemoveAll(dbPath)

		contractAddrs := createContracts(b, k, addrs, ctx)
		msgs := generateContractCallMsgs(b, addrs, contractAddrs, method, memSize, bond)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			vm.HandleMsgContract(ctx, msgs[i], k)
		}

		ms.Commit()
	}
}

func benchmarkContractCallFuncFactoryRepeatHash(method string, dataSize, repeat int) func(b *testing.B) {
	return func(b *testing.B) {
		os.RemoveAll(dbPath)
		k, addrs, ms, ctx := prepareTest(b, dbPath)
		defer os.RemoveAll(dbPath)

		contractAddrs := createContracts(b, k, addrs, ctx)
		msgs := generateContractCallMsgsRepeatHash(b, addrs, contractAddrs, method, dataSize, repeat)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			vm.HandleMsgContract(ctx, msgs[i], k)
		}

		ms.Commit()
	}
}

func Benchmark_handleMsgContract_Call_Memory(b *testing.B) {
	tests := []struct {
		dataSize int
		bond     int
	}{
		{64, 1000000},
		{256, 1000000},
		{512, 1000000},
		{1024, 1000000},
		{1024 * 2, 1000000},
		{1024 * 4, 1000000},
		{1024 * 6, 1000000},
		{1024 * 8, 1000000},
		{1024 * 10, 1000000},
		{1024 * 512, 1000000},
		{1024 * 1024 * 1, 1000000},
	}

	for _, test := range tests {
		b.Run(fmt.Sprintf("memSize_%s", translateSize(test.dataSize)), benchmarkContractCallFuncFactory(contractCallMethod, test.dataSize, test.bond))
		time.Sleep(time.Second * 3)
	}
}

func Benchmark_handleMsgContract_Call_Hashing(b *testing.B) {
	tests := []struct {
		dataSize int
		repeat   int
	}{
		{64, 1},
		{64, 10},
		{64, 20},
		{64, 30},
		{64, 40},
		{64, 100},
		{64, 10},
		{128, 10},
		{512, 10},
		{1024, 10},
		{1024 * 1024, 10},
	}

	for _, test := range tests {
		b.Run(fmt.Sprintf("dataSize:%s_repeatHash:%d", translateSize(test.dataSize), test.repeat), benchmarkContractCallFuncFactoryRepeatHash(hashTest, test.dataSize, test.repeat))
		time.Sleep(time.Second * 3)
	}
}
