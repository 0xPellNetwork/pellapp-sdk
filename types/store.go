package types

import (
	"context"
	"fmt"
	"sync"

	storetypes "cosmossdk.io/store/types"
)

// StoreProvider is an interface that provides access to the commit and query multi-store.
type StoreProvider interface {
	CommitMultiStore() storetypes.CommitMultiStore
	QueryMultiStore() storetypes.MultiStore
}

// QueryManager is an interface that defines methods for querying data from the store.
type QueryManager interface {
	Get(ctx context.Context, storeKey storetypes.StoreKey, key []byte) ([]byte, error)
}

// TxManager is an interface that extends QueryManager and defines methods for modifying data in the store.
type TxManager interface {
	QueryManager
	Set(ctx context.Context, storeKey storetypes.StoreKey, key, value []byte) (storetypes.CommitID, error)
	Delete(ctx context.Context, storeKey storetypes.StoreKey, key []byte) (storetypes.CommitID, error)
}

// appDataManager is a struct that implements the TxManager interface.
type appDataManager struct {
	provider StoreProvider
	mtx      sync.RWMutex
}

// Get retrieves a value from the store using the provided store key and key.
func (m *appDataManager) Get(ctx context.Context, storeKey storetypes.StoreKey, key []byte) ([]byte, error) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	store := m.provider.QueryMultiStore().GetKVStore(storeKey)
	if store == nil {
		return nil, fmt.Errorf("store %s not found", storeKey)
	}
	return store.Get(key), nil
}

// Set stores a value in the store using the provided store key and key, returning the commit ID and error.
func (m *appDataManager) Set(ctx context.Context,
	storeKey storetypes.StoreKey,
	key, value []byte,
) (storetypes.CommitID, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	store := m.provider.CommitMultiStore().GetKVStore(storeKey)
	if store == nil {
		return storetypes.CommitID{}, fmt.Errorf("store %s not found", storeKey)
	}
	store.Set(key, value)
	return m.provider.CommitMultiStore().Commit(), nil
}

// Delete removes a value from the store using the provided store key and key, returning the commit ID and error.
func (m *appDataManager) Delete(ctx context.Context, storeKey storetypes.StoreKey, key []byte) (storetypes.CommitID, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	store := m.provider.CommitMultiStore().GetKVStore(storeKey)
	if store == nil {
		return storetypes.CommitID{}, fmt.Errorf("store %s not found", storeKey)
	}
	store.Delete(key)

	return m.provider.CommitMultiStore().Commit(), nil
}

// AppTxManager is a struct that embeds appDataManager and provides transaction management functionality.
type AppTxManager struct {
	*appDataManager
}

// NewAppTxManager creates a new instance of AppTxManager with the provided StoreProvider.
type AppQueryManager struct {
	*appDataManager
}

// NewAppTxManager creates a new instance of AppTxManager with the provided StoreProvider.
func NewAppTxManager(provider StoreProvider) *AppTxManager {
	return &AppTxManager{
		appDataManager: &appDataManager{
			provider: provider,
		},
	}
}

// AppQueryManager is a struct that embeds appDataManager and provides query management functionality.
func NewAppQueryManager(provider StoreProvider) *AppQueryManager {
	return &AppQueryManager{
		appDataManager: &appDataManager{
			provider: provider,
		},
	}
}
