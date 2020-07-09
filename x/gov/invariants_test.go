package gov

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/okex/okchain/x/gov/keeper"
	"github.com/okex/okchain/x/gov/types"
)

func TestModuleAccountInvariant(t *testing.T) {
	ctx, _, gk, _, crisisKeeper := keeper.CreateTestInput(t, false, 1000)
	govHandler := NewHandler(gk)

	initialDeposit := sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 50)}
	content := types.NewTextProposal("Test", "description")
	newProposalMsg, err := NewMsgSubmitProposal(content, initialDeposit, keeper.Addrs[0])
	require.Nil(t, err)
	res, err := govHandler(ctx, newProposalMsg)
	require.Nil(t, err)
	proposalID := types.GetProposalIDFromBytes(res.Data)

	newDepositMsg := NewMsgDeposit(keeper.Addrs[0], proposalID,
		sdk.DecCoins{sdk.NewInt64DecCoin(sdk.DefaultBondDenom, 100)})
	res, err = govHandler(ctx, newDepositMsg)
	require.Nil(t, err)

	invariant := ModuleAccountInvariant(gk)
	_, broken := invariant(ctx)
	require.False(t, broken)

	// todo: check diff after RegisterInvariants
	RegisterInvariants(&crisisKeeper, gk)
}
