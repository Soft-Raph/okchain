package distribution

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gogo/protobuf/grpc"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/okex/okchain/x/distribution/client/cli"
	"github.com/okex/okchain/x/distribution/client/rest"
	"github.com/okex/okchain/x/distribution/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic is a struct of app module basics object
type AppModuleBasic struct{}

// Name returns module name
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers module codec
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis gives a validity check to module genesis
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, bz json.RawMessage) error {
	var data GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// RegisterRESTRoutes registers rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, StoreKey)
}

// GetTxCmd gets the root tx command of this module
func (AppModuleBasic) GetTxCmd(ctx client.Context) *cobra.Command {
	return cli.GetTxCmd(StoreKey, ctx.Codec)
}

// GetQueryCmd gets the root query command of this module
func (AppModuleBasic) GetQueryCmd(ctx client.Context) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, ctx.Codec)
}

// AppModule is a struct of app module
type AppModule struct {
	AppModuleBasic
	keeper       Keeper
	supplyKeeper types.BankKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper Keeper, supplyKeeper types.BankKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
		supplyKeeper:   supplyKeeper,
	}
}

// Name returns module name
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants registers invariants
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	RegisterInvariants(ir, am.keeper)
}

func (am AppModule) RegisterQueryService(grpc.Server) {}

// Route returns the message routing key for the distribution module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, NewHandler(am.keeper))
}

// NewHandler returns module handler
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

// QuerierRoute returns module querier route name
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}

// NewQuerierHandler returns module querier
func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

// InitGenesis initializes module genesis
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, am.supplyKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports module genesis
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock is invoked on the beginning of each block
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.keeper)
}

// EndBlock is invoked on the end of each block
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
