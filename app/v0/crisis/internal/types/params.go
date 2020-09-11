package types

import (
	"fmt"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/params"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName
)

var (
	// ParamStoreKeyConstantFee - key for constant fee parameter
	ParamStoreKeyConstantFee = []byte("ConstantFee")
)

// ParamKeyTable - type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		params.NewParamSetPair(ParamStoreKeyConstantFee, sdk.Coin{}, validateConstantFee),
	)
}

func validateConstantFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("validateConstantFee invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return fmt.Errorf("invalid constant fee: %s", v)
	}

	return nil
}
