package types

import (
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

type VersionInfo struct {
	UpgradeInfo sdk.UpgradeConfig `json:"upgrade_info"`
	Success     bool              `json:"success"`
}

func NewVersionInfo(upgradeConfig sdk.UpgradeConfig, success bool) VersionInfo {
	return VersionInfo{
		upgradeConfig,
		success,
	}
}
