package types

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MinNamePrice is Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("dpc", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// NewWhois returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}

type TokenPool struct {
	SupplyBill int64 `json:"supply_bill"`
	Supply     int64 `json:"supply"`
	BorrowBill int64 `json:"borrow_bill"`
	Borrow     int64 `json:"borrow"`

	// last liquidate blockNumber
	liquidateIndex uint64
}

// implement fmt.Stringer
func (tp TokenPool) String() string {
	return strings.TrimSpace(fmt.Sprintf(`SupplyBill: %d
Supply: %d
BorrowBill: %d
Borrow: %d`, tp.SupplyBill, tp.Supply, tp.BorrowBill, tp.Borrow))
}

// GetCash Cash = Supply - Borrow
func (tp TokenPool) GetCash() int64 {
	return tp.Supply - tp.Borrow
}

type BillBank struct {
	// internal account for token bill(deposit)
	AccountDepositBills []byte `json:"account_deposit_bills"`
	// internal account for token bill(borrow)
	AccountBorrowBills []byte `json:"account_deposit_bills"`

	Pools []byte `json:"account_deposit_bills"`

	Oracle []byte `json:"oracle"`

	// BlockNumber simulate
	BlockNumber uint64 `json:"block_number"`
	// borrowRate every block
	borrowRate uint64 `json:"borrow_rate"`
}



func NewBillBank() BillBank {
	return BillBank{
		AccountDepositBills: NewAccountBills().Bytes(),
		AccountBorrowBills:  NewAccountBills().Bytes(),
		Pools: 		NewPools().Bytes(),
		Oracle:     NewOracle().Bytes(),
		BlockNumber: 1,
		borrowRate:  1,
	}
}

//BillBank
func (b *BillBank) liquidate(symbol string) {
	pool := b.getPool(symbol)

	growth := b.calculateGrowth(symbol)
	pool.Supply += growth
	pool.Borrow += growth

	// update pool
	pool.liquidateIndex = b.BlockNumber
	pools, err := b.GetPools()
	if err != nil {
		return
	}
	pools.SubPools[symbol] = pool
}

func (b *BillBank) calculateGrowth(symbol string) int64 {
	pool := b.getPool(symbol)

	var growth int64 = 0
	borrow := pool.Borrow
	if borrow != 0 {
		// Compound interest
		// formula:
		//		b: borrow
		//		r: rate
		//		n: block number
		//		b = b * (1+r)^n
		borrow = int64(float64(borrow) * math.Pow(
			1+float64(b.borrowRate)/100,
			float64(b.BlockNumber-pool.liquidateIndex),
		))
		growth = borrow - pool.Borrow
	}
	return growth
}

func (b BillBank) getPool(symbol string) (pool TokenPool) {
	var ok bool
	pools, err := b.GetPools()
	if err != nil {
		return TokenPool{}
	}
	if pool, ok = pools.SubPools[symbol]; !ok {
		log.Panicf("not support token: %v", symbol)
	}
	return
}

func (b BillBank) NetValueOf(userAcc sdk.AccAddress) int64 {
	user := userAcc.String()
	var supplyValue int64 = 0
	depositBills, err := b.GetAccountDepositBills()
	if err != nil {
		return 0
	}
	if acc, ok := depositBills.Bills[user]; ok {
		for sym, bill := range acc {
			if bill != 0 {
				supplyValue += b.SupplyValueOf(sym, userAcc)
			}
		}
	}

	var borrowValue int64 = 0
	borrowBills, err := b.GetAccountDepositBills()
	if err != nil {
		return 0
	}
	if acc, ok := borrowBills.Bills[user]; ok {
		for sym, bill := range acc {
			if bill != 0 {
				borrowValue += b.BorrowValueOf(sym, userAcc)
			}
		}
	}

	return supplyValue - borrowValue
}



//BillBank Borrow method
func (b BillBank) BorrowBalanceOf(symbol string, userAcc sdk.AccAddress) int64 {
	user := userAcc.String()
	pool := b.getPool(symbol)


	borrowBills, err := b.GetAccountDepositBills()
	if err != nil {
		return 0
	}
	// check bill
	if _, ok := borrowBills.Bills[user]; !ok {
		return 0
	}
	var bill int64 = 0
	if b, ok := borrowBills.Bills[user][symbol]; ok {
		bill = b
	}
	if bill == 0 {
		return 0
	}

	// calcuate amount
	// current block liquidated, growth is zero
	growth := b.calculateGrowth(symbol)
	return int64(float64(bill) * float64(pool.Borrow + growth) / float64(pool.BorrowBill))
}

func (b BillBank) BorrowValueOf(symbol string, userAcc sdk.AccAddress) int64 {
	return b.BorrowValueEstimate(
		b.BorrowBalanceOf(symbol, userAcc),
		symbol,
	)
}

func (b *BillBank) BorrowValueEstimate(amount int64, symbol string) int64 {
	oracle, err := b.GetOracle()
	if err != nil {
		return 0
	}
	return amount * oracle.GetPrice(symbol)/1000000
}

func (b *BillBank) Borrow(amount sdk.Coins, symbol string, userAcc sdk.AccAddress) error {
	user := userAcc.String()
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	coin := amount.AmountOf(symbol).Int64()
	// check cash of pool
	if coin > pool.GetCash() {
		return fmt.Errorf("not enough token for borrow. amount: %v, cash: %v", amount, pool.GetCash())
	}

	// calcuate bill
	bill := coin
	if pool.BorrowBill != 0 && pool.Borrow != 0 {
		bill = coin * (pool.BorrowBill / pool.Borrow)
	}


	borrowBills, err := b.GetAccountBorrowBills()
	if err != nil {
		return fmt.Errorf("there is no deposit bills")
	}
	// update user account bill
	if accountBorrow, ok := borrowBills.Bills[user]; ok {
		if _, ok := accountBorrow[symbol]; ok {
			borrowBills.Bills[user][symbol] += bill
		} else {
			borrowBills.Bills[user][symbol] = bill
		}
	} else {
		borrowBills.Bills[user] = map[DSymbol]DBill{symbol: bill}
	}

	// update borrow
	pool.BorrowBill += bill
	pool.Borrow += coin

	pools, err := b.GetPools()
	if err != nil {
		return fmt.Errorf("there is no pools")
	}
	pools.SubPools[symbol] = pool
	//TODO set change to billbank
	b.AccountBorrowBills = borrowBills.Bytes()
	b.Pools = pools.Bytes()

	return nil
}

func (b *BillBank) Repay(amount sdk.Coins, symbol string, userAcc sdk.AccAddress) error {
	user := userAcc.String()
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check borrow
	accountAmount := b.BorrowBalanceOf(symbol, userAcc)
	coin := amount.AmountOf(symbol).Int64()
	if coin > accountAmount {
		return fmt.Errorf("too much amount to repay. user: %v, need repay: %v", user, accountAmount)
	}

	// calculate bill
	bill := int64(float64(coin) * (float64(pool.BorrowBill) / float64(pool.Borrow)))

	borrowBills, err := b.GetAccountBorrowBills()
	if err != nil {
		return fmt.Errorf("there is no borrow bills")
	}

	// update user account borrow
	borrowBills.Bills[user][symbol] -= bill

	// update borrow
	pool.BorrowBill -= bill
	pool.Borrow -= coin

	pools, err := b.GetPools()
	if err != nil {
		return fmt.Errorf("there is no pools")
	}

	pools.SubPools[symbol] = pool


	b.AccountBorrowBills = borrowBills.Bytes()
	b.Pools = pools.Bytes()
	return nil
}

//BillBank supply methods
func (b *BillBank) SupplyBalanceOf(symbol string, userAcc sdk.AccAddress) int64 {
	user := userAcc.String()
	pool := b.getPool(symbol)

	depositBills, err := b.GetAccountDepositBills()
	if err != nil {
		return 0
	}
	// check bill
	if _, ok := depositBills.Bills[user]; !ok {
		return 0
	}
	var bill int64 = 0
	if b, ok := depositBills.Bills[user][symbol]; ok {
		bill = b
	}
	if bill == 0 {
		return 0
	}

	// calcuate amount
	// current block liquidated, growth is zero
	growth := b.calculateGrowth(symbol)
	return int64(float64(bill) * (float64(pool.Supply + growth) / float64(pool.SupplyBill)))
}

func (b *BillBank) SupplyValueOf(symbol string, userAcc sdk.AccAddress) int64 {
	oracle, err := b.GetOracle()
	if err != nil {
		return 0
	}
	return b.SupplyBalanceOf(symbol, userAcc) * oracle.GetPrice(symbol)/1000000
}

func (b *BillBank) Deposit(amount sdk.Coins, symbol string, userAcc sdk.AccAddress) error {
	user := userAcc.String()
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// calcuate bill
	coin := amount.AmountOf(symbol).Int64()
	bill := coin
	if pool.SupplyBill != 0 && pool.Supply != 0 {
		bill = int64(float64(coin) * float64(pool.SupplyBill / pool.Supply))
	}

	depostiBills, err := b.GetAccountDepositBills()
	if err != nil {
		return fmt.Errorf("there is no bills")
	}

	// update user account bill
	if accountBill, ok := depostiBills.Bills[user]; ok {
		if _, ok := accountBill[symbol]; ok {
			depostiBills.Bills[user][symbol] += bill
		} else {
			depostiBills.Bills[user][symbol] = bill
		}
	} else {
		depostiBills.Bills[user] = map[DSymbol]DBill{symbol: bill}
	}

	pools, err := b.GetPools()
	if err != nil {
		return fmt.Errorf("there is no pools")
	}
	// update pool
	pool.SupplyBill += bill
	pool.Supply += coin
	pools.SubPools[symbol] = pool


	b.AccountDepositBills = depostiBills.Bytes()
	b.Pools = pools.Bytes()

	return nil
}

func (b *BillBank) Withdraw(amount sdk.Coins, symbol string, userAcc sdk.AccAddress) (err error) {
	user := userAcc.String()
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check account balance
	accountAmount := b.SupplyBalanceOf(symbol, userAcc)
	coin := amount.AmountOf(symbol).Int64()
	if coin > accountAmount {
		return fmt.Errorf("not enough amount for withdraw. user: %v, acutal amount: %v", userAcc.String(), accountAmount)
	}
	// check balance of supply
	if coin > pool.GetCash() {
		return fmt.Errorf("not enough token for withdraw. amount: %v, cash %v", amount, pool.GetCash())
	}

	// calcuate bill
	bill := int64(float64(coin) * float64(pool.SupplyBill / pool.Supply))

	// update user account bill
	depositBills, err := b.GetAccountDepositBills()
	if err != nil {
		return fmt.Errorf("there is no deposit bills")
	}

	depositBills.Bills[user][symbol] -= bill

	pools, err := b.GetPools()
	if err != nil {
		return fmt.Errorf("there is no pools")
	}
	// update pool
	pool.SupplyBill -= bill
	pool.Supply -= coin
	pools.SubPools[symbol] = pool

	b.AccountBorrowBills = depositBills.Bytes()
	b.Pools = pools.Bytes()

	return
}

func (b BillBank)GetAccountDepositBills() (AccountBills, error) {
	o := AccountBills{}
	err := json.Unmarshal(b.AccountDepositBills, &o)
	if err != nil {
		return NewAccountBills(), err
	}
	return o, nil
}

func (b BillBank)GetAccountBorrowBills() (AccountBills, error) {
	o := AccountBills{}
	err := json.Unmarshal(b.AccountBorrowBills, &o)
	if err != nil {
		return NewAccountBills(), err
	}
	return o, nil
}

func (b BillBank)GetPools() (Pools, error) {
	o := Pools{}
	err := json.Unmarshal(b.Pools, &o)
	if err != nil {
		return NewPools(), err
	}
	return o, nil
}

func (b BillBank)GetOracle() (Oracle, error) {
	o := Oracle{}
	err := json.Unmarshal(b.Oracle, &o)
	if err != nil {
		//fmt.Println("shit fuck", err.Error()) //for test set price
		return NewOracle(), err
	}
	return o, nil
}

type AccountBills struct {
	Bills map[DUser]map[DSymbol]DBill `json:"account_bills"`
}

func NewAccountBills() AccountBills {
	return AccountBills{map[DUser]map[DSymbol]DBill{"init":{"init2":999}}}
}

func (ab AccountBills) Bytes() []byte{
	initBytes, err := json.Marshal(ab)
	if err != nil {
		return []byte{}
	}
	return initBytes
}

type Pools struct {
	SubPools map[DSymbol]TokenPool `json:"sub_pools"`
}

func NewPools() Pools {
	return Pools{map[DSymbol]TokenPool{
		"ETH": TokenPool{},
		"DAI": TokenPool{},}}
}

func (p Pools) Bytes() []byte{
	initBytes, err := json.Marshal(p)
	if err != nil {
		return []byte{}
	}
	return initBytes
}

//Oracle, maybe use chainlink later
type Oracle struct {
	TokensPrice map[DSymbol]DPrice `json:"tokensPrice"`
}

func NewOracle() Oracle {
	return Oracle{map[DSymbol]DPrice{"init":999}}
}

// implement fmt.Stringer
func (o Oracle) String() string {
	mjson,_ :=json.Marshal(o.TokensPrice)
	mString :=string(mjson)
	return mString
}

func (o Oracle) Bytes() []byte{
	initBytes, err := json.Marshal(o)
	if err != nil {
		return []byte{}
	}
	return initBytes
}


func (o Oracle) GetPrice(symbol string) int64 {
	if v, ok := o.TokensPrice[symbol]; ok {
		return v
	}
	return 0
}

func (o *Oracle) SetPrice(symbol string, price int64) {
	o.TokensPrice[symbol] = price
}