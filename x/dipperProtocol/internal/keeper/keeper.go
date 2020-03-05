package keeper

import (
	"github.com/Dipper-Protocol/x/dipperProtocol/internal/types"
	"github.com/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Protocol/types"
	"github.com/Dipper-Protocol/x/bank"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	CoinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the dipperProtocol Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

//Dipper Bank
func (k Keeper) GetBillBank(ctx sdk.Context) types.BillBank {
	store := ctx.KVStore(k.storeKey)
	if !k.IsObjectPresent(ctx, types.DipperBank){
		return types.NewBillBank()
	}
	bz := store.Get([]byte(types.DipperBank))
	var billBank = types.NewBillBank()
	k.cdc.MustUnmarshalBinaryBare(bz, &billBank)
	return billBank
}

func (k Keeper) SetBillBank(ctx sdk.Context, bb types.BillBank) {
	//if len(oracle.TokensPrice) == 0 {
	//	return
	//}
	store := ctx.KVStore(k.storeKey)
	//oracle, _ := bb.GetOracle()
	//fmt.Println("in the end the oracle is", oracle) //for test set price
	store.Set([]byte(types.DipperBank), k.cdc.MustMarshalBinaryBare(bb))
}

//NetValueOf
func (k Keeper) GetNetValueOf(ctx sdk.Context, user sdk.AccAddress) int64 {
	return k.GetBillBank(ctx).NetValueOf(user)
}

//Borrow methods
func (k Keeper)GetBorrowBalanceOf(ctx sdk.Context, symbol string, user sdk.AccAddress) int64 {
	return k.GetBillBank(ctx).BorrowBalanceOf(symbol, user)
}

func (k Keeper)GetBorrowValueOf(ctx sdk.Context, symbol string, user sdk.AccAddress) int64 {
	return k.GetBillBank(ctx).BorrowValueOf(symbol, user)
}

func (k Keeper)GetBorrowValueEstimate(ctx sdk.Context, amount int64, symbol string) int64{
	bank := k.GetBillBank(ctx)
	return bank.BorrowValueEstimate(amount, symbol)
}

func (k Keeper)BankBorrow(ctx sdk.Context, amount sdk.Coins, symbol string, user sdk.AccAddress) error{
	bank := k.GetBillBank(ctx)
	err := bank.Borrow(amount, symbol, user)
	if err != nil {
		return err
	}
	k.SetBillBank(ctx, bank)
	return nil
}

func (k Keeper)BankRepay(ctx sdk.Context, amount sdk.Coins, symbol string, user sdk.AccAddress) error{
	bank := k.GetBillBank(ctx)
	err := bank.Repay(amount, symbol, user)
	if err != nil {
		return err
	}
	k.SetBillBank(ctx, bank)
	return nil
}


//Supply methods
func (k Keeper)GetSupplyBalanceOf(ctx sdk.Context, symbol string, user sdk.AccAddress) int64 {
	bank := k.GetBillBank(ctx)
	return bank.SupplyBalanceOf(symbol, user)
}

func (k Keeper)GetSupplyValueOf(ctx sdk.Context, symbol string, user sdk.AccAddress) int64 {
	bank := k.GetBillBank(ctx)
	return bank.SupplyValueOf(symbol, user)
}

func (k Keeper)BankDeposit(ctx sdk.Context, amount sdk.Coins, symbol string, user sdk.AccAddress) error {
	bank := k.GetBillBank(ctx)
	err := bank.Deposit(amount, symbol, user)
	if err != nil {
		return err
	}
	k.SetBillBank(ctx, bank)
	return nil
}

func (k Keeper)BankWithdraw(ctx sdk.Context, amount sdk.Coins, symbol string, user sdk.AccAddress) error{
	bank := k.GetBillBank(ctx)
	err := bank.Withdraw(amount, symbol, user)
	if err != nil {
		return err
	}
	k.SetBillBank(ctx, bank)
	return nil
}

//Orcale methods
// Gets the entire Whois metadata struct for a name
func (k Keeper) GetBankOracle(ctx sdk.Context) types.Oracle {
	oracle, err := k.GetBillBank(ctx).GetOracle()
	if err != nil {
		return types.NewOracle()
	}
	return oracle
}

func (k Keeper) SetBankOracle(ctx sdk.Context, oracle types.Oracle) {
	bank := k.GetBillBank(ctx)
	bank.Oracle = oracle.Bytes()
	//realbank, _ := bank.GetOracle()
	//fmt.Println("original oracle is", oracle,  "original oracle byte is", oracle.Bytes() ,"real store is", realbank) //for test set price
	k.SetBillBank(ctx, bank)
}

func (k Keeper)GetOraclePrice(ctx sdk.Context,symbol string) int64 {
	oracle := k.GetBankOracle(ctx)
	return oracle.GetPrice(symbol)
}

func (k Keeper)SetOraclePrice(ctx sdk.Context, symbol string, price int64) {
	oracle := k.GetBankOracle(ctx)
	oracle.SetPrice(symbol, price)
	//fmt.Println("oracle's price", oracle) //for test set price
	k.SetBankOracle(ctx, oracle)
}

func (k Keeper) IsObjectPresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}