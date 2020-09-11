package types

import (
	sdkerrors "github.com/Dipper-Labs/Dipper-Protocol/types/errors"
)

var (
	ErrNoSender         = sdkerrors.New(ModuleName, 1, "sender address is empty")
	ErrUnknownInvariant = sdkerrors.New(ModuleName, 2, "unknown invariant")
)
