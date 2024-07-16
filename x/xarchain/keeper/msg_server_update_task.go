package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"xarchain/x/xarchain/types"
)

func (k msgServer) UpdateTask(goCtx context.Context, msg *types.MsgUpdateTask) (*types.MsgUpdateTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var task = types.Task{
		Creator: msg.Creator,
		Id:      msg.Id,
		Title:   msg.Title,
		Status:  "picked",
	}
	val, found := k.GetTask(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetTask(ctx, task)
	return &types.MsgUpdateTaskResponse{}, nil
}
