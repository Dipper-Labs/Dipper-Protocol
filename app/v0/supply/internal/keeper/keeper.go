package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply/exported"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply/internal/types"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

// Keeper of the supply store
type Keeper struct {
	cdc       *codec.Codec
	storeKey  sdk.StoreKey
	ak        types.AccountKeeper
	bk        types.BankKeeper
	permAddrs map[string]types.PermissionsForAddress
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, ak types.AccountKeeper, bk types.BankKeeper, maccPerms map[string][]string) Keeper {
	// set the addresses
	permAddrs := make(map[string]types.PermissionsForAddress)
	for name, perms := range maccPerms {
		permAddrs[name] = types.NewPermissionsForAddress(name, perms)
	}

	return Keeper{
		cdc:       cdc,
		storeKey:  key,
		ak:        ak,
		bk:        bk,
		permAddrs: permAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

// GetSupply retrieves the Supply from store
func (k Keeper) GetSupply(ctx sdk.Context) (supply exported.SupplyI) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(SupplyKey)
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

// SetSupply sets the Supply to store
func (k Keeper) SetSupply(ctx sdk.Context, supply exported.SupplyI) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(SupplyKey, b)
}

// ValidatePermissions validates that the module account has been granted
// permissions within its set of allowed permissions.
func (k Keeper) ValidatePermissions(macc exported.ModuleAccountI) error {
	permAddr := k.permAddrs[macc.GetName()]
	for _, perm := range macc.GetPermissions() {
		if !permAddr.HasPermission(perm) {
			return fmt.Errorf("invalid module permission %s", perm)
		}
	}
	return nil
}

// for Vesting
// GetVesting get vesting
func (k Keeper) GetVesting(ctx sdk.Context, addr sdk.AccAddress) (exist bool, vesting types.Vesting) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(VestingStoreKey(addr))
	if bz == nil {
		return false, vesting
	}

	err := k.cdc.UnmarshalBinaryBare(bz, &vesting)
	if err != nil {
		panic(err)
	}

	return true, vesting
}

// SetAccount implements sdk.AccountKeeper.
func (k Keeper) SetVesting(ctx sdk.Context, vesting types.Vesting) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryBare(vesting)
	if err != nil {
		panic(err)
	}
	store.Set(VestingStoreKey(vesting.Address), bz)
}

// GetAllVestings returns all vestings in the supplyKeeper.
func (k Keeper) GetAllVestings(ctx sdk.Context) (vestings []types.Vesting) {
	appendVesting := func(vesting types.Vesting) (stop bool) {
		vestings = append(vestings, vesting)
		return false
	}
	k.IterateVestings(ctx, appendVesting)
	return
}

// RemoveAccount removes an account for the account mapper store.
func (k Keeper) RemoveVesting(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(VestingStoreKey(address))
}

// IterateVestings
func (k Keeper) IterateVestings(ctx sdk.Context, process func(types.Vesting) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, VestingStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		var vesting types.Vesting
		err := k.cdc.UnmarshalBinaryBare(val, &vesting)
		if err != nil {
			panic(err)
		}

		if process(vesting) {
			return
		}
		iter.Next()
	}
}
