package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "xarchain/testutil/keeper"
	"xarchain/testutil/nullify"
	"xarchain/x/xarchain/keeper"
	"xarchain/x/xarchain/types"
)

func createTestCblock(keeper keeper.Keeper, ctx context.Context) types.Cblock {
	item := types.Cblock{}
	keeper.SetCblock(ctx, item)
	return item
}

func TestCblockGet(t *testing.T) {
	keeper, ctx := keepertest.XarchainKeeper(t)
	item := createTestCblock(keeper, ctx)
	rst, found := keeper.GetCblock(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestCblockRemove(t *testing.T) {
	keeper, ctx := keepertest.XarchainKeeper(t)
	createTestCblock(keeper, ctx)
	keeper.RemoveCblock(ctx)
	_, found := keeper.GetCblock(ctx)
	require.False(t, found)
}
