package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateCblock{}

func NewMsgCreateCblock(creator string, blocknumber uint64) *MsgCreateCblock {
	return &MsgCreateCblock{
		Creator:     creator,
		Blocknumber: blocknumber,
	}
}

func (msg *MsgCreateCblock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateCblock{}

func NewMsgUpdateCblock(creator string, blocknumber uint64) *MsgUpdateCblock {
	return &MsgUpdateCblock{
		Creator:     creator,
		Blocknumber: blocknumber,
	}
}

func (msg *MsgUpdateCblock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteCblock{}

func NewMsgDeleteCblock(creator string) *MsgDeleteCblock {
	return &MsgDeleteCblock{
		Creator: creator,
	}
}

func (msg *MsgDeleteCblock) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
