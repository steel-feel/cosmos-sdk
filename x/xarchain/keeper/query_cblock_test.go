package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "xarchain/testutil/keeper"
	"xarchain/testutil/nullify"
	"xarchain/x/xarchain/types"
)

func TestCblockQuery(t *testing.T) {
	keeper, ctx := keepertest.XarchainKeeper(t)
	item := createTestCblock(keeper, ctx)
	tests := []struct {
		desc     string
		request  *types.QueryGetCblockRequest
		response *types.QueryGetCblockResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetCblockRequest{},
			response: &types.QueryGetCblockResponse{Cblock: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Cblock(ctx, tc.request)
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
