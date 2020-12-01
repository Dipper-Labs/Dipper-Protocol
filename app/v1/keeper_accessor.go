package v1

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/gov"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/staking"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply"
)

// GovKeeper return govKeeper
func (p *ProtocolV1) GovKeeper() gov.Keeper {
	return p.govKeeper
}

// SetGovKeeper set govKeeper
func (p *ProtocolV1) SetGovKeeper(gk gov.Keeper) {
	p.govKeeper = gk
}

// StakingKeeper return stakingKeeper
func (p *ProtocolV1) StakingKeeper() staking.Keeper {
	return p.stakingKeeper
}

// AccountKeeper return accountKeeper
func (p *ProtocolV1) AccountKeeper() auth.AccountKeeper {
	return p.accountKeeper
}

// SupplyKeeper return supplyKeeper
func (p *ProtocolV1) SupplyKeeper() supply.Keeper {
	return p.supplyKeeper
}
