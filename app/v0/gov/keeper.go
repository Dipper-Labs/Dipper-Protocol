package gov

import (
	"fmt"
	"time"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/gov/types"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/guardian"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/supply/exported"
	"github.com/Dipper-Labs/Dipper-Protocol/codec"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"

	"github.com/tendermint/tendermint/libs/log"
)

// Keeper defines the governance module Keeper
type Keeper struct {
	// The reference to the Paramstore to get and set gov specific params
	paramSpace ParamSubspace

	// The SupplyKeeper to reduce the supply of the network
	supplyKeeper SupplyKeeper

	// The reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	sk StakingKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Proposal router
	router Router

	gk guardian.Keeper
	pk sdk.ProtocolKeeper
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
//
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace ParamSubspace,
	supplyKeeper SupplyKeeper, sk StakingKeeper, gk guardian.Keeper, pk sdk.ProtocolKeeper,
) Keeper {

	// ensure governance module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// It is vital to seal the governance proposal router here as to not allow
	// further handlers to be registered after the keeper is created since this
	// could create invalid or non-deterministic behavior.

	return Keeper{
		storeKey:     key,
		paramSpace:   paramSpace,
		supplyKeeper: supplyKeeper,
		sk:           sk,
		gk:           gk,
		cdc:          cdc,
		pk:           pk,
	}
}

func (keeper *Keeper) SetRouter(rtr Router) {
	keeper.router = rtr
	rtr.Seal()
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

// GetGovernanceAccount returns the governance ModuleAccount
func (keeper Keeper) GetGovernanceAccount(ctx sdk.Context) exported.ModuleAccountI {
	return keeper.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// Params

// Returns the current DepositParams from the global param store
func (keeper Keeper) GetDepositParams(ctx sdk.Context) DepositParams {
	var depositParams DepositParams
	keeper.paramSpace.Get(ctx, ParamStoreKeyDepositParams, &depositParams)
	return depositParams
}

// Returns the current VotingParams from the global param store
func (keeper Keeper) GetVotingParams(ctx sdk.Context) VotingParams {
	var votingParams VotingParams
	keeper.paramSpace.Get(ctx, ParamStoreKeyVotingParams, &votingParams)
	return votingParams
}

// Returns the current TallyParam from the global param store
func (keeper Keeper) GetTallyParams(ctx sdk.Context) TallyParams {
	var tallyParams TallyParams
	keeper.paramSpace.Get(ctx, ParamStoreKeyTallyParams, &tallyParams)
	return tallyParams
}

func (keeper Keeper) setDepositParams(ctx sdk.Context, depositParams DepositParams) {
	keeper.paramSpace.Set(ctx, ParamStoreKeyDepositParams, &depositParams)
}

func (keeper Keeper) setVotingParams(ctx sdk.Context, votingParams VotingParams) {
	keeper.paramSpace.Set(ctx, ParamStoreKeyVotingParams, &votingParams)
}

func (keeper Keeper) setTallyParams(ctx sdk.Context, tallyParams TallyParams) {
	keeper.paramSpace.Set(ctx, ParamStoreKeyTallyParams, &tallyParams)
}

// ProposalQueues

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(types.ActiveProposalQueueKey(proposalID, endTime), bz)
}

// RemoveFromActiveProposalQueue removes a proposalID from the Active Proposal Queue
func (keeper Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.ActiveProposalQueueKey(proposalID, endTime))
}

// InsertInactiveProposalQueue Inserts a ProposalID into the inactive proposal queue at endTime
func (keeper Keeper) InsertInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(types.InactiveProposalQueueKey(proposalID, endTime), bz)
}

// RemoveFromInactiveProposalQueue removes a proposalID from the Inactive Proposal Queue
func (keeper Keeper) RemoveFromInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.InactiveProposalQueueKey(proposalID, endTime))
}

// Iterators

// IterateProposals iterates over the all the proposals and performs a callback function
func (keeper Keeper) IterateProposals(ctx sdk.Context, cb func(proposal types.Proposal) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.ProposalsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var proposal types.Proposal
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &proposal)

		if cb(proposal) {
			break
		}
	}
}

// IterateActiveProposalsQueue iterates over the proposals in the active proposal queue
// and performs a callback function
func (keeper Keeper) IterateActiveProposalsQueue(ctx sdk.Context, endTime time.Time, cb func(proposal types.Proposal) (stop bool)) {
	iterator := keeper.ActiveProposalQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types.SplitActiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateInactiveProposalsQueue iterates over the proposals in the inactive proposal queue
// and performs a callback function
func (keeper Keeper) IterateInactiveProposalsQueue(ctx sdk.Context, endTime time.Time, cb func(proposal types.Proposal) (stop bool)) {
	iterator := keeper.InactiveProposalQueueIterator(ctx, endTime)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID, _ := types.SplitInactiveProposalQueueKey(iterator.Key())
		proposal, found := keeper.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// IterateAllDeposits iterates over the all the stored deposits and performs a callback function
func (keeper Keeper) IterateAllDeposits(ctx sdk.Context, cb func(deposit types.Deposit) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DepositsKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// IterateDeposits iterates over the all the proposals deposits and performs a callback function
func (keeper Keeper) IterateDeposits(ctx sdk.Context, proposalID uint64, cb func(deposit types.Deposit) (stop bool)) {
	iterator := keeper.GetDepositsIterator(ctx, proposalID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}

// IterateAllVotes iterates over the all the stored votes and performs a callback function
func (keeper Keeper) IterateAllVotes(ctx sdk.Context, cb func(vote types.Vote) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// IterateVotes iterates over the all the proposals votes and performs a callback function
func (keeper Keeper) IterateVotes(ctx sdk.Context, proposalID uint64, cb func(vote types.Vote) (stop bool)) {
	iterator := keeper.GetVotesIterator(ctx, proposalID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// ActiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Active Queue that expire by endTime
func (keeper Keeper) ActiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(ActiveProposalQueuePrefix, sdk.PrefixEndBytes(types.ActiveProposalByTimeKey(endTime)))
}

// InactiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Inactive Queue that expire by endTime
func (keeper Keeper) InactiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(InactiveProposalQueuePrefix, sdk.PrefixEndBytes(types.InactiveProposalByTimeKey(endTime)))
}

func (keeper Keeper) SoftwareUpgradeProposalExist(ctx sdk.Context) bool {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.SoftwareUpgradeKey)
	return len(bz) != 0
}

func (keeper Keeper) SoftwareUpgradeSet(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)
	store.Set(types.SoftwareUpgradeKey, types.SoftwareUpgradeValue)
}

func (keeper Keeper) SoftwareUpgradeClear(ctx sdk.Context) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.SoftwareUpgradeKey)
}

func (keeper Keeper) Getcdc() *codec.Codec {
	return keeper.cdc
}
