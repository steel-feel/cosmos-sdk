package xarchain_test

import (
	"testing"

	keepertest "xarchain/testutil/keeper"
	"xarchain/testutil/nullify"
	xarchain "xarchain/x/xarchain/module"
	"xarchain/x/xarchain/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		Cblock: &types.Cblock{
			Blocknumber: 85,
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.XarchainKeeper(t)
	xarchain.InitGenesis(ctx, k, genesisState)
	got := xarchain.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.Cblock, got.Cblock)
	// this line is used by starport scaffolding # genesis/test/assert
}
