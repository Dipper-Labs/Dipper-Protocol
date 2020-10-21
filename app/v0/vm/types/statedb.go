package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/tendermint/tendermint/crypto"

	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/auth/types"
	"github.com/Dipper-Labs/Dipper-Protocol/app/v0/vm/common/math"
	"github.com/Dipper-Labs/Dipper-Protocol/hexutil"
	"github.com/Dipper-Labs/Dipper-Protocol/store/prefix"
	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

var (
	zeroBalance = sdk.ZeroInt().BigInt()
)

type revision struct {
	id           int
	journalIndex int
}

type CommitStateDB struct {
	// TODO: We need to store the context as part of the structure itself opposed
	// to being passed as a parameter (as it should be) in order to implement the
	// StateDB interface. Perhaps there is a better way.
	ctx sdk.Context

	ak auth.AccountKeeper

	storageKey sdk.StoreKey

	// maps that hold 'live' objects, which will get modified while processing a
	// state transition
	stateObjects      map[string]*stateObject
	stateObjectsDirty map[string]struct{}

	// The refund counter, also used by state transitioning.
	refund uint64

	thash, bhash sdk.Hash
	txIndex      int
	logs         map[sdk.Hash][]*Log

	// TODO: Determine if we actually need this as we do not need preimages in
	// the SDK, but it seems to be used elsewhere in Geth.
	preimages map[sdk.Hash][]byte

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memo-ized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionID int

	// mutex for state deep copying
	lock sync.Mutex
}

// NewCommitStateDB returns a reference to a newly initialized CommitStateDB
// which implements Geth's state.StateDB interface.
//
// CONTRACT: Stores used for state must be cache-wrapped as the ordering of the
// key/value space matters in determining the merkle root.
//func NewCommitStateDB(ctx sdk.Context, ak auth.AccountKeeper, storageKey, codeKey sdk.StoreKey) *CommitStateDB {
func NewCommitStateDB(ak auth.AccountKeeper, storageKey sdk.StoreKey) *CommitStateDB {
	return &CommitStateDB{
		ak:                ak,
		storageKey:        storageKey,
		stateObjects:      make(map[string]*stateObject),
		stateObjectsDirty: make(map[string]struct{}),
		logs:              make(map[sdk.Hash][]*Log),
		preimages:         make(map[sdk.Hash][]byte),
		journal:           newJournal(),
	}
}

func NewStateDB(db *CommitStateDB) *CommitStateDB {
	return &CommitStateDB{
		ak:                db.ak,
		storageKey:        db.storageKey,
		stateObjects:      make(map[string]*stateObject),
		stateObjectsDirty: make(map[string]struct{}),
		logs:              make(map[sdk.Hash][]*Log),
		preimages:         make(map[sdk.Hash][]byte),
		journal:           newJournal(),
	}
}

// WithContext returns a Database with an updated sdk context
func (csdb *CommitStateDB) WithContext(ctx sdk.Context) *CommitStateDB {
	csdb.ctx = ctx
	return csdb
}

// ContractCreatedEvent emit event of contract created
// nolint
func (csdb *CommitStateDB) ContractCreatedEvent(addr sdk.AccAddress) {
	csdb.ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeContractCreated,
			sdk.NewAttribute(AttributeKeyAddress, addr.String()),
		),
	})
}

// ContractCalledEvent emit event of contract called
// nolint
func (csdb *CommitStateDB) ContractCalledEvent(addr sdk.AccAddress) {
	csdb.ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			EventTypeContractCalled,
			sdk.NewAttribute(AttributeKeyAddress, addr.String()),
		),
	})
}

func (csdb *CommitStateDB) WithTxHash(txHash []byte) *CommitStateDB {
	csdb.thash = sdk.BytesToHash(txHash)
	return csdb
}

// ----------------------------------------------------------------------------
// Setters
// ----------------------------------------------------------------------------

// SetBalance sets the balance of an account.
func (csdb *CommitStateDB) SetBalance(addr sdk.AccAddress, amount *big.Int) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetBalance(amount)
	}
}

// AddBalance adds amount to the account associated with addr.
func (csdb *CommitStateDB) AddBalance(addr sdk.AccAddress, amount *big.Int) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.AddBalance(amount)
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (csdb *CommitStateDB) SubBalance(addr sdk.AccAddress, amount *big.Int) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SubBalance(amount)
	}
}

// SetNonce sets the nonce (sequence number) of an account.
func (csdb *CommitStateDB) SetNonce(addr sdk.AccAddress, nonce uint64) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetNonce(nonce)
	}
}

// SetState sets the storage state with a key, value pair for an account.
func (csdb *CommitStateDB) SetState(addr sdk.AccAddress, key, value sdk.Hash) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetState(key, value)
	}
}

// SetCode sets the code for a given account.
func (csdb *CommitStateDB) SetCode(addr sdk.AccAddress, code []byte) {
	so := csdb.GetOrNewStateObject(addr)
	if so != nil {
		so.SetCode(sdk.BytesToHash(crypto.Sha256(code)), code)
	}
}

// AddLog adds a new log to the state and sets the log metadata from the state.
func (csdb *CommitStateDB) AddLog(log *Log) {
	csdb.journal.append(addLogChange{txhash: csdb.thash})

	log.TxHash = csdb.thash
	log.BlockHash = csdb.bhash
	log.TxIndex = uint(csdb.txIndex)
	log.Index = csdb.updateLogIndexByOne(false)
	csdb.logs[csdb.thash] = append(csdb.logs[csdb.thash], log)
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (csdb *CommitStateDB) AddPreimage(hash sdk.Hash, preimage []byte) {
	if _, ok := csdb.preimages[hash]; !ok {
		csdb.journal.append(addPreimageChange{hash: hash})

		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		csdb.preimages[hash] = pi
	}
}

// Preimages returns a list of SHA3 preimages that have been submitted.
func (csdb *CommitStateDB) Preimages() map[sdk.Hash][]byte {
	return csdb.preimages
}

// AddRefund adds gas to the refund counter.
func (csdb *CommitStateDB) AddRefund(gas uint64) {
	csdb.journal.append(refundChange{prev: csdb.refund})
	csdb.refund += gas
}

// SubRefund removes gas from the refund counter. It will panic if the refund
// counter goes below zero.
func (csdb *CommitStateDB) SubRefund(gas uint64) {
	csdb.journal.append(refundChange{prev: csdb.refund})
	if gas > csdb.refund {
		panic("refund counter below zero")
	}

	csdb.refund -= gas
}

// ----------------------------------------------------------------------------
// Getters
// ----------------------------------------------------------------------------

// GetBalance retrieves the balance from the given address or 0 if object not
// found.
func (csdb *CommitStateDB) GetBalance(addr sdk.AccAddress) *big.Int {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.Balance()
	}

	return zeroBalance
}

// GetNonce returns the nonce (sequence number) for a given account.
func (csdb *CommitStateDB) GetNonce(addr sdk.AccAddress) uint64 {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.Nonce()
	}

	return 0
}

// TxIndex returns the current transaction index set by Prepare.
func (csdb *CommitStateDB) TxIndex() int {
	return csdb.txIndex
}

// BlockHash returns the current block hash set by Prepare.
func (csdb *CommitStateDB) BlockHash() sdk.Hash {
	return csdb.bhash
}

// GetCode returns the code for a given account.
func (csdb *CommitStateDB) GetCode(addr sdk.AccAddress) []byte {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.Code()
	}

	return nil
}

// GetCodeSize returns the code size for a given account.
func (csdb *CommitStateDB) GetCodeSize(addr sdk.AccAddress) int {
	so := csdb.getStateObject(addr)
	if so == nil {
		return 0
	}

	if so.code != nil {
		return len(so.code)
	}

	// TODO: we may need to cache these lookups directly
	return len(so.Code())
}

// GetCodeHash returns the code hash for a given account.
func (csdb *CommitStateDB) GetCodeHash(addr sdk.AccAddress) sdk.Hash {
	so := csdb.getStateObject(addr)
	if so == nil {
		return sdk.Hash{}
	}

	return sdk.BytesToHash(so.CodeHash())
}

// GetState retrieves a value from the given account's storage store.
func (csdb *CommitStateDB) GetState(addr sdk.AccAddress, hash sdk.Hash) sdk.Hash {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.GetState(hash)
	}

	return sdk.Hash{}
}

// GetCommittedState retrieves a value from the given account's committed
// storage.
func (csdb *CommitStateDB) GetCommittedState(addr sdk.AccAddress, hash sdk.Hash) sdk.Hash {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.GetCommittedState(hash)
	}

	return sdk.Hash{}
}

// GetLogs returns the current logs for a given hash in the state.
func (csdb *CommitStateDB) GetLogs(hash sdk.Hash) (logs []*Log) {
	r, ok := csdb.logs[hash]
	if ok {
		return r
	}

	ctx := csdb.ctx
	store := prefix.NewStore(ctx.KVStore(csdb.storageKey), KeyPrefixLogs)
	d := store.Get(hash.Bytes())
	err := json.Unmarshal(d, &logs)
	if err != nil {
		ctx.Logger().Error(err.Error())
		return
	}

	return
}

// Logs returns all the current logs in the state.
func (csdb *CommitStateDB) Logs() []*Log { // todo: is should get all logs from store?
	logs := make([]*Log, 0, len(csdb.logs))
	for _, lgs := range csdb.logs {
		logs = append(logs, lgs...)
	}

	return logs
}

func (csdb *CommitStateDB) ClearLogs() {
	for k := range csdb.logs {
		delete(csdb.logs, k)
	}
}

// GetRefund returns the current value of the refund counter.
func (csdb *CommitStateDB) GetRefund() uint64 {
	return csdb.refund
}

// HasSuicided returns if the given account for the specified address has been
// killed.
func (csdb *CommitStateDB) HasSuicided(addr sdk.AccAddress) bool {
	so := csdb.getStateObject(addr)
	if so != nil {
		return so.suicided
	}

	return false
}

// ----------------------------------------------------------------------------
// Persistence
// ----------------------------------------------------------------------------

// Commit writes the state to the appropriate KVStores. For each state object
// in the cache, it will either be removed, or have it's code set and/or it's
// state (storage) updated. In addition, the state object (account) itself will
// be written. Finally, the root hash (version) will be returned.
// nolint
func (csdb *CommitStateDB) Commit(deleteEmptyObjects bool) (root sdk.Hash, err error) {
	defer csdb.clearJournalAndRefund()

	// remove dirty state object entries based on the journal
	for addr := range csdb.journal.dirties {
		csdb.stateObjectsDirty[addr] = struct{}{}
	}

	// set the state objects
	for addr, so := range csdb.stateObjects {
		_, isDirty := csdb.stateObjectsDirty[addr]

		switch {
		case so.suicided || (isDirty && deleteEmptyObjects && so.empty()):
			// If the state object has been removed, don't bother syncing it and just
			// remove it from the store.
			csdb.deleteStateObject(so)

		case isDirty:
			// write any contract code associated with the state object
			if so.code != nil && so.dirtyCode {
				so.commitCode()
				so.dirtyCode = false
			}

			// update the object in the KVStore
			csdb.updateStateObject(so)
		}

		delete(csdb.stateObjectsDirty, addr)
	}

	// NOTE: Ethereum returns the trie merkle root here, but as commitment
	// actually happens in the BaseApp at EndBlocker, we do not know the root at
	// this time.
	return
}

func (csdb *CommitStateDB) commitLogs() {
	ctx := csdb.ctx
	store := prefix.NewStore(ctx.KVStore(csdb.storageKey), KeyPrefixLogs)

	hs := make([]string, 0, len(csdb.logs))
	for h := range csdb.logs {
		hs = append(hs, h.String())
	}
	sort.Strings(hs)

	for _, h := range hs {
		hash := sdk.HexToHash(h)
		d, err := json.Marshal(csdb.logs[hash])
		if err != nil {
			ctx.Logger().Error(err.Error())
			continue
		}

		ctx.Logger().Debug(fmt.Sprintf("set logs, txHash: %s, logs: %s", hash.String(), string(d)))
		store.Set(hash.Bytes(), d)
	}
}

func (csdb *CommitStateDB) updateLogIndexByOne(isSubtract bool) uint64 {
	ctx := csdb.ctx
	store := prefix.NewStore(ctx.KVStore(csdb.storageKey), KeyPrefixLogsIndex)

	value := big.NewInt(0)
	if store.Has(LogIndexKey) {
		d := store.Get(LogIndexKey)
		value.SetBytes(d)

		if isSubtract {
			if value.Uint64() == 0 {
				ctx.Logger().Error("current logIndex is 0, can not to be Subtracted")
				return 0
			}
			value.SetUint64(value.Uint64() - 1)
		} else {
			if value.Uint64() == math.MaxUint64 {
				ctx.Logger().Error("current logIndex will out of range")
				return value.Uint64()
			}
			value.SetUint64(value.Uint64() + 1)
		}
	}

	store.Set(LogIndexKey, value.Bytes())

	return value.Uint64()
}

// Finalise finalizes the state objects (accounts) state by setting their state,
// removing the csdb destructed objects and clearing the journal as well as the
// refunds.
func (csdb *CommitStateDB) Finalise(deleteEmptyObjects bool) {
	for addr := range csdb.journal.dirties {
		so, exist := csdb.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx:
			// 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			//
			// That tx goes out of gas, and although the notion of 'touched' does not
			// exist there, the touch-event will still be recorded in the journal.
			// Since ripeMD is a special snowflake, it will persist in the journal even
			// though the journal is reverted. In this special circumstance, it may
			// exist in journal.dirties but not in stateObjects. Thus, we can safely
			// ignore it here.
			continue
		}

		if so.suicided || (deleteEmptyObjects && so.empty()) {
			csdb.deleteStateObject(so)
		} else {
			// Set all the dirty state storage items for the state object in the
			// KVStore and finally set the account in the account mapper.
			so.commitState()
			csdb.updateStateObject(so)
		}

		csdb.stateObjectsDirty[addr] = struct{}{}
	}

	csdb.commitLogs()
	csdb.ClearLogs()

	// invalidate journal because reverting across transactions is not allowed
	csdb.clearJournalAndRefund()
}

// IntermediateRoot returns the current root hash of the state. It is called in
// between transactions to get the root hash that goes into transaction
// receipts.
//
// NOTE: The SDK has not concept or method of getting any intermediate merkle
// root as commitment of the merkle-ized tree doesn't happen until the
// BaseApps' EndBlocker.
func (csdb *CommitStateDB) IntermediateRoot(deleteEmptyObjects bool) sdk.Hash {
	csdb.Finalise(deleteEmptyObjects)

	return sdk.Hash{}
}

// updateStateObject writes the given state object to the store.
func (csdb *CommitStateDB) updateStateObject(so *stateObject) {
	csdb.ak.SetAccount(csdb.ctx, so.account)
}

// deleteStateObject removes the given state object from the state store.
func (csdb *CommitStateDB) deleteStateObject(so *stateObject) {
	so.deleted = true
	csdb.ak.RemoveAccount(csdb.ctx, so.account)
}

// ----------------------------------------------------------------------------
// Snapshotting
// ----------------------------------------------------------------------------

// Snapshot returns an identifier for the current revision of the state.
func (csdb *CommitStateDB) Snapshot() int {
	id := csdb.nextRevisionID
	csdb.nextRevisionID++

	csdb.validRevisions = append(
		csdb.validRevisions,
		revision{
			id:           id,
			journalIndex: csdb.journal.length(),
		},
	)

	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (csdb *CommitStateDB) RevertToSnapshot(revID int) {
	// find the snapshot in the stack of valid snapshots
	idx := sort.Search(len(csdb.validRevisions), func(i int) bool {
		return csdb.validRevisions[i].id >= revID
	})

	if idx == len(csdb.validRevisions) || csdb.validRevisions[idx].id != revID {
		panic(fmt.Errorf("revision ID %v cannot be reverted", revID))
	}

	snapshot := csdb.validRevisions[idx].journalIndex

	// replay the journal to undo changes and remove invalidated snapshots
	csdb.journal.revert(csdb, snapshot)
	csdb.validRevisions = csdb.validRevisions[:idx]
}

// ----------------------------------------------------------------------------
// Auxiliary
// ----------------------------------------------------------------------------

// Database retrieves the low level database supporting the lower level trie
// ops. It is not used in Ethermint, so it returns nil.
//func (csdb *CommitStateDB) Database() ethstate.Database {
//	return nil
//}

// Empty returns whether the state object is either non-existent or empty
// according to the EIP161 specification (balance = nonce = code = 0).
func (csdb *CommitStateDB) Empty(addr sdk.AccAddress) bool {
	so := csdb.getStateObject(addr)
	return so == nil || so.empty()
}

// Exist reports whether the given account address exists in the state. Notably,
// this also returns true for suicided accounts.
func (csdb *CommitStateDB) Exist(addr sdk.AccAddress) bool {
	return csdb.getStateObject(addr) != nil
}

// Error returns the first non-nil error the StateDB encountered.
func (csdb *CommitStateDB) Error() error {
	return csdb.dbErr
}

// Suicide marks the given account as suicided and clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after Suicide.
func (csdb *CommitStateDB) Suicide(addr sdk.AccAddress) bool {
	so := csdb.getStateObject(addr)
	if so == nil {
		return false
	}

	csdb.journal.append(suicideChange{
		account:     &addr,
		prev:        so.suicided,
		prevBalance: sdk.NewIntFromBigInt(so.Balance()),
	})

	so.markSuicided()
	so.SetBalance(new(big.Int))

	return true
}

// Reset clears out all ephemeral state objects from the state db, but keeps
// the underlying account mapper and store keys to avoid reloading data for the
// next operations.
func (csdb *CommitStateDB) Reset(_ sdk.Hash) error {
	csdb.stateObjects = make(map[string]*stateObject)
	csdb.stateObjectsDirty = make(map[string]struct{})
	csdb.thash = sdk.Hash{}
	csdb.bhash = sdk.Hash{}
	csdb.txIndex = 0
	csdb.logs = make(map[sdk.Hash][]*Log)
	csdb.preimages = make(map[sdk.Hash][]byte)

	csdb.clearJournalAndRefund()
	return nil
}

// UpdateAccounts updates the nonce and coin balances of accounts
func (csdb *CommitStateDB) UpdateAccounts() {
	for addr, so := range csdb.stateObjects {
		addr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			continue
		}
		accI := csdb.ak.GetAccount(csdb.ctx, addr)
		acc, ok := accI.(*types.BaseAccount)
		if ok {
			if (so.Balance() != acc.GetCoins().AmountOf(sdk.NativeTokenName).BigInt()) || (so.Nonce() != acc.GetSequence()) {
				// If queried account's balance or nonce are invalid, update the account pointer
				so.account = acc
			}
		}

	}
}

// ClearStateObjects clears cache of state objects to handle account changes outside of the EVM
func (csdb *CommitStateDB) ClearStateObjects() {
	csdb.stateObjects = make(map[string]*stateObject)
	csdb.stateObjectsDirty = make(map[string]struct{})
}

func (csdb *CommitStateDB) clearJournalAndRefund() {
	csdb.journal = newJournal()
	csdb.validRevisions = csdb.validRevisions[:0]
	csdb.refund = 0
}

// Prepare sets the current transaction hash and index and block hash which is
// used when the EVM emits new state logs.
func (csdb *CommitStateDB) Prepare(thash, bhash sdk.Hash, txi int) {
	csdb.thash = thash
	csdb.bhash = bhash
	csdb.txIndex = txi
}

// CreateAccount explicitly creates a state object. If a state object with the
// address already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might
// arise that a contract does the following:
//
//   1. sends funds to sha(account ++ (nonce + 1))
//   2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (csdb *CommitStateDB) CreateAccount(addr sdk.AccAddress) {
	newobj, prevobj := csdb.createObject(addr)
	if prevobj != nil {
		newobj.setBalance(sdk.NewIntFromBigInt(prevobj.Balance()))
	}
}

// Copy creates a deep, independent copy of the state.
//
// NOTE: Snapshots of the copied state cannot be applied to the copy.
func (csdb *CommitStateDB) Copy() *CommitStateDB {
	csdb.lock.Lock()
	defer csdb.lock.Unlock()

	// copy all the basic fields, initialize the memory ones
	state := &CommitStateDB{
		ctx:               csdb.ctx,
		ak:                csdb.ak,
		storageKey:        csdb.storageKey,
		stateObjects:      make(map[string]*stateObject, len(csdb.journal.dirties)),
		stateObjectsDirty: make(map[string]struct{}, len(csdb.journal.dirties)),
		refund:            csdb.refund,
		logs:              make(map[sdk.Hash][]*Log, len(csdb.logs)),
		preimages:         make(map[sdk.Hash][]byte),
		journal:           newJournal(),
	}

	// copy the dirty states, logs, and preimages
	for addr := range csdb.journal.dirties {
		// There is a case where an object is in the journal but not in the
		// stateObjects: OOG after touch on ripeMD prior to Byzantium. Thus, we
		// need to check for nil.
		//
		// Ref: https://github.com/ethereum/go-ethereum/pull/16485#issuecomment-380438527
		if object, exist := csdb.stateObjects[addr]; exist {
			state.stateObjects[addr] = object.deepCopy(state)
			state.stateObjectsDirty[addr] = struct{}{}
		}
	}

	// Above, we don't copy the actual journal. This means that if the copy is
	// copied, the loop above will be a no-op, since the copy's journal is empty.
	// Thus, here we iterate over stateObjects, to enable copies of copies.
	for addr := range csdb.stateObjectsDirty {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = csdb.stateObjects[addr].deepCopy(state)
			state.stateObjectsDirty[addr] = struct{}{}
		}
	}

	// copy logs
	for hash, logs := range csdb.logs {
		cpy := make([]*Log, len(logs))
		for i, l := range logs {
			cpy[i] = new(Log)
			*cpy[i] = *l
		}
		state.logs[hash] = cpy
	}

	// copy pre-images
	for hash, preimage := range csdb.preimages {
		state.preimages[hash] = preimage
	}

	return state
}

// ForEachStorage iterates over each storage items, all invokes the provided
// callback on each key, value pair .
func (csdb *CommitStateDB) ForEachStorage(addr sdk.AccAddress, cb func(key, value sdk.Hash) bool) error {
	so := csdb.getStateObject(addr)
	if so == nil {
		return nil
	}

	store := prefix.NewStore(csdb.ctx.KVStore(csdb.storageKey), KeyPrefixStorage)
	iter := sdk.KVStorePrefixIterator(store, so.Address().Bytes())

	for ; iter.Valid(); iter.Next() {
		key := sdk.BytesToHash(iter.Key())
		value := iter.Value()

		if value, dirty := so.dirtyStorage[key]; dirty {
			cb(key, value)
			continue
		}

		cb(key, sdk.BytesToHash(value))
	}

	iter.Close()
	return nil
}

// GetOrNewStateObject retrieves a state object or create a new state object if
// nil.
func (csdb *CommitStateDB) GetOrNewStateObject(addr sdk.AccAddress) StateObject {
	so := csdb.getStateObject(addr)
	if so == nil || so.deleted {
		so, _ = csdb.createObject(addr)
	}

	return so
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (csdb *CommitStateDB) createObject(addr sdk.AccAddress) (newObj, prevObj *stateObject) {
	prevObj = csdb.getStateObject(addr)

	acc := csdb.ak.NewAccountWithAddress(csdb.ctx, sdk.AccAddress(addr.Bytes()))
	newObj = newObject(csdb, acc)
	newObj.setNonce(0) // sets the object to dirty

	if prevObj == nil {
		csdb.journal.append(createObjectChange{account: &addr})
	} else {
		csdb.journal.append(resetObjectChange{prev: prevObj})
	}

	csdb.setStateObject(newObj)
	return newObj, prevObj
}

// setError remembers the first non-nil error it is called with.
func (csdb *CommitStateDB) setError(err error) {
	if csdb.dbErr == nil {
		csdb.dbErr = err
	}
}

// getStateObject attempts to retrieve a state object given by the address.
// Returns nil and sets an error if not found.
func (csdb *CommitStateDB) getStateObject(addr sdk.AccAddress) (stateObject *stateObject) {
	// prefer 'live' (cached) objects
	if so := csdb.stateObjects[addr.String()]; so != nil {
		if so.deleted {
			return nil
		}

		return so
	}

	// otherwise, attempt to fetch the account from the account mapper
	acc := csdb.ak.GetAccount(csdb.ctx, addr.Bytes())
	if acc == nil {
		csdb.setError(fmt.Errorf("no account found for address: %X", addr.Bytes()))
		return nil
	}

	// insert the state object into the live set
	so := newObject(csdb, acc)
	csdb.setStateObject(so)

	return so
}

func (csdb *CommitStateDB) setStateObject(so *stateObject) {
	csdb.stateObjects[so.Address().String()] = so
}

func (csdb *CommitStateDB) ExportStateObjects(params QueryStateParams) (sos SOs) {
	var so SO

	for _, stateObject := range csdb.stateObjects {
		if params.ContractOnly {
			if len(stateObject.CodeHash()) == 0 {
				continue
			}
		}

		so.Address = stateObject.address
		so.BaseAccount = *stateObject.account
		so.OriginStorage = stateObject.originStorage.Copy()
		so.DirtyStorage = stateObject.dirtyStorage.Copy()
		so.DirtyCode = stateObject.dirtyCode
		so.Suicided = stateObject.suicided
		so.Deleted = stateObject.deleted
		if !params.ShowCode {
			so.Code = nil
		} else {
			so.Code = append(so.Code[:], stateObject.code...)
		}

		sos = append(sos, so)
	}

	return sos
}

// ExportState used to export vm state to genesis file
func (csdb *CommitStateDB) ExportState() (s GenesisState) {
	s.Storage = csdb.exportStorage()
	s.Codes = csdb.exportCodes()
	s.VMLogs = csdb.exportLogs()
	return
}

// ImportState used to import vm state from genesis file
func (csdb *CommitStateDB) ImportState(s GenesisState) {
	err := csdb.importCodes(s.Codes)
	if err != nil {
		panic(err)
	}

	err = csdb.importStorage(s.Storage)
	if err != nil {
		panic(err)
	}

	err = csdb.importLogs(s.VMLogs)
	if err != nil {
		panic(err)
	}
}

func (csdb *CommitStateDB) exportCodes() map[string]sdk.Code {
	store := prefix.NewStore(csdb.ctx.KVStore(csdb.storageKey), KeyPrefixCode)
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	codes := make(map[string]sdk.Code, 10240)
	for ; iter.Valid(); iter.Next() {
		codes[hexutil.Encode(iter.Key())] = iter.Value()
	}

	return codes
}

func (csdb *CommitStateDB) importCodes(codes map[string]sdk.Code) error {
	store := prefix.NewStore(csdb.ctx.KVStore(csdb.storageKey), KeyPrefixCode)

	for k, v := range codes {
		K, err := hexutil.Decode(k)
		if err != nil {
			return err
		}
		store.Set(K, v)
	}

	return nil
}

func (csdb *CommitStateDB) exportStorage() (gs []Storage) {
	store := prefix.NewStore(csdb.ctx.KVStore(csdb.storageKey), KeyPrefixStorage)
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		gs = append(gs, Storage{
			Key:   iter.Key(),
			Value: iter.Value(),
		})
	}

	return
}

func (csdb *CommitStateDB) importStorage(gs []Storage) error {
	store := prefix.NewStore(csdb.ctx.KVStore(csdb.storageKey), KeyPrefixStorage)

	for _, item := range gs {
		store.Set(item.Key, item.Value)
	}

	return nil
}

func (csdb *CommitStateDB) exportLogs() (vmLogs VMLogs) {
	store := prefix.NewStore(csdb.ctx.KVStore(csdb.storageKey), KeyPrefixLogs)
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	vmLogs.Logs = make(map[string]string, 10240)

	for ; iter.Valid(); iter.Next() {
		vmLogs.Logs[hexutil.Encode(iter.Key())] = string(iter.Value())
	}

	logIndexStore := prefix.NewStore(csdb.ctx.KVStore(csdb.storageKey), KeyPrefixLogsIndex)
	vmLogs.LogIndex = -1
	if logIndexStore.Has(LogIndexKey) {
		vmLogs.LogIndex = big.NewInt(0).SetBytes(logIndexStore.Get(LogIndexKey)).Int64()
	}

	return
}

func (csdb *CommitStateDB) importLogs(vmLogs VMLogs) error {
	store := prefix.NewStore(csdb.ctx.KVStore(csdb.storageKey), KeyPrefixLogs)

	if vmLogs.LogIndex == -1 {
		return nil
	}

	for txHashStr, logs := range vmLogs.Logs {
		txHash, err := hexutil.Decode(txHashStr)
		if err != nil {
			return err
		}

		store.Set(txHash, []byte(logs))
	}

	logIndexStore := prefix.NewStore(csdb.ctx.KVStore(csdb.storageKey), KeyPrefixLogsIndex)

	logIndex := big.NewInt(vmLogs.LogIndex)
	logIndexStore.Set(LogIndexKey, logIndex.Bytes())

	return nil
}

// for simulation
func (csdb *CommitStateDB) GetAllHotContractAddrs() (accs []sdk.AccAddress) {
	for _, obj := range csdb.stateObjects {
		if len(obj.code) != 0 || len(obj.account.CodeHash) != 0 {
			accs = append(accs, obj.address)
		}
	}
	return
}
