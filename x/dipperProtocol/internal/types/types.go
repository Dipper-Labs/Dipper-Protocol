package types

import (
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
	SupplyBill int64 `json:"supplyBill"`
	Supply     int64 `json:"supply"`
	BorrowBill int64 `json:"BorrowBill"`
	Borrow     int64 `json:"Borrow"`

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
func (tp *TokenPool) GetCash() int64 {
	return tp.Supply - tp.Borrow
}


type BillBank struct {
	// internal account for token bill(deposit)
	AccountDepositBills map[tUser]map[tSymbol]tBill `json:"accountDepositBills"`
	// internal account for token bill(borrow)
	AccountBorrowBills map[tUser]map[tSymbol]tBill `json:"accountDepositBills"`

	Pools map[tSymbol]TokenPool `json:"AccountDepositBills"`

	Oracler Oracle `json:"AccountDepositBills"`

	// BlockNumber simulate
	BlockNumber uint64 `json:"AccountDepositBills"`
	// borrowRate every block
	borrowRate uint64 `json:"AccountDepositBills"`
}

func NewBillBank() *BillBank {
	return &BillBank{
		AccountDepositBills: map[tUser]map[tSymbol]tBill{},
		AccountBorrowBills:  map[tUser]map[tSymbol]tBill{},
		Pools: map[tSymbol]TokenPool{
			"ETH": TokenPool{},
			"DAI": TokenPool{},
		},
		Oracler:     NewOracle(),
		BlockNumber: 1,
		borrowRate:  100000,
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
	b.Pools[symbol] = pool
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
			float64(1000000+b.borrowRate),
			float64(b.BlockNumber-pool.liquidateIndex),
		))
		growth = borrow - pool.Borrow
	}
	return growth
}

func (b *BillBank) getPool(symbol string) (pool TokenPool) {
	var ok bool
	if pool, ok = b.Pools[symbol]; !ok {
		log.Panicf("not support token: %v", symbol)
	}
	return
}

func (b *BillBank) NetValueOf(userAcc sdk.AccAddress) int64 {
	user := userAcc.String()
	var supplyValue int64 = 0
	if acc, ok := b.AccountDepositBills[user]; ok {
		for sym, bill := range acc {
			if bill != 0 {
				supplyValue += b.SupplyValueOf(sym, userAcc)
			}
		}
	}

	var borrowValue int64 = 0
	if acc, ok := b.AccountBorrowBills[user]; ok {
		for sym, bill := range acc {
			if bill != 0 {
				borrowValue += b.BorrowValueOf(sym, userAcc)
			}
		}
	}

	return supplyValue - borrowValue
}



//BillBank Borrow method
func (b *BillBank) BorrowBalanceOf(symbol string, userAcc sdk.AccAddress) int64 {
	user := userAcc.String()
	pool := b.getPool(symbol)

	// check bill
	if _, ok := b.AccountBorrowBills[user]; !ok {
		return 0
	}
	var bill int64 = 0
	if b, ok := b.AccountBorrowBills[user][symbol]; ok {
		bill = b
	}
	if bill == 0 {
		return 0
	}

	// calcuate amount
	// current block liquidated, growth is zero
	growth := b.calculateGrowth(symbol)
	return bill * (pool.Borrow + growth) / pool.BorrowBill
}

func (b *BillBank) BorrowValueOf(symbol string, userAcc sdk.AccAddress) int64 {
	return b.BorrowValueEstimate(
		b.BorrowBalanceOf(symbol, userAcc),
		symbol,
	)
}

func (b *BillBank) BorrowValueEstimate(amount int64, symbol string) int64 {
	return amount * b.Oracler.GetPrice(symbol)
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

	// update user account bill
	if accountBorrow, ok := b.AccountBorrowBills[user]; ok {
		if _, ok := accountBorrow[symbol]; ok {
			b.AccountBorrowBills[user][symbol] += bill
		} else {
			b.AccountBorrowBills[user][symbol] = bill
		}
	} else {
		b.AccountBorrowBills[user] = map[tSymbol]tBill{symbol: bill}
	}

	// update borrow
	pool.BorrowBill += bill
	pool.Borrow += coin
	b.Pools[symbol] = pool

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
	bill := coin * (pool.BorrowBill / pool.Borrow)

	// update user account borrow
	b.AccountBorrowBills[user][symbol] -= bill

	// update borrow
	pool.BorrowBill -= bill
	pool.Borrow -= coin
	b.Pools[symbol] = pool

	return nil
}

//BillBank supply methods
func (b *BillBank) SupplyBalanceOf(symbol string, userAcc sdk.AccAddress) int64 {
	user := userAcc.String()
	pool := b.getPool(symbol)

	// check bill
	if _, ok := b.AccountDepositBills[user]; !ok {
		return 0
	}
	var bill int64 = 0
	if b, ok := b.AccountDepositBills[user][symbol]; ok {
		bill = b
	}
	if bill == 0 {
		return 0
	}

	// calcuate amount
	// current block liquidated, growth is zero
	growth := b.calculateGrowth(symbol)
	return bill * ((pool.Supply + growth) / pool.SupplyBill)
}

func (b *BillBank) SupplyValueOf(symbol string, userAcc sdk.AccAddress) int64 {
	return b.SupplyBalanceOf(symbol, userAcc) * b.Oracler.GetPrice(symbol)
}

func (b *BillBank) Deposit(amount sdk.Coins, symbol string, userAcc sdk.AccAddress) error {
	user := userAcc.String()
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// calcuate bill
	coin := amount.AmountOf(symbol).Int64()
	bill := coin
	if pool.SupplyBill != 0 && pool.Supply != 0 {
		bill = coin * (pool.SupplyBill / pool.Supply)
	}

	// update user account bill
	if accountBill, ok := b.AccountDepositBills[user]; ok {
		if _, ok := accountBill[symbol]; ok {
			b.AccountDepositBills[user][symbol] += bill
		} else {
			b.AccountDepositBills[user][symbol] = bill
		}
	} else {
		b.AccountDepositBills[user] = map[tSymbol]tBill{symbol: bill}
	}

	// update pool
	pool.SupplyBill += bill
	pool.Supply += coin
	b.Pools[symbol] = pool

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
	bill := coin * (pool.SupplyBill / pool.Supply)

	// update user account bill
	b.AccountDepositBills[user][symbol] -= bill

	// update pool
	pool.SupplyBill -= bill
	pool.Supply -= coin
	b.Pools[symbol] = pool

	return
}


//Oracle, maybe use chainlink later
type Oracle struct {
	TokensPrice map[tSymbol]tPrice `json:"tokensPrice"`
}

func NewOracle() Oracle {
	return Oracle{map[tSymbol]tPrice{}}
}

// implement fmt.Stringer
func (o Oracle) String() string {
	return "DipperBank" //TODO add some orcale info.
}

func (o *Oracle) GetPrice(symbol string) int64 {
	if v, ok := o.TokensPrice[symbol]; ok {
		return v
	}
	return 0
}

func (o *Oracle) SetPrice(symbol string, price int64) {
	o.TokensPrice[symbol] = price
}