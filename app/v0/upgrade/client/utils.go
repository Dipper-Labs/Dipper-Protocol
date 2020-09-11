package client

import (
	"fmt"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/upgrade/types"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

type UpgradeInfoOutput struct {
	CurrentVersion    types.VersionInfo `json:"current_version"`
	LastFailedVersion uint64            `json:"last_failed_version"`
	UpgradeInProgress sdk.UpgradeConfig `json:"upgrade_in_progress"`
}

func NewUpgradeInfoOutput(currentVersion types.VersionInfo, lastFailedVersion uint64, upgradeInProgress sdk.UpgradeConfig) UpgradeInfoOutput {
	return UpgradeInfoOutput{
		currentVersion,
		lastFailedVersion,
		upgradeInProgress,
	}
}

func (p UpgradeInfoOutput) String() string {
	success := "fail"
	if p.CurrentVersion.Success {
		success = "success"
	}
	return fmt.Sprintf(`Upgrade Info:
  Current Version[%v]:  %s     
  Last Failed Version:  %v
  Upgrade In Progress:  %s`,
		success, p.CurrentVersion.UpgradeInfo, p.LastFailedVersion, p.UpgradeInProgress)
}
