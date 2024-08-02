package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateSyncblock{}

func NewMsgCreateSyncblock(
    creator string,
    chainId string,
    from uint64,
    to uint64,
    
) *MsgCreateSyncblock {
  return &MsgCreateSyncblock{
		Creator : creator,
		ChainId: chainId,
		From: from,
        To: to,
        
	}
}

func (msg *MsgCreateSyncblock) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  	if err != nil {
  		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  	}
  return nil
}

var _ sdk.Msg = &MsgUpdateSyncblock{}

func NewMsgUpdateSyncblock(
    creator string,
    chainId string,
    from uint64,
    to uint64,
    
) *MsgUpdateSyncblock {
  return &MsgUpdateSyncblock{
		Creator: creator,
        ChainId: chainId,
        From: from,
        To: to,
        
	}
}

func (msg *MsgUpdateSyncblock) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  }
   return nil
}

var _ sdk.Msg = &MsgDeleteSyncblock{}

func NewMsgDeleteSyncblock(
    creator string,
    chainId string,
    
) *MsgDeleteSyncblock {
  return &MsgDeleteSyncblock{
		Creator: creator,
		ChainId: chainId,
        
	}
}

func (msg *MsgDeleteSyncblock) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Creator)
  if err != nil {
    return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
  }
  return nil
}
