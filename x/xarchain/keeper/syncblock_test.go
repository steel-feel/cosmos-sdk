package keeper_test

import (
	"context"
	"strconv"
	"testing"

	"xarchain/x/xarchain/keeper"
	"xarchain/x/xarchain/types"
	keepertest "xarchain/testutil/keeper"
	"xarchain/testutil/nullify"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNSyncblock(keeper keeper.Keeper, ctx context.Context, n int) []types.Syncblock {
	items := make([]types.Syncblock, n)
	for i := range items {
		items[i].ChainId = strconv.Itoa(i)
        
		keeper.SetSyncblock(ctx, items[i])
	}
	return items
}

func TestSyncblockGet(t *testing.T) {
	keeper, ctx := keepertest.XarchainKeeper(t)
	items := createNSyncblock(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetSyncblock(ctx,
		    item.ChainId,
            
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestSyncblockRemove(t *testing.T) {
	keeper, ctx := keepertest.XarchainKeeper(t)
	items := createNSyncblock(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveSyncblock(ctx,
		    item.ChainId,
            
		)
		_, found := keeper.GetSyncblock(ctx,
		    item.ChainId,
            
		)
		require.False(t, found)
	}
}

func TestSyncblockGetAll(t *testing.T) {
	keeper, ctx := keepertest.XarchainKeeper(t)
	items := createNSyncblock(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllSyncblock(ctx)),
	)
}
