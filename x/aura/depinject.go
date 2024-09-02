package aura

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	modulev1 "github.com/ondoprotocol/usdy-noble/api/aura/module/v1"
	"github.com/ondoprotocol/usdy-noble/x/aura/keeper"
	"github.com/ondoprotocol/usdy-noble/x/aura/types"
)

var (
	_ appmodule.AppModule = AppModule{}
)

func (AppModule) IsOnePerModuleType() {}

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type AuraModuleInputs struct {
	depinject.In

	Config *modulev1.Module

	Cdc          codec.Codec
	StoreService store.KVStoreService
	BankKeeper   types.BankKeeper
}

type ModuleOutputs struct {
	depinject.Out

	Keeper *keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in AuraModuleInputs) ModuleOutputs {
	if in.Config.Denom == "" {
		panic("denom for x/aura module must be set")
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Config.Denom,
		in.BankKeeper,
	)
	m := NewAppModule(k)

	return ModuleOutputs{Keeper: k, Module: m}
}
