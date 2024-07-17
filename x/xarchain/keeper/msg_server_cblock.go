package keeper

import (
	"context"

	"xarchain/x/xarchain/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateCblock(goCtx context.Context, msg *types.MsgCreateCblock) (*types.MsgCreateCblockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetCblock(ctx)
	if isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "already set")
	}

	var cblock = types.Cblock{
		Creator:     msg.Creator,
		Blocknumber: msg.Blocknumber,
	}

	k.SetCblock(
		ctx,
		cblock,
	)
	return &types.MsgCreateCblockResponse{}, nil
}

func (k msgServer) UpdateCblock(goCtx context.Context, msg *types.MsgUpdateCblock) (*types.MsgUpdateCblockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	_, isFound := k.GetCblock(ctx)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}

	// // Checks if the msg creator is the same as the current owner
	// if msg.Creator != valFound.Creator {
	// 	return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	// }

	var cblock = types.Cblock{
		Creator:    msg.Creator,
		Blocknumber: msg.Blocknumber,
	}

	k.SetCblock(ctx, cblock)

	return &types.MsgUpdateCblockResponse{}, nil
}

func (k msgServer) DeleteCblock(goCtx context.Context, msg *types.MsgDeleteCblock) (*types.MsgDeleteCblockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	_, isFound := k.GetCblock(ctx)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "not set")
	}

	// Checks if the msg creator is the same as the current owner
	// if msg.Creator != valFound.Creator {
	// 	return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	// }

	k.RemoveCblock(ctx)

	return &types.MsgDeleteCblockResponse{}, nil
}
