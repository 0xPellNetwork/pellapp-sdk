package store

import (
	storetypes "cosmossdk.io/store/types"
)

// StoreManager encapsulates all storage-related operations
type StoreManager struct {
	cms storetypes.CommitMultiStore
	qms storetypes.MultiStore
}

// NewStoreManager creates a new storage manager
func NewStoreManager(cms storetypes.CommitMultiStore, qms storetypes.MultiStore) *StoreManager {
	return &StoreManager{
		cms: cms,
		qms: qms,
	}
}

// HasQueryStore checks if the query store is set
func (sm *StoreManager) HasQueryStore() bool {
	return sm.qms != nil
}

// GetKVStore returns the KV store for a specific store key.
// If useQueryStore is true, it will use the query multistore (if available).
func (sm *StoreManager) GetKVStore(key storetypes.StoreKey, useQueryStore bool) storetypes.KVStore {
	if useQueryStore && sm.HasQueryStore() {
		return sm.qms.GetKVStore(key)
	}
	return sm.cms.GetKVStore(key)
}

// Set writes the given key-value pair to the store specified by storeKey.
func (sm *StoreManager) Set(storeKey storetypes.StoreKey, key, value []byte) {
	store := sm.cms.GetKVStore(storeKey)
	store.Set(key, value)
}

// Get retrieves a value for the given key from the store specified by storeKey.
// If useQueryStore is true, it will attempt to use the query store for better performance.
func (sm *StoreManager) Get(storeKey storetypes.StoreKey, key []byte, useQueryStore bool) []byte {
	store := sm.GetKVStore(storeKey, useQueryStore)
	return store.Get(key)
}

// Delete removes a key-value pair from the store specified by storeKey.
func (sm *StoreManager) Delete(storeKey storetypes.StoreKey, key []byte) {
	store := sm.cms.GetKVStore(storeKey)
	store.Delete(key)
}

// Has checks if a key exists in the store specified by storeKey.
// If useQueryStore is true, it will attempt to use the query store.
func (sm *StoreManager) Has(storeKey storetypes.StoreKey, key []byte, useQueryStore bool) bool {
	store := sm.GetKVStore(storeKey, useQueryStore)
	return store.Has(key)
}

// Iterator returns an iterator over a domain of keys in the store specified by storeKey.
// If useQueryStore is true, it will attempt to use the query store.
func (sm *StoreManager) Iterator(storeKey storetypes.StoreKey, start, end []byte, useQueryStore bool) storetypes.Iterator {
	store := sm.GetKVStore(storeKey, useQueryStore)
	return store.Iterator(start, end)
}

// BatchOperation executes multiple store operations atomically.
// The provided function is given a cached multi-store to perform operations on.
// If the function returns an error, none of the operations will be committed.
// Otherwise, all operations will be written to the main store.
func (sm *StoreManager) BatchOperation(fn func(ms storetypes.CacheMultiStore) error) error {
	// Create a cached store for atomic operations
	cms := sm.cms.CacheMultiStore()

	// Execute operations on the cached store
	err := fn(cms)
	if err != nil {
		return err
	}

	// If successful, write all operations to the main store
	cms.Write()
	return nil
}

// ForEach iterates over all key-value pairs in a store and applies the handler function.
// Returns early if the handler returns false.
func (sm *StoreManager) ForEach(storeKey storetypes.StoreKey, handler func(key, value []byte) bool, useQueryStore bool) {
	store := sm.GetKVStore(storeKey, useQueryStore)
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		if !handler(iter.Key(), iter.Value()) {
			break
		}
	}
}

// FindByPrefix retrieves all values with keys matching the given prefix.
func (sm *StoreManager) FindByPrefix(storeKey storetypes.StoreKey, prefix []byte, useQueryStore bool) [][]byte {
	store := sm.GetKVStore(storeKey, useQueryStore)

	iter := store.Iterator(prefix, storetypes.PrefixEndBytes(prefix))
	defer iter.Close()

	var results [][]byte
	for ; iter.Valid(); iter.Next() {
		results = append(results, iter.Value())
	}

	return results
}

// Commit commits all changes to the underlying CommitMultiStore.
// This is useful after performing multiple Set/Delete operations.
func (sm *StoreManager) Commit() storetypes.CommitID {
	return sm.cms.Commit()
}
