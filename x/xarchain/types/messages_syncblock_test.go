package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"xarchain/testutil/sample"
)

func TestMsgCreateSyncblock_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateSyncblock
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateSyncblock{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateSyncblock{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateSyncblock_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateSyncblock
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateSyncblock{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateSyncblock{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgDeleteSyncblock_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteSyncblock
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteSyncblock{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteSyncblock{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
