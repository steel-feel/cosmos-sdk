package keeper

import (
	"context"

	"xarchain/x/xarchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateIntent(goCtx context.Context, msg *types.MsgUpdateIntent) (*types.MsgUpdateIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgUpdateIntentResponse{}, nil
}
