package keeper_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	keepertest "xarchain/testutil/keeper"
	"xarchain/x/xarchain/keeper"
	"xarchain/x/xarchain/types"
)

func TestCblockMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.XarchainKeeper(t)
	srv := keeper.NewMsgServerImpl(k)

	expected := &types.MsgCreateCblock{}
	_, err := srv.CreateCblock(ctx, expected)
	require.NoError(t, err)
	_, found := k.GetCblock(ctx)
	require.True(t, found)
	// require.Equal(t, expected.Creator, rst.Creator)
}

func TestCblockMsgServerUpdate(t *testing.T) {

	tests := []struct {
		desc    string
		request *types.MsgUpdateCblock
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateCblock{},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateCblock{},
			err:     sdkerrors.ErrUnauthorized,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.XarchainKeeper(t)
			srv := keeper.NewMsgServerImpl(k)
			expected := &types.MsgCreateCblock{}
			_, err := srv.CreateCblock(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateCblock(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetCblock(ctx)
				require.True(t, found)
				// require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestCblockMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteCblock
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeleteCblock{Creator: creator},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeleteCblock{Creator: "B"},
			err:     sdkerrors.ErrUnauthorized,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.XarchainKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.CreateCblock(ctx, &types.MsgCreateCblock{})
			require.NoError(t, err)
			_, err = srv.DeleteCblock(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetCblock(ctx)
				require.False(t, found)
			}
		})
	}
}
