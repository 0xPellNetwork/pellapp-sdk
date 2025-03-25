package testutil

import (
	"context"
	"os"
	"testing"

	cosmoslog "cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
)

// DefaultContext creates a sdktypes.Context with a fresh MemDB that can be used in tests.
func DefaultContext(key, tkey storetypes.StoreKey) sdktypes.Context {
	db := dbm.NewMemDB()
	logger := log.NewLogger(os.Stdout)
	clogger := cosmoslog.NewCustomLogger(*(logger.Impl().(*zerolog.Logger)))
	cms := store.NewCommitMultiStore(db, clogger, metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, storetypes.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	ctx := sdktypes.NewContext(context.Background(), cms, log.NewNopLogger())

	return ctx
}

// DefaultContextWithKeys creates a sdktypes.Context with a fresh MemDB, mounting the providing keys for usage in the multistore.
// This function is intended to be used for testing purposes only.
func DefaultContextWithKeys(
	keys map[string]*storetypes.KVStoreKey,
	transKeys map[string]*storetypes.TransientStoreKey,
	memKeys map[string]*storetypes.MemoryStoreKey,
) sdktypes.Context {
	db := dbm.NewMemDB()
	logger := log.NewLogger(os.Stdout)
	clogger := cosmoslog.NewCustomLogger(*(logger.Impl().(*zerolog.Logger)))
	cms := store.NewCommitMultiStore(db, clogger, metrics.NewNoOpMetrics())

	for _, key := range keys {
		cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	}

	for _, tKey := range transKeys {
		cms.MountStoreWithDB(tKey, storetypes.StoreTypeTransient, db)
	}

	for _, memkey := range memKeys {
		cms.MountStoreWithDB(memkey, storetypes.StoreTypeMemory, db)
	}

	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	return sdktypes.NewContext(context.Background(), cms, log.NewNopLogger())
}

type TestContext struct {
	Ctx sdktypes.Context
	DB  *dbm.MemDB
	CMS store.CommitMultiStore
}

func DefaultContextWithDB(t testing.TB, key, tkey storetypes.StoreKey) TestContext {
	db := dbm.NewMemDB()
	logger := log.NewLogger(os.Stdout)
	clogger := cosmoslog.NewCustomLogger(*(logger.Impl().(*zerolog.Logger)))
	cms := store.NewCommitMultiStore(db, clogger, metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, storetypes.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	assert.NoError(t, err)

	ctx := sdktypes.NewContext(context.Background(), cms, log.NewNopLogger())

	return TestContext{ctx, db, cms}
}
