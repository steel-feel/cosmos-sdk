package keeper

import (
	"context"

	"xarchain/x/xarchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateIntent(goCtx context.Context, msg *types.MsgCreateIntent) (*types.MsgCreateIntentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	var intent = types.Intent{
		To:      msg.To,
		Value:   msg.Value,
		From:    msg.From,
		Data:    msg.Data,
		Status:  "proposed",
		ChainId: msg.ChainId,
		Intentid: msg.Intentid,
	}

	id := k.AppendIntent(ctx, intent)

	return &types.MsgCreateIntentResponse{
		Id: id,
	}, nil
}
