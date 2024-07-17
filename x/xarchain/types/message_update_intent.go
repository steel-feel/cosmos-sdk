package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateIntent{}

func NewMsgUpdateIntent(creator string, status string, id uint64) *MsgUpdateIntent {
	return &MsgUpdateIntent{
		Creator: creator,
		Status:  status,
		Id:      id,
	}
}

func (msg *MsgUpdateIntent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
