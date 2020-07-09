package params

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkparamskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	sdkparamstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/okex/okchain/x/params/types"
)

// Keeper is the struct of params keeper
type Keeper struct {
	cdc codec.Marshaler
	sdkparamskeeper.Keeper
	// the reference to the Paramstore to get and set gov specific params
	paramSpace sdkparamstypes.Subspace
	// the reference to the DelegationSet and ValidatorSet to get information about validators and delegators
	sk StakingKeeper
	// the reference to the CoinKeeper to modify balances
	ck BankKeeper
	// the reference to the GovKeeper to insert waiting queue
	gk GovKeeper
}

// NewKeeper creates a new instance of params keeper
func NewKeeper(cdc codec.Marshaler, key *sdk.KVStoreKey, tkey *sdk.TransientStoreKey) (
	k Keeper) {
	k = Keeper{
		Keeper: sdkparamskeeper.NewKeeper(cdc, key, tkey),
	}
	k.cdc = cdc
	k.paramSpace = k.Subspace(DefaultParamspace).WithKeyTable(types.ParamKeyTable())
	return k
}

// SetStakingKeeper hooks the staking keeper into params keeper
func (keeper *Keeper) SetStakingKeeper(sk StakingKeeper) {
	keeper.sk = sk
}

// SetBankKeeper hooks the bank keeper into params keeper
func (keeper *Keeper) SetBankKeeper(ck BankKeeper) {
	keeper.ck = ck
}

// SetGovKeeper hooks the gov keeper into params keeper
func (keeper *Keeper) SetGovKeeper(gk GovKeeper) {
	keeper.gk = gk
}

// SetParams sets the params into the store
func (keeper *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	keeper.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the params info from the store
func (keeper Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	keeper.paramSpace.GetParamSet(ctx, &params)
	return params
}
