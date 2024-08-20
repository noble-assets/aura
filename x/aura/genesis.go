package aura

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/keeper"
	"github.com/ondoprotocol/usdy-noble/x/aura/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/types/blocklist"
)

func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genesis types.GenesisState) {
	if err := k.SetBlocklistOwner(ctx, genesis.BlocklistState.Owner); err != nil {
		panic(err)
	}
	if err := k.SetBlocklistPendingOwner(ctx, genesis.BlocklistState.PendingOwner); err != nil {
		panic(err)
	}
	for _, account := range genesis.BlocklistState.BlockedAddresses {
		address, _ := sdk.AccAddressFromBech32(account)
		if err := k.SetBlockedAddress(ctx, address); err != nil {
			panic(err)
		}
	}

	if err := k.SetPaused(ctx, genesis.Paused); err != nil {
		panic(err)
	}
	if err := k.SetOwner(ctx, genesis.Owner); err != nil {
		panic(err)
	}
	if err := k.SetPendingOwner(ctx, genesis.PendingOwner); err != nil {
		panic(err)
	}
	for _, burner := range genesis.Burners {
		if err := k.SetBurner(ctx, burner.Address, burner.Allowance); err != nil {
			panic(err)
		}
	}
	for _, minter := range genesis.Minters {
		if err := k.SetMinter(ctx, minter.Address, minter.Allowance); err != nil {
			panic(err)
		}
	}
	for _, pauser := range genesis.Pausers {
		if err := k.SetPauser(ctx, pauser); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	blocklistOwner, err := k.GetBlocklistOwner(ctx)
	if err != nil {
		panic(err)
	}
	blocklistPendingOwner, err := k.GetBlocklistPendingOwner(ctx)
	if err != nil {
		panic(err)
	}
	blockedAddresses, err := k.GetBlockedAddresses(ctx)
	if err != nil {
		panic(err)
	}
	paused, err := k.GetPaused(ctx)
	if err != nil {
		panic(err)
	}
	owner, err := k.GetOwner(ctx)
	if err != nil {
		panic(err)
	}
	pendingOwner, err := k.GetPendingOwner(ctx)
	if err != nil {
		panic(err)
	}
	burners, err := k.GetBurners(ctx)
	if err != nil {
		panic(err)
	}
	minters, err := k.GetMinters(ctx)
	if err != nil {
		panic(err)
	}
	pausers, err := k.GetPausers(ctx)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		BlocklistState: blocklist.GenesisState{
			Owner:            blocklistOwner,
			PendingOwner:     blocklistPendingOwner,
			BlockedAddresses: blockedAddresses,
		},
		Paused:       paused,
		Owner:        owner,
		PendingOwner: pendingOwner,
		Burners:      burners,
		Minters:      minters,
		Pausers:      pausers,
	}
}
