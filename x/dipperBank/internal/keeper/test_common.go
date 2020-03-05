package keeper

import (
	"github.com/Dipper-Protocol/x/dipperBank/internal/types"
	"github.com/Dipper-Protocol/codec"
	"github.com/Dipper-Protocol/store"
	sdk "github.com/Dipper-Protocol/types"
	"github.com/Dipper-Protocol/x/auth"
	"github.com/Dipper-Protocol/x/bank"
	"github.com/Dipper-Protocol/x/params"
	"github.com/Dipper-Protocol/x/staking"
	"github.com/Dipper-Protocol/x/supply"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tm-db"
)


//nolint: deadcode unused
var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delPk2   = ed25519.GenPrivKey().PubKey()
	delPk3   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
	delAddr2 = sdk.AccAddress(delPk2.Address())
	delAddr3 = sdk.AccAddress(delPk3.Address())

	ValOpPk1    = ed25519.GenPrivKey().PubKey()
	ValOpPk2    = ed25519.GenPrivKey().PubKey()
	ValOpPk3    = ed25519.GenPrivKey().PubKey()
	ValOpAddr1  = sdk.ValAddress(ValOpPk1.Address())
	ValOpAddr2  = sdk.ValAddress(ValOpPk2.Address())
	ValOpAddr3  = sdk.ValAddress(ValOpPk3.Address())
	valAccAddr1 = sdk.AccAddress(ValOpPk1.Address()) // generate acc addresses for these validator keys too
	valAccAddr2 = sdk.AccAddress(ValOpPk2.Address())
	valAccAddr3 = sdk.AccAddress(ValOpPk3.Address())

	ValConsPk11  = ed25519.GenPrivKey().PubKey()
	ValConsPk12  = ed25519.GenPrivKey().PubKey()
	ValConsPk13  = ed25519.GenPrivKey().PubKey()
	ValConsAddr1 = sdk.ConsAddress(ValConsPk11.Address())
	ValConsAddr2 = sdk.ConsAddress(ValConsPk12.Address())
	ValConsAddr3 = sdk.ConsAddress(ValConsPk13.Address())

	// TODO move to common testing package for all modules
	// test addresses
	TestAddrs = []sdk.AccAddress{
		delAddr1, delAddr2, delAddr3,
		valAccAddr1, valAccAddr2, valAccAddr3,
	}

	emptyDelAddr sdk.AccAddress
	emptyValAddr sdk.ValAddress
	emptyPubkey  crypto.PubKey
	stakeDenom   = "stake"
	feeDenom     = "fee"
)
// test common should produce a staking keeper, a supply keeper, a bank keeper, an auth keeper, a validatorvesting keeper, a context,

func CreateTestInput(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, auth.AccountKeeper, bank.Keeper, staking.Keeper, supply.Keeper, Keeper) {

	initTokens := sdk.TokensFromConsensusPower(initPower)

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	dipperBank := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(dipperBank, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	require.Nil(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foo-chain"}, isCheckTx, log.NewNopLogger())

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
	validatorVestingAcc := supply.NewEmptyModuleAccount(types.ModuleName)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true
	blacklistedAddrs[validatorVestingAcc.GetAddress().String()] = true

	cdc := MakeTestCodec()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)

	stakingParams := staking.NewParams(time.Hour, 100, uint16(7), sdk.DefaultBondDenom)

	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		types.ModuleName:          {supply.Burner},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)

	stakingKeeper := staking.NewKeeper(cdc, keyStaking, tkeyParams,supplyKeeper, pk.Subspace(staking.DefaultParamspace), staking.DefaultCodespace)
	stakingKeeper.SetParams(ctx, stakingParams)

	keeper := NewKeeper(bankKeeper, dipperBank, cdc)

	initCoins := sdk.NewCoins(sdk.NewCoin(stakingKeeper.BondDenom(ctx), initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(stakingKeeper.BondDenom(ctx), initTokens.MulRaw(int64(len(TestAddrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range TestAddrs {
		_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		require.Nil(t, err)
	}

	// set module accounts
	//keeper.supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	//keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	//keeper.supplyKeeper.SetModuleAccount(ctx, bondPool)

	return ctx, accountKeeper, bankKeeper, stakingKeeper, supplyKeeper, keeper
}


func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	//vesting.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}
