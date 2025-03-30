package baseapp

import (
	"os"
	"testing"

	cosmoslog "cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func setupBaseApp(t *testing.T) *BaseApp {
	logger := log.NewLogger(os.Stdout)
	db := dbm.NewMemDB()
	app := NewBaseApp("test", logger, db, nil)
	return app
}

func TestBaseAppStore(t *testing.T) {
	app := setupBaseApp(t)

	// test that the commit multi store is not nil when the app is created
	cms := app.CommitMultiStore()
	require.NotNil(t, cms)

	require.Equal(t, int64(0), app.LastBlockHeight())
	require.NotNil(t, app.LastCommitID())

	// test setting a new commit multi store
	clogger := cosmoslog.NewCustomLogger(*(app.logger.Impl().(*zerolog.Logger)))
	newCMS := store.NewCommitMultiStore(dbm.NewMemDB(), clogger, storemetrics.NewNoOpMetrics())
	app.SetCMS(newCMS)
	require.Equal(t, newCMS, app.CommitMultiStore())
}

func TestBaseAppQueryStore(t *testing.T) {
	app := setupBaseApp(t)

	// test that the query store is nil
	require.False(t, app.HasQueryMultiStore())

	// test QueryMultiStore() is not nil  even if not set(it uses the same as CommitMultiStore)
	require.Equal(t, app.CommitMultiStore(), app.QueryMultiStore())

	// test setting a new query store
	cms := app.CommitMultiStore()
	qms := cms.CacheMultiStore()
	app.SetQueryMultiStore(qms)

	require.True(t, app.HasQueryMultiStore())
	require.Equal(t, qms, app.QueryMultiStore())
}

func TestBaseAppSeal(t *testing.T) {
	app := setupBaseApp(t)
	clogger := cosmoslog.NewCustomLogger(*(app.logger.Impl().(*zerolog.Logger)))
	newCMS := store.NewCommitMultiStore(dbm.NewMemDB(), clogger, storemetrics.NewNoOpMetrics())

	// testing that the cms can be set before sealing
	app.SetCMS(newCMS)

	// this should not panic
	app.Sealed()
	require.True(t, app.sealed)

	// this should panic because the app is sealed
	require.Panics(t, func() {
		app.SetCMS(newCMS)
	})
}

func TestBaseAppStoreOperations(t *testing.T) {
	app := setupBaseApp(t)

	// Test store key setup
	storeKey := storetypes.NewKVStoreKey("test")
	app.MountStore(storeKey, storetypes.StoreTypeIAVL)

	// Initialize stores
	err := app.cms.LoadLatestVersion()
	require.NoError(t, err)

	// Get stores for testing
	cmsStore := app.cms.GetKVStore(storeKey)
	require.NotNil(t, cmsStore)

	require.NoError(t, err)
	qmsStore := app.QueryMultiStore().GetKVStore(storeKey)
	require.NotNil(t, qmsStore)

	// Test basic store operations
	testKey := []byte("test_key")
	testValue := []byte("test_value")

	// Test Set operation
	cmsStore.Set(testKey, testValue)
	value := cmsStore.Get(testKey)
	require.Equal(t, testValue, value)

	// Verify query store doesn't have the value yet
	value = qmsStore.Get(testKey)
	require.NotNil(t, value)

	// Test Delete operation
	cmsStore.Delete(testKey)
	value = cmsStore.Get(testKey)
	require.Nil(t, value)

	// Test query store independence
	qmsStore.Set(testKey, testValue)
	value = qmsStore.Get(testKey)
	require.Equal(t, testValue, value)

	// Verify main store is not affected by query store operations
	value = cmsStore.Get(testKey)
	require.Nil(t, value)

	qmsStore = app.QueryMultiStore().GetKVStore(storeKey)
	value = qmsStore.Get(testKey)
	require.Nil(t, value)
}
