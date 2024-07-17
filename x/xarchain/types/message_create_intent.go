package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateIntent{}

func NewMsgCreateIntent(creator string, from string, to string, data string, value string, chainId string) *MsgCreateIntent {
	return &MsgCreateIntent{
		Creator: creator,
		From:    from,
		To:      to,
		Data:    data,
		Value:   value,
		ChainId: chainId,
	}
}

func (msg *MsgCreateIntent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
