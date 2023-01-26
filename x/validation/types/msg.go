package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	types2 "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	"github.com/golang/protobuf/proto"
)

type MsgCustom struct {
	Sender sdk.AccAddress
	Value  int
}

func NewMsgCustom(sender sdk.AccAddress, value int) *MsgCustom {
	return &MsgCustom{
		Sender: sender,
		Value:  value,
	}
}

func (msg MsgCustom) Route() string {
	return "custom"
}

func (msg MsgCustom) Type() string {
	return "custom_message"
}

func (msg MsgCustom) ValidateBasic() error {
	if msg.Sender.Empty() {
		return types2.ErrInvalidAccountAddress
	}
	if msg.Value <= 0 {
		return sdkerrors.New("error", 0, "")
	}

	return nil
}

func (msg MsgCustom) GetSignBytes() []byte {
	cdcAmino := codec.NewLegacyAmino()
	bz := cdcAmino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCustom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// Reset resets the message
func (msg *MsgCustom) Reset() {
	*msg = MsgCustom{}
}

// String returns a human-readable representation of the message
func (msg *MsgCustom) String() string {
	return proto.CompactTextString(msg)
}

// ProtoMessage implements the protobuf proto.Message interface
func (msg MsgCustom) ProtoMessage() {}
