package types

import (
	"context"
	"os"
	"testing"

	cosmoslog "cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/types"
	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextMultiStore(t *testing.T) {
	// create test keys
	key := types.NewKVStoreKey("test_key")
	tkey := types.NewTransientStoreKey("test_transient_key")

	// Create test environment
	db := dbm.NewMemDB()
	logger := log.NewLogger(os.Stdout)
	clogger := cosmoslog.NewCustomLogger(*(logger.Impl().(*zerolog.Logger)))
	cms := store.NewCommitMultiStore(db, clogger, metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, storetypes.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	require.NoError(t, err)
	ctx := NewContext(context.Background(), cms, logger)

	// test get MultiStore
	ms := ctx.MultiStore()
	require.NotNil(t, ms, "MultiStore should not be nil")

	// get KVStore from context
	store := ctx.KVStore(key)
	require.NotNil(t, store, "KVStore should not be nil")

	// test set and get value from KVStore
	key1 := []byte("key1")
	value1 := []byte("value1")
	store.Set(key1, value1)

	got := store.Get(key1)
	assert.Equal(t, value1, got, "Retrieved value should match stored value")
}

func TestContextCacheMultiStore(t *testing.T) {
	// create test keys
	key := types.NewKVStoreKey("test_key")
	tkey := types.NewTransientStoreKey("test_transient_key")

	// Create test environment
	db := dbm.NewMemDB()
	logger := log.NewLogger(os.Stdout)
	clogger := cosmoslog.NewCustomLogger(*(logger.Impl().(*zerolog.Logger)))
	cms := store.NewCommitMultiStore(db, clogger, metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, storetypes.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	require.NoError(t, err)
	ctx := NewContext(context.Background(), cms, logger)

	// get original store and set initial data
	store := ctx.KVStore(key)
	key1 := []byte("key1")
	value1 := []byte("value1")
	store.Set(key1, value1)

	// create cache context
	cacheCtx, writeCache := ctx.CacheContext()

	// modify data in cache store
	cacheStore := cacheCtx.KVStore(key)
	value2 := []byte("value2")
	cacheStore.Set(key1, value2)

	// check original store data is not changed
	got := store.Get(key1)
	assert.Equal(t, value1, got, "Original store should maintain old value")

	// check cache store data is changed
	got = cacheStore.Get(key1)
	assert.Equal(t, value2, got, "Cache store should have new value")

	// write cache
	writeCache()

	// check original store data is changed
	got = store.Get(key1)
	assert.Equal(t, value2, got, "Original store should have new value after write cache")
}
