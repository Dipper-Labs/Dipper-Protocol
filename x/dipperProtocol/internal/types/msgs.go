package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

//MsgBankRepay define someone repay money to bank
type MsgBankBorrow struct {
	Amount sdk.Coins `json:"amount"`
	Symbol string `json:"symbol"`
	Owner sdk.AccAddress `json:"owner"`
}

func NewMsgBankBorrow(amount sdk.Coins, symbol string, owner sdk.AccAddress) MsgBankBorrow {
	return MsgBankBorrow{
		Amount: amount,
		Symbol: symbol,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgBankBorrow) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBankBorrow) Type() string { return "repay_money" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBankBorrow) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if !msg.Amount.IsAllPositive() || len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Amount and/or Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBankBorrow) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBankBorrow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//MsgBankRepay define someone repay money to bank
type MsgBankRepay struct {
	Amount sdk.Coins `json:"amount"`
	Symbol string `json:"symbol"`
	Owner sdk.AccAddress `json:"owner"`
}

func NewMsgBankRepay(amount sdk.Coins, symbol string, owner sdk.AccAddress) MsgBankRepay {
	return MsgBankRepay{
		Amount: amount,
		Symbol: symbol,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgBankRepay) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBankRepay) Type() string { return "repay_money" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBankRepay) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if !msg.Amount.IsAllPositive() || len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Amount and/or Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBankRepay) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBankRepay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//MsgBankBorrow define someone deposit money to bank
type MsgBankDeposit struct {
	Amount sdk.Coins `json:"amount"`
	Symbol string `json:"symbol"`
	Owner sdk.AccAddress `json:"owner"`
}

func NewMsgBankDeposit(amount sdk.Coins, symbol string, owner sdk.AccAddress) MsgBankDeposit {
	return MsgBankDeposit{
		Amount: amount,
		Symbol: symbol,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgBankDeposit) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBankDeposit) Type() string { return "deposit_money" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBankDeposit) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if !msg.Amount.IsAllPositive() || len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Amount and/or Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBankDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBankDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//MsgBankWithdraw define someone withdraw money from bank
type MsgBankWithdraw struct {
	Amount sdk.Coins `json:"amount"`
	Symbol string `json:"symbol"`
	Owner sdk.AccAddress `json:"owner"`
}

func NewMsgBankWithdraw(amount sdk.Coins, symbol string, owner sdk.AccAddress) MsgBankWithdraw {
	return MsgBankWithdraw{
		Amount: amount,
		Symbol: symbol,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgBankWithdraw) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBankWithdraw) Type() string { return "withdraw_money" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBankWithdraw) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if !msg.Amount.IsAllPositive() || len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("Amount and/or Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBankWithdraw) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBankWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

//MsgBankBorrow define someone who has been chose to set oracle price.
type MsgSetOraclePrice struct {
	//Name string `json:"name"`
	Symbol string `json:"symbol"`
	Price string `json:"price"`
	Owner sdk.AccAddress `json:"owner"`
}

func NewMsgSetOraclePrice(name string, symbol string, amount string, owner sdk.AccAddress) MsgSetOraclePrice {
	return MsgSetOraclePrice{
		//Name: name,
		Symbol: symbol,
		Price: amount,
		Owner:  owner,
	}
}

// Route should return the name of the module
func (msg MsgSetOraclePrice) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetOraclePrice) Type() string { return "set_oracle_price" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetOraclePrice) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Price) == 0 || len(msg.Symbol) == 0{
		return sdk.ErrUnknownRequest("Price and/or Name and/or Symbol cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetOraclePrice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetOraclePrice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
