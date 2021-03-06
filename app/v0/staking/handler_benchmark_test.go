package staking_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/Dipper-Labs/Dipper-Protocol/app"
	"github.com/Dipper-Labs/Dipper-Protocol/app/protocol"
	v0 "github.com/Dipper-Labs/Dipper-Protocol/app/v0"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/bank"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/distribution"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/params"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/slashing"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/staking"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply"
	"github.com/Dipper-Labs/Dipper-Protocol/store"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

const (
	dbPath    = "/tmp/testdata"
	varAddr   = "dipvaloper18dyhj6ncf5r9m5ecseeyv8xjmeyu3ug0pvrjjn"
	varPubkey = "dipvalconspub1zcjduepqe88trqmzgwa044v0k0emax74nxewtcea87muae2g9qpw5smhlc3qqe8ck7"
)

var (
	delegateAmount = sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(sdk.NativeTokenFraction*100))
)

func prepareTest(b *testing.B, dbPath string) (k staking.Keeper, addrs []sdk.AccAddress, ms store.CommitMultiStore, ctx sdk.Context, valAddr sdk.ValAddress) {
	cdc := app.MakeLatestCodec()

	// params
	paramsKeeper := params.NewKeeper(cdc, protocol.Keys[params.StoreKey], protocol.TKeys[params.TStoreKey])
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := paramsKeeper.Subspace(staking.DefaultParamspace)
	slashingSubspace := paramsKeeper.Subspace(slashing.DefaultParamspace)
	distributionSubspace := paramsKeeper.Subspace(distribution.DefaultParamspace)

	// module account perms for supply keeper
	var maccPerms = map[string][]string{
		distribution.ModuleName:   nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	}

	// keeper
	accountKeeper := auth.NewAccountKeeper(cdc, protocol.Keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, bankSubspace, nil)
	supplyKeeper := supply.NewKeeper(cdc, protocol.Keys[supply.StoreKey], accountKeeper, bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		cdc,
		protocol.Keys[staking.StoreKey],
		protocol.TKeys[staking.TStoreKey],
		supplyKeeper,
		stakingSubspace)
	slashingKeeper := slashing.NewKeeper(cdc, protocol.Keys[slashing.StoreKey], stakingKeeper, slashingSubspace)
	distributionKeeper := distribution.NewKeeper(
		cdc,
		protocol.Keys[distribution.StoreKey],
		distributionSubspace,
		stakingKeeper,
		supplyKeeper,
		auth.FeeCollectorName,
		v0.ModuleAccountAddrs())

	// setup stakingKeeper hooks
	stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(distributionKeeper.Hooks(), slashingKeeper.Hooks()),
	)

	// new db
	os.RemoveAll(dbPath)
	db, err := sdk.NewLevelDB("application", dbPath)
	require.Nil(b, err)

	// new store
	ms = store.NewCommitMultiStore(db)
	ms.SetPruning(store.PruneSyncable)

	// mount stores
	ms.MountStoreWithDB(protocol.Keys[params.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.Keys[auth.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.Keys[supply.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.Keys[staking.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.Keys[slashing.StoreKey], sdk.StoreTypeIAVL, nil)
	ms.MountStoreWithDB(protocol.Keys[distribution.StoreKey], sdk.StoreTypeIAVL, nil)

	// mount tstores
	ms.MountStoreWithDB(protocol.TKeys[params.TStoreKey], sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(protocol.TKeys[staking.TStoreKey], sdk.StoreTypeTransient, nil)

	// load store
	ms.LoadLatestVersion()

	// new context
	logger := log.NewNopLogger()
	ctx = sdk.NewContext(ms, abci.Header{}, false, logger)

	// account init balance
	coins := sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(sdk.NativeTokenFraction*1000)))

	// generate accounts
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

	// setup validators
	valAddr, err = sdk.ValAddressFromBech32(varAddr)
	require.Nil(b, err)
	valPubkey, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, varPubkey)
	require.Nil(b, err)
	valDesc := staking.NewDescription("mock_moniker", "mock_identity", "mock_website", "mock_details")
	validator := staking.NewValidator(valAddr, valPubkey, valDesc)
	stakingKeeper.SetValidator(ctx, validator)
	stakingKeeper.AddValidatorTokensAndShares(ctx, validator, sdk.NewInt(sdk.NativeTokenFraction*1000000), true)

	// setup staking params
	params := staking.DefaultParams()
	params.MaxEntries = uint16(b.N + 100)
	params.MaxLever = sdk.NewDec(10000000)
	stakingKeeper.SetParams(ctx, params)

	// setup distribution current rewards
	rewardsAmount := sdk.NewDecCoins(sdk.NewCoins(sdk.NewCoin(sdk.NativeTokenName, sdk.NewInt(sdk.NativeTokenFraction*10))))
	currentRewardsP0 := distribution.NewValidatorCurrentRewards(rewardsAmount, 0)
	currentRewardsP1 := distribution.NewValidatorCurrentRewards(rewardsAmount, 1)
	distributionKeeper.SetValidatorCurrentRewards(ctx, valAddr, currentRewardsP0)
	distributionKeeper.SetValidatorCurrentRewards(ctx, valAddr, currentRewardsP1)

	// setup distribution historical rewards
	historicalRewards := distribution.NewValidatorHistoricalRewards(rewardsAmount, 1)
	distributionKeeper.SetValidatorHistoricalRewards(ctx, valAddr, 0, historicalRewards)
	distributionKeeper.SetValidatorHistoricalRewards(ctx, valAddr, 1, historicalRewards)

	//setup distribution outstanding rewards
	distributionKeeper.SetValidatorOutstandingRewards(ctx, valAddr, sdk.DecCoins{})

	// setup distribution FeePool
	distributionKeeper.SetFeePool(ctx, distribution.FeePool{})

	return stakingKeeper, addrs, ms, ctx, valAddr
}

func Benchmark_handleMsgDelegate(b *testing.B) {
	k, addrs, ms, ctx, valAddr := prepareTest(b, dbPath)

	// setup delegate msgs
	var msgs []staking.MsgDelegate
	for _, addr := range addrs {
		msg := staking.NewMsgDelegate(addr, valAddr, delegateAmount)
		msgs = append(msgs, msg)
	}

	// reset timer
	b.ResetTimer()

	// benchmark test
	for i := 0; i < b.N; i++ {
		staking.HandleMsgDelegate(ctx, msgs[i], k)
	}

	ms.Commit()
}

func Benchmark_handleMsgUnDelegate(b *testing.B) {
	k, addrs, ms, ctx, valAddr := prepareTest(b, dbPath)

	// do delegate
	var delegateMsgs []staking.MsgDelegate
	for _, addr := range addrs {
		msg := staking.NewMsgDelegate(addr, valAddr, delegateAmount)
		delegateMsgs = append(delegateMsgs, msg)
	}

	for i := 0; i < b.N; i++ {
		staking.HandleMsgDelegate(ctx, delegateMsgs[i], k)
	}

	// do undelegate
	var undelegateMsgs []staking.MsgUndelegate
	for _, addr := range addrs {
		msg := staking.NewMsgUndelegate(addr, valAddr, delegateAmount)
		undelegateMsgs = append(undelegateMsgs, msg)
	}

	ms.Commit()
	// reset timer
	b.ResetTimer()

	// benchmark test
	for i := 0; i < b.N; i++ {
		staking.HandleMsgUndelegate(ctx, undelegateMsgs[i], k)
	}

	ms.Commit()
}
