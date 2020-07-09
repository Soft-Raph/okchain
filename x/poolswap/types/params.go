package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/params"
)

// FeeRate defines swap fee rate
var (
	defaultFeeRate = sdk.NewDecWithPrec(3, 3)
)

// Default parameter namespace
const (
	DefaultParamspace = ModuleName

	GenerateTokenType = 2
)

// Parameter store keys
var (
	KeyFeeRate = []byte("FeeRate")
)

// ParamKeyTable for swap module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for swap at genesis
type Params struct {
	FeeRate sdk.Dec `json:"fee_rate"`
}

// NewParams creates a new Params object
func NewParams(feeRate sdk.Dec) Params {
	return Params{
		FeeRate: feeRate,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Poolswap Params:
  TradeFeeRate: %s`, p.FeeRate)
}

func tmpValidate(value interface{}) error {
	return nil
}
// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyFeeRate, &p.FeeRate, tmpValidate},
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(defaultFeeRate)
}
