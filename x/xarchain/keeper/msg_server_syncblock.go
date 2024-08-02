package keeper

import (
	"context"

    "xarchain/x/xarchain/types"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)


func (k msgServer) CreateSyncblock(goCtx context.Context,  msg *types.MsgCreateSyncblock) (*types.MsgCreateSyncblockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value already exists
    _, isFound := k.GetSyncblock(
        ctx,
        msg.ChainId,
        )
    if isFound {
        return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
    }

    var syncblock = types.Syncblock{
        Creator: msg.Creator,
        ChainId: msg.ChainId,
        From: msg.From,
        To: msg.To,
        
    }

   k.SetSyncblock(
   		ctx,
   		syncblock,
   	)
	return &types.MsgCreateSyncblockResponse{}, nil
}

func (k msgServer) UpdateSyncblock(goCtx context.Context,  msg *types.MsgUpdateSyncblock) (*types.MsgUpdateSyncblockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, isFound := k.GetSyncblock(
        ctx,
        msg.ChainId,
    )
    if !isFound {
        return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the msg creator is the same as the current owner
    if msg.Creator != valFound.Creator {
        return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

    var syncblock = types.Syncblock{
		Creator: msg.Creator,
		ChainId: msg.ChainId,
        From: msg.From,
		To: msg.To,
		
	}

	k.SetSyncblock(ctx, syncblock)

	return &types.MsgUpdateSyncblockResponse{}, nil
}

func (k msgServer) DeleteSyncblock(goCtx context.Context,  msg *types.MsgDeleteSyncblock) (*types.MsgDeleteSyncblockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // Check if the value exists
    valFound, isFound := k.GetSyncblock(
        ctx,
        msg.ChainId,
    )
    if !isFound {
        return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
    }

    // Checks if the msg creator is the same as the current owner
    if msg.Creator != valFound.Creator {
        return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
    }

	k.RemoveSyncblock(
	    ctx,
	msg.ChainId,
    )

	return &types.MsgDeleteSyncblockResponse{}, nil
}
