package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/params"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// Parameter store keys
var (
	KeyMintDenom           = []byte("MintDenom")
	KeyInflationRateChange = []byte("InflationRateChange")
	KeyInflationMax        = []byte("InflationMax")
	KeyInflationMin        = []byte("InflationMin")
	KeyGoalBonded          = []byte("GoalBonded")
	KeyBlocksPerYear       = []byte("BlocksPerYear")
	KeyMaxProvisions       = []byte("MaxProvisions")
)

// mint parameters
type Params struct {
	MintDenom           string  `json:"mint_denom" yaml:"mint_denom"`                       // type of coin to mint
	InflationRateChange sdk.Dec `json:"inflation_rate_change" yaml:"inflation_rate_change"` // maximum annual change in inflation rate
	InflationMax        sdk.Dec `json:"inflation_max" yaml:"inflation_max"`                 // maximum inflation rate
	InflationMin        sdk.Dec `json:"inflation_min" yaml:"inflation_min"`                 // minimum inflation rate
	GoalBonded          sdk.Dec `json:"goal_bonded" yaml:"goal_bonded"`                     // goal of percent bonded atoms
	BlocksPerYear       uint64  `json:"blocks_per_year" yaml:"blocks_per_year"`             // expected blocks per year
	MaxProvisions       sdk.Dec `json:"max_provisions" yaml:"max_provisions"`
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, inflationRateChange, inflationMax,
	inflationMin, goalBonded sdk.Dec, blocksPerYear uint64, maxProvisions sdk.Dec) Params {

	return Params{
		MintDenom:           mintDenom,
		InflationRateChange: inflationRateChange,
		InflationMax:        inflationMax,
		InflationMin:        inflationMin,
		GoalBonded:          goalBonded,
		BlocksPerYear:       blocksPerYear,
		MaxProvisions:       maxProvisions,
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:           sdk.DefaultBondDenom,
		InflationRateChange: sdk.NewDecWithPrec(6, 2),
		InflationMax:        sdk.NewDecWithPrec(10, 2),
		InflationMin:        sdk.NewDecWithPrec(4, 2),
		GoalBonded:          sdk.NewDecWithPrec(67, 2),
		BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
		MaxProvisions:       sdk.NewDec(350000000).Mul(sdk.NewDec(1000000000000)),
	}
}

// validate params
func (p Params) Validate() error {
	if p.InflationMax.LT(p.InflationMin) {
		return fmt.Errorf("mint parameter Max inflation must be greater than or equal to min inflation")
	}

	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationRateChange(p.InflationRateChange); err != nil {
		return err
	}
	if err := validateInflationMax(p.InflationMax); err != nil {
		return err
	}
	if err := validateInflationMin(p.InflationMin); err != nil {
		return err
	}
	if err := validateGoalBonded(p.GoalBonded); err != nil {
		return err
	}
	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}
	if p.InflationMax.LT(p.InflationMin) {
		return fmt.Errorf(
			"max inflation (%s) must be greater than or equal to min inflation (%s)",
			p.InflationMax, p.InflationMin,
		)
	}

	return nil
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		params.NewParamSetPair(KeyInflationRateChange, &p.InflationRateChange, validateInflationRateChange),
		params.NewParamSetPair(KeyInflationMax, &p.InflationMax, validateInflationMax),
		params.NewParamSetPair(KeyInflationMin, &p.InflationMin, validateInflationMin),
		params.NewParamSetPair(KeyGoalBonded, &p.GoalBonded, validateGoalBonded),
		params.NewParamSetPair(KeyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
		params.NewParamSetPair(KeyMaxProvisions, &p.MaxProvisions, validateMaxProvisions),
	}
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}
	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}

func validateInflationRateChange(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation rate change cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("inflation rate change too large: %s", v)
	}

	return nil
}

func validateInflationMax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("max inflation cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("max inflation too large: %s", v)
	}

	return nil
}

func validateInflationMin(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min inflation cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("min inflation too large: %s", v)
	}

	return nil
}

func validateGoalBonded(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("goal bonded cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("goal bonded must be <= 1: %s", v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}

func validateMaxProvisions(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("max provisions cannot be negative: %s", v)
	}

	return nil
}
