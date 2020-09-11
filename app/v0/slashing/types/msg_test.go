package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Dipper-Labs/Dipper-Protocol/types"
)

func TestMsgUnjailGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("abcd")
	msg := NewMsgUnjail(sdk.ValAddress(addr))
	bytes := msg.GetSignBytes()
	require.Equal(
		t,
		`{"type":"dip/MsgUnjail","value":{"address":"dipvaloper1v93xxeqcn5rzv"}}`,
		string(bytes),
	)
}
