package upgrade

import (
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/staking"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/staking/exported"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/upgrade/types"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
	sdkerrors "github.com/Dipper-Labs/Dipper-Protocol/types/errors"
)

// Keeper defines the keeper struct of the upgrade store
type Keeper struct {
	storeKey       sdk.StoreKey
	cdc            *codec.Codec
	protocolKeeper sdk.ProtocolKeeper
	sk             staking.Keeper
}

// NewKeeper creates a new upgrade keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, protocolKeeper sdk.ProtocolKeeper, sk staking.Keeper) Keeper {
	keeper := Keeper{
		key,
		cdc,
		protocolKeeper,
		sk,
	}
	return keeper
}

// AddNewVersionInfo adds new version info
func (k Keeper) AddNewVersionInfo(ctx sdk.Context, versionInfo types.VersionInfo) {
	kvStore := ctx.KVStore(k.storeKey)

	versionInfoBytes, err := k.cdc.MarshalBinaryLengthPrefixed(versionInfo)
	if err != nil {
		panic(err)
	}
	kvStore.Set(types.GetProposalIDKey(versionInfo.UpgradeInfo.ProposalID), versionInfoBytes)

	proposalIDBytes, err := k.cdc.MarshalBinaryLengthPrefixed(versionInfo.UpgradeInfo.ProposalID)
	if err != nil {
		panic(err)
	}

	if versionInfo.Success {
		kvStore.Set(types.GetSuccessVersionKey(versionInfo.UpgradeInfo.Protocol.Version), proposalIDBytes)
	} else {
		kvStore.Set(types.GetFailedVersionKey(versionInfo.UpgradeInfo.Protocol.Version, versionInfo.UpgradeInfo.ProposalID), proposalIDBytes)
	}
}

// SetSignal sets signal for upgrade
func (k Keeper) SetSignal(ctx sdk.Context, protocol uint64, address string) {
	kvStore := ctx.KVStore(k.storeKey)
	cmsgBytes, err := k.cdc.MarshalBinaryLengthPrefixed(true)
	if err != nil {
		panic(err)
	}
	kvStore.Set(types.GetSignalKey(protocol, address), cmsgBytes)
}

// GetSignal gets signal
func (k Keeper) GetSignal(ctx sdk.Context, protocol uint64, address string) bool {
	kvStore := ctx.KVStore(k.storeKey)
	flagBytes := kvStore.Get(types.GetSignalKey(protocol, address))
	if flagBytes != nil {
		var flag bool
		err := k.cdc.UnmarshalBinaryLengthPrefixed(flagBytes, &flag)
		if err != nil {
			panic(err)
		}
		return true
	}
	return false
}

// DeleteSignal removes signal
func (k Keeper) DeleteSignal(ctx sdk.Context, protocol uint64, address string) bool {
	if ok := k.GetSignal(ctx, protocol, address); ok {
		kvStore := ctx.KVStore(k.storeKey)
		kvStore.Delete(types.GetSignalKey(protocol, address))
		return true
	}
	return false
}

// IterateBondedValidatorsByPower iterates bonded validators by power
func (k Keeper) IterateBondedValidatorsByPower(ctx sdk.Context,
	fn func(index int64, validator exported.ValidatorI) (stop bool)) {
	k.sk.IterateBondedValidatorsByPower(ctx, fn)
}

// GetCurrentVersion gets current version
func (k Keeper) GetCurrentVersion(ctx sdk.Context) uint64 {
	return k.protocolKeeper.GetCurrentVersion(ctx)
}

// SetAppUpgradeConfig sets app upgrade config for test
// deprecated
func (k Keeper) SetAppUpgradeConfig(ctx sdk.Context, proposalID, version, upgradeHeight uint64, software string,
) error {
	if _, found := k.protocolKeeper.GetUpgradeConfig(ctx); found {
		return sdkerrors.Wrap(sdkerrors.ErrInternal, "failed. an app upgrade config is existed, only one entry is permitted")
	}

	appUpgradeConfig := sdk.NewUpgradeConfig(
		proposalID,
		sdk.NewProtocolDefinition(version, software, upgradeHeight, sdk.NewDecWithPrec(7, 1)),
	)
	k.protocolKeeper.SetUpgradeConfig(ctx, appUpgradeConfig)
	return nil
}
