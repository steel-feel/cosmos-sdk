package keeper

import (
	"context"

	"xarchain/x/xarchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateTask(goCtx context.Context, msg *types.MsgCreateTask) (*types.MsgCreateTaskResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	var task = types.Task{
		Creator: msg.Creator,
		Title:   msg.Title,
		Status:  "proposed",
		Abci:    "none",
	}

	id := k.AppendTask(ctx, task)

	return &types.MsgCreateTaskResponse{
		Id: id,
	}, nil
}
