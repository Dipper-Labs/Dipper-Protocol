package params

import (
	"fmt"

	sdk "github.com/Dipper-Protocol/types"
	govtypes "github.com/Dipper-Protocol/x/gov/types"
)

func NewParamChangeProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) sdk.Error {
		switch c := content.(type) {
		case ParameterChangeProposal:
			return handleParameterChangeProposal(ctx, k, c)

		default:
			errMsg := fmt.Sprintf("unrecognized param proposal content type: %T", c)
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}

func handleParameterChangeProposal(ctx sdk.Context, k Keeper, p ParameterChangeProposal) sdk.Error {
	for _, c := range p.Changes {
		ss, ok := k.GetSubspace(c.Subspace)
		if !ok {
			return sdk.ErrUnknownRequest(c.Subspace)
			//return sdkerrors.Wrap(ErrUnknownSubspace, c.Subspace)
		}

		k.Logger(ctx).Info(
			fmt.Sprintf("attempt to set new parameter value; key: %s, value: %s", c.Key, c.Value),
		)

		if err := ss.Update(ctx, []byte(c.Key), []byte(c.Value)); err != nil {
			errMsg := fmt.Sprintf("key: %s, value: %s, err: %s", c.Key, c.Value, err.Error())
			return sdk.ErrUnknownRequest(errMsg)
			//return sdkerrors.Wrapf(ErrSettingParameter, "key: %s, value: %s, err: %s", c.Key, c.Value, err.Error())
		}
	}

	return nil
}
