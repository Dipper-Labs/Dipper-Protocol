package dipperProtocol

import (
	sdk "github.com/Dipper-Protocol/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	WhoisRecords []byte `json:"whois_records"`
	//BillBank BillBank `json:"bill_bank"`
}

func NewGenesisState(whoIsRecords []byte) GenesisState {
	return GenesisState{WhoisRecords: nil}//, BillBank: BillBank{}}
}

func ValidateGenesis(data GenesisState) error {
	//for _, record := range data.WhoisRecords {
	//	if record.Owner == nil {
	//		return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Owner", record.Value)
	//	}
	//	if record.Value == "" {
	//		return fmt.Errorf("invalid WhoisRecord: Owner: %s. Error: Missing Value", record.Owner)
	//	}
	//	if record.Price == nil {
	//		return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Price", record.Value)
	//	}
	//}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		WhoisRecords: []byte{},
		//BillBank: types.BillBank{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	//for _, record := range data.WhoisRecords {
		//keeper.SetWhois(ctx, record.Value, record)
	//}
	//keeper.SetBillBank(ctx, data.BillBank)
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	//var records []Whois
	//iterator := k.GetNamesIterator(ctx)
	//for ; iterator.Valid(); iterator.Next() {
	//
	//	name := string(iterator.Key())
	//	whois := k.GetWhois(ctx, name)
	//	records = append(records, whois)
	//}
	//billBank := types.NewBillBank()
	return GenesisState{WhoisRecords: []byte{}}//, BillBank:billBank}
}
