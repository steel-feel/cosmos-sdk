package keeper_test

import (
    "strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"xarchain/x/xarchain/types"
	"xarchain/testutil/nullify"
	keepertest "xarchain/testutil/keeper"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestSyncblockQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.XarchainKeeper(t)
	msgs := createNSyncblock(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetSyncblockRequest
		response *types.QueryGetSyncblockResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetSyncblockRequest{
			    ChainId: msgs[0].ChainId,
                
			},
			response: &types.QueryGetSyncblockResponse{Syncblock: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetSyncblockRequest{
			    ChainId: msgs[1].ChainId,
                
			},
			response: &types.QueryGetSyncblockResponse{Syncblock: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetSyncblockRequest{
			    ChainId:strconv.Itoa(100000),
                
			},
			err:     status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Syncblock(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestSyncblockQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.XarchainKeeper(t)
	msgs := createNSyncblock(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllSyncblockRequest {
		return &types.QueryAllSyncblockRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.SyncblockAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Syncblock), step)
			require.Subset(t,
            	nullify.Fill(msgs),
            	nullify.Fill(resp.Syncblock),
            )
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.SyncblockAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Syncblock), step)
			require.Subset(t,
            	nullify.Fill(msgs),
            	nullify.Fill(resp.Syncblock),
            )
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.SyncblockAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Syncblock),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.SyncblockAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
