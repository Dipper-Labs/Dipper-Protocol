package types

import sdk "github.com/Dipper-Labs/Dipper-Protocol/types"

type Vesting struct {
	Address   sdk.AccAddress `json:"address"`
	Amount    sdk.Coins      `json:"Amount"`
	StartTime int64          `json:"start_time"`
	EndTime   int64          `json:"end_time"`
}
