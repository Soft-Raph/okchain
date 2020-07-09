package gov

import (
	"encoding/json"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"

	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okchain/x/common/version"
	"github.com/okex/okchain/x/gov/client"
	sdkgovcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	"github.com/okex/okchain/x/gov/client/rest"
	"github.com/okex/okchain/x/gov/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines app module basics object
type AppModuleBasic struct {
	proposalHandlers []client.ProposalHandler // proposal handlers which live in governance cli and rest
}

// NewAppModuleBasic creates a new AppModuleBasic object
func NewAppModuleBasic(proposalHandlers ...client.ProposalHandler) AppModuleBasic {
	return AppModuleBasic{
		proposalHandlers: proposalHandlers,
	}
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name gets basic module name
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis sets default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis validates module genesis
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, bz json.RawMessage) error {
	var data GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers rest routes
func (a AppModuleBasic) RegisterRESTRoutes(ctx sdkclient.Context, rtr *mux.Router) {
	proposalRESTHandlers := make([]rest.ProposalRESTHandler, len(a.proposalHandlers))
	for i, proposalHandler := range a.proposalHandlers {
		proposalRESTHandlers[i] = proposalHandler.RESTHandler(ctx)
	}

	rest.RegisterRoutes(ctx, rtr, proposalRESTHandlers)
}

// GetTxCmd gets the root tx command of this module
func (a AppModuleBasic) GetTxCmd(ctx sdkclient.Context) *cobra.Command {
	proposalCLIHandlers := make([]*cobra.Command, len(a.proposalHandlers))
	for i, proposalHandler := range a.proposalHandlers {
		proposalCLIHandlers[i] = proposalHandler.CLIHandler(ctx)
	}

	return sdkgovcli.NewTxCmd(ctx, proposalCLIHandlers)
}

// GetQueryCmd gets the root query command of this module
func (AppModuleBasic) GetQueryCmd(ctx sdkclient.Context) *cobra.Command {
	return sdkgovcli.GetQueryCmd(StoreKey, ctx.Codec)
}

// AppModule defines app module
type AppModule struct {
	AppModuleBasic
	keeper     Keeper
	accKeeper  sdkGovTypes.AccountKeeper
	bankKeeper sdkGovTypes.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(_ version.ProtocolVersionType, keeper Keeper, accKeeper sdkGovTypes.AccountKeeper,
	bankKeeper sdkGovTypes.BankKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		accKeeper:      accKeeper,
		bankKeeper:     bankKeeper,
	}
}

// Name gets module name
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, am.keeper)
}

// Route gets module message route name
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

// NewHandler gets module handler
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute module querier route name
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler gets module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper.Keeper)
}

func (am AppModule) RegisterQueryService(grpc.Server) {}

// InitGenesis inits module genesis
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.accKeeper, am.bankKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports module genesis
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock implements module begin-block
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock implements module end-block
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}
