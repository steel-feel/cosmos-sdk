package keeper_test

import (
    "strconv"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

    keepertest "xarchain/testutil/keeper"
    "xarchain/x/xarchain/keeper"
    "xarchain/x/xarchain/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestSyncblockMsgServerCreate(t *testing.T) {
	k, ctx := keepertest.XarchainKeeper(t)
	srv := keeper.NewMsgServerImpl(k)
	creator := "A"
	for i := 0; i < 5; i++ {
		expected := &types.MsgCreateSyncblock{Creator: creator,
		    ChainId: strconv.Itoa(i),
            
		}
		_, err := srv.CreateSyncblock(ctx, expected)
		require.NoError(t, err)
		rst, found := k.GetSyncblock(ctx,
		    expected.ChainId,
            
		)
		require.True(t, found)
		require.Equal(t, expected.Creator, rst.Creator)
	}
}

func TestSyncblockMsgServerUpdate(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgUpdateSyncblock
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgUpdateSyncblock{Creator: creator,
			    ChainId: strconv.Itoa(0),
                
			},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgUpdateSyncblock{Creator: "B",
			    ChainId: strconv.Itoa(0),
                
			},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgUpdateSyncblock{Creator: creator,
			    ChainId: strconv.Itoa(100000),
                
			},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.XarchainKeeper(t)
			srv := keeper.NewMsgServerImpl(k)
			expected := &types.MsgCreateSyncblock{Creator: creator,
			    ChainId: strconv.Itoa(0),
                
			}
			_, err := srv.CreateSyncblock(ctx, expected)
			require.NoError(t, err)

			_, err = srv.UpdateSyncblock(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				rst, found := k.GetSyncblock(ctx,
				    expected.ChainId,
                    
				)
				require.True(t, found)
				require.Equal(t, expected.Creator, rst.Creator)
			}
		})
	}
}

func TestSyncblockMsgServerDelete(t *testing.T) {
	creator := "A"

	tests := []struct {
		desc    string
		request *types.MsgDeleteSyncblock
		err     error
	}{
		{
			desc:    "Completed",
			request: &types.MsgDeleteSyncblock{Creator: creator,
			    ChainId: strconv.Itoa(0),
                
			},
		},
		{
			desc:    "Unauthorized",
			request: &types.MsgDeleteSyncblock{Creator: "B",
			    ChainId: strconv.Itoa(0),
                
			},
			err:     sdkerrors.ErrUnauthorized,
		},
		{
			desc:    "KeyNotFound",
			request: &types.MsgDeleteSyncblock{Creator: creator,
			    ChainId: strconv.Itoa(100000),
                
			},
			err:     sdkerrors.ErrKeyNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.XarchainKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			_, err := srv.CreateSyncblock(ctx, &types.MsgCreateSyncblock{Creator: creator,
			    ChainId: strconv.Itoa(0),
                
			})
			require.NoError(t, err)
			_, err = srv.DeleteSyncblock(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				_, found := k.GetSyncblock(ctx,
				    tc.request.ChainId,
                    
				)
				require.False(t, found)
			}
		})
	}
}
