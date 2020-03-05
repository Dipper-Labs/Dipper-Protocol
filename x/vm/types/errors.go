package types

import (
	sdk "github.com/Dipper-Protocol/types"
)

const (
	CodespaceType = ModuleName
) 

var (
	ErrNoPayload                = sdk.NewError(CodespaceType, 1, "no payload")
	ErrOutOfGas                 = sdk.NewError(CodespaceType, 2, "out of gas")
	ErrCodeStoreOutOfGas        = sdk.NewError(CodespaceType, 3, "contract creation code storage out of gas")
	ErrDepth                    = sdk.NewError(CodespaceType, 4, "max call depth exceeded")
	ErrTraceLimitReached        = sdk.NewError(CodespaceType, 5, "the number of logs reached the specified limit")
	ErrNoCompatibleInterpreter  = sdk.NewError(CodespaceType, 6, "no compatible interpreter")
	ErrEmptyInputs              = sdk.NewError(CodespaceType, 7, "empty input")
	ErrInsufficientBalance      = sdk.NewError(CodespaceType, 8, "insufficient balance for transfer")
	ErrContractAddressCollision = sdk.NewError(CodespaceType, 9, "contract address collision")
	ErrNoCodeExist              = sdk.NewError(CodespaceType, 10, "no code exists")
	ErrMaxCodeSizeExceeded      = sdk.NewError(CodespaceType, 11, "evm: max code size exceeded")
	ErrWriteProtection          = sdk.NewError(CodespaceType, 12, "vm: write protection")
	ErrReturnDataOutOfBounds    = sdk.NewError(CodespaceType, 13, "evm: return data out of bounds")
	ErrExecutionReverted        = sdk.NewError(CodespaceType, 14, "evm: execution reverted")
	ErrInvalidJump              = sdk.NewError(CodespaceType, 15, "evm: invalid jump destination")
	ErrGasUintOverflow          = sdk.NewError(CodespaceType, 16, "gas uint64 overflow")
)
