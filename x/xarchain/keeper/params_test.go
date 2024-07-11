package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "xarchain/testutil/keeper"
	"xarchain/x/xarchain/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.XarchainKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
