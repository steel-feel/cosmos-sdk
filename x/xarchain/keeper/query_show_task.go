package keeper

import (
	"context"

	"xarchain/x/xarchain/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ShowTask(goCtx context.Context, req *types.QueryShowTaskRequest) (*types.QueryShowTaskResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	task, found := k.GetTask(ctx, req.Id)
    if !found {
        return nil, sdkerrors.ErrKeyNotFound
    }

	return &types.QueryShowTaskResponse{Task: task}, nil
}
