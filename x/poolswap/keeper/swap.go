package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerror "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/okex/okchain/x/poolswap/types"
)

// IsTokenExist check token is exist
func (k Keeper) IsTokenExist(ctx sdk.Context, token string) error {
	isExist := k.tokenKeeper.TokenExist(ctx, token)
	if !isExist {
		return sdkerror.Wrap(sdkerror.ErrInternal, "Failed: token does not exist")
	}

	t := k.tokenKeeper.GetTokenInfo(ctx, token)
	if t.Type == types.GenerateTokenType {
		return sdkerror.Wrap(sdkerror.ErrInvalidCoins, "Failed to create exchange with pool token")
	}
	return nil

}
