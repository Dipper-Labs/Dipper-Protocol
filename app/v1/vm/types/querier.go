package types

import (
	"fmt"

	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

const (
	QueryParameters = "params"
	QueryState      = "state"
	QueryCode       = "code"
	QueryStorage    = "storage"
	QueryTxLogs     = "logs"
	EstimateGas     = "estimate_gas"
	QueryCall       = "call"
)

// QueryLogsResult - for query logs
type QueryLogsResult struct {
	Logs []*Log `json:"logs"`
}

func (q QueryLogsResult) String() string {
	return fmt.Sprintf("%+v", q.Logs)
}

// QueryStorageResult - for query storage
type QueryStorageResult struct {
	Value sdk.Hash `json:"value"`
}

func (q QueryStorageResult) String() string {
	return q.Value.String()
}

// SimulationResult - for Gas Estimate
type SimulationResult struct {
	Code   uint8
	ErrMsg string
	Gas    uint64
	Res    string
}

func (r SimulationResult) String() string {
	return fmt.Sprintf("Code = %d\nErrMsg = %s\nGas = %d\nRes = %s", r.Code, r.ErrMsg, r.Gas, r.Res)
}

// QueryStateParams - for query vm db state
type QueryStateParams struct {
	ShowCode     bool `json:"show_code" yaml:"show_code"`
	ContractOnly bool `json:"contract_only" yaml:"contract_only"`
}
