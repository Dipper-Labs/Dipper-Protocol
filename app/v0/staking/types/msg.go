package types

import (
	"bytes"
	"encoding/json"

	"github.com/tendermint/tendermint/crypto"
	"gopkg.in/yaml.v2"

	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	sdkerrors "github.com/Dipper-Labs/Dipper-Protocol/types/errors"
)

const (
	TypeMsgCreateValidator = "create_validator"
	TypeMsgEditValidator   = "edit_validator"
	TypeMsgDelegate        = "delegate"
	TypeMsgBeginRedelegate = "begin_redelegate"
	TypeMsgUndelegate      = "begin_unbonding"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgCreateValidator{}
	_ sdk.Msg = &MsgEditValidator{}
	_ sdk.Msg = &MsgDelegate{}
	_ sdk.Msg = &MsgUndelegate{}
	_ sdk.Msg = &MsgBeginRedelegate{}
)

//______________________________________________________________________

// MsgCreateValidator - struct for bonding transactions
type MsgCreateValidator struct {
	Description       Description     `json:"description" yaml:"description"`
	Commission        CommissionRates `json:"commission" yaml:"commission"`
	MinSelfDelegation sdk.Int         `json:"min_self_delegation" yaml:"min_self_delegation"`
	DelegatorAddress  sdk.AccAddress  `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress  sdk.ValAddress  `json:"validator_address" yaml:"validator_address"`
	PubKey            crypto.PubKey   `json:"pubkey" yaml:"pubkey"`
	Value             sdk.Coin        `json:"value" yaml:"value"`
}

type msgCreateValidatorJSON struct {
	Description       Description     `json:"description" yaml:"description"`
	Commission        CommissionRates `json:"commission" yaml:"commission"`
	MinSelfDelegation sdk.Int         `json:"min_self_delegation" yaml:"min_self_delegation"`
	DelegatorAddress  sdk.AccAddress  `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress  sdk.ValAddress  `json:"validator_address" yaml:"validator_address"`
	PubKey            string          `json:"pubkey" yaml:"pubkey"`
	Value             sdk.Coin        `json:"value" yaml:"value"`
}

// NewMsgCreateValidator creates a new MsgCreateValidator instance.
// Delegator address and validator address are the same.
func NewMsgCreateValidator(
	valAddr sdk.ValAddress, pubKey crypto.PubKey, selfDelegation sdk.Coin,
	description Description, commission CommissionRates, minSelfDelegation sdk.Int,
) MsgCreateValidator {

	return MsgCreateValidator{
		Description:       description,
		DelegatorAddress:  sdk.AccAddress(valAddr),
		ValidatorAddress:  valAddr,
		PubKey:            pubKey,
		Value:             selfDelegation,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
	}
}

//nolint
func (msg MsgCreateValidator) Route() string { return RouterKey }
func (msg MsgCreateValidator) Type() string  { return TypeMsgCreateValidator }

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgCreateValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	addrs := []sdk.AccAddress{msg.DelegatorAddress}

	if !bytes.Equal(msg.DelegatorAddress.Bytes(), msg.ValidatorAddress.Bytes()) {
		// if validator addr is not same as delegator addr, validator must sign
		// msg as well
		addrs = append(addrs, sdk.AccAddress(msg.ValidatorAddress))
	}
	return addrs
}

// MarshalJSON implements the json.Marshaler interface to provide custom JSON
// serialization of the MsgCreateValidator type.
func (msg MsgCreateValidator) MarshalJSON() ([]byte, error) {
	return json.Marshal(msgCreateValidatorJSON{
		Description:       msg.Description,
		Commission:        msg.Commission,
		DelegatorAddress:  msg.DelegatorAddress,
		ValidatorAddress:  msg.ValidatorAddress,
		PubKey:            sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, msg.PubKey),
		Value:             msg.Value,
		MinSelfDelegation: msg.MinSelfDelegation,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface to provide custom
// JSON deserialization of the MsgCreateValidator type.
func (msg *MsgCreateValidator) UnmarshalJSON(bz []byte) error {
	var msgCreateValJSON msgCreateValidatorJSON
	if err := json.Unmarshal(bz, &msgCreateValJSON); err != nil {
		return err
	}

	msg.Description = msgCreateValJSON.Description
	msg.Commission = msgCreateValJSON.Commission
	msg.DelegatorAddress = msgCreateValJSON.DelegatorAddress
	msg.ValidatorAddress = msgCreateValJSON.ValidatorAddress
	var err error
	msg.PubKey, err = sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, msgCreateValJSON.PubKey)
	if err != nil {
		return err
	}
	msg.Value = msgCreateValJSON.Value
	msg.MinSelfDelegation = msgCreateValJSON.MinSelfDelegation

	return nil
}

// MarshalYAML implements a custom marshal yaml function due to consensus pubkey.
func (msg MsgCreateValidator) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		Description       Description
		Commission        CommissionRates
		MinSelfDelegation sdk.Int
		DelegatorAddress  sdk.AccAddress
		ValidatorAddress  sdk.ValAddress
		PubKey            string
		Value             sdk.Coin
	}{
		Description:       msg.Description,
		Commission:        msg.Commission,
		MinSelfDelegation: msg.MinSelfDelegation,
		DelegatorAddress:  msg.DelegatorAddress,
		ValidatorAddress:  msg.ValidatorAddress,
		PubKey:            sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, msg.PubKey),
		Value:             msg.Value,
	})

	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgCreateValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	if msg.DelegatorAddress.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr
	}
	if !sdk.AccAddress(msg.ValidatorAddress).Equals(msg.DelegatorAddress) {
		return ErrBadValidatorAddr
	}
	if !msg.Value.Amount.IsPositive() {
		return ErrBadDelegationAmount
	}
	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}
	if msg.Commission == (CommissionRates{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}
	if err := msg.Commission.Validate(); err != nil {
		return err
	}
	if !msg.MinSelfDelegation.IsPositive() {
		return ErrMinSelfDelegationInvalid
	}
	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return ErrSelfDelegationBelowMinimum
	}

	return nil
}

// MsgEditValidator - struct for editing a validator
type MsgEditValidator struct {
	Description
	ValidatorAddress sdk.ValAddress `json:"address" yaml:"address"`

	// We pass a reference to the new commission rate and min self delegation as it's not mandatory to
	// update. If not updated, the deserialized rate will be zero with no way to
	// distinguish if an update was intended.
	//
	// REF: #2373
	CommissionRate    *sdk.Dec `json:"commission_rate" yaml:"commission_rate"`
	MinSelfDelegation *sdk.Int `json:"min_self_delegation" yaml:"min_self_delegation"`
}

func NewMsgEditValidator(valAddr sdk.ValAddress, description Description, newRate *sdk.Dec, newMinSelfDelegation *sdk.Int) MsgEditValidator {
	return MsgEditValidator{
		Description:       description,
		CommissionRate:    newRate,
		ValidatorAddress:  valAddr,
		MinSelfDelegation: newMinSelfDelegation,
	}
}

func (msg MsgEditValidator) Route() string { return RouterKey }

func (msg MsgEditValidator) Type() string { return TypeMsgEditValidator }

func (msg MsgEditValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddress)}
}

// get the bytes for the message signer to sign on
func (msg MsgEditValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgEditValidator) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr
	}
	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}
	if msg.MinSelfDelegation != nil && !msg.MinSelfDelegation.IsPositive() {
		return ErrMinSelfDelegationInvalid
	}
	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(sdk.OneDec()) || msg.CommissionRate.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "commission rate must be between 0 and 1 (inclusive)")
		}
	}

	return nil
}

// MsgDelegate - struct for bonding transactions
type MsgDelegate struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	Amount           sdk.Coin       `json:"amount" yaml:"amount"`
}

func NewMsgDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) MsgDelegate {
	return MsgDelegate{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		Amount:           amount,
	}
}

func (msg MsgDelegate) Route() string { return RouterKey }

func (msg MsgDelegate) Type() string { return TypeMsgDelegate }

func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}

// get the bytes for the message signer to sign on
func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgDelegate) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadDelegationAmount
	}
	return nil
}

//______________________________________________________________________

// MsgDelegate - struct for bonding transactions
type MsgBeginRedelegate struct {
	DelegatorAddress    sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorSrcAddress sdk.ValAddress `json:"validator_src_address" yaml:"validator_src_address"`
	ValidatorDstAddress sdk.ValAddress `json:"validator_dst_address" yaml:"validator_dst_address"`
	Amount              sdk.Coin       `json:"amount" yaml:"amount"`
}

func NewMsgBeginRedelegate(delAddr sdk.AccAddress, valSrcAddr,
	valDstAddr sdk.ValAddress, amount sdk.Coin) MsgBeginRedelegate {

	return MsgBeginRedelegate{
		DelegatorAddress:    delAddr,
		ValidatorSrcAddress: valSrcAddr,
		ValidatorDstAddress: valDstAddr,
		Amount:              amount,
	}
}

func (msg MsgBeginRedelegate) Route() string { return RouterKey }

func (msg MsgBeginRedelegate) Type() string { return TypeMsgBeginRedelegate }

func (msg MsgBeginRedelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddress}
}

// get the bytes for the message signer to sign on
func (msg MsgBeginRedelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgBeginRedelegate) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.ValidatorSrcAddress.Empty() {
		return ErrEmptyValidatorAddr
	}
	if msg.ValidatorDstAddress.Empty() {
		return ErrEmptyValidatorAddr
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadSharesAmount
	}
	return nil
}

// MsgUndelegate - struct for unbonding transactions
type MsgUndelegate struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	Amount           sdk.Coin       `json:"amount" yaml:"amount"`
}

func NewMsgUndelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) MsgUndelegate {
	return MsgUndelegate{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		Amount:           amount,
	}
}

func (msg MsgUndelegate) Route() string { return RouterKey }

func (msg MsgUndelegate) Type() string { return TypeMsgUndelegate }

func (msg MsgUndelegate) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.DelegatorAddress} }

// get the bytes for the message signer to sign on
func (msg MsgUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgUndelegate) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrEmptyDelegatorAddr
	}
	if msg.ValidatorAddress.Empty() {
		return ErrEmptyValidatorAddr
	}
	if msg.Amount.Amount.LTE(sdk.ZeroInt()) {
		return ErrBadSharesAmount
	}
	return nil
}
