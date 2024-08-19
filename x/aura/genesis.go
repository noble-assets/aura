package aura

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/keeper"
	"github.com/ondoprotocol/usdy-noble/x/aura/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/types/blocklist"
)

func InitGenesis(ctx sdk.Context, k *keeper.Keeper, genesis types.GenesisState) {
	k.SetBlocklistOwner(ctx, genesis.BlocklistState.Owner)
	k.SetBlocklistPendingOwner(ctx, genesis.BlocklistState.PendingOwner)
	for _, account := range genesis.BlocklistState.BlockedAddresses {
		address, _ := sdk.AccAddressFromBech32(account)
		k.SetBlockedAddress(ctx, address)
	}

	k.SetPaused(ctx, genesis.Paused)
	k.SetOwner(ctx, genesis.Owner)
	k.SetPendingOwner(ctx, genesis.PendingOwner)
	for _, burner := range genesis.Burners {
		k.SetBurner(ctx, burner.Address, burner.Allowance)
	}
	for _, minter := range genesis.Minters {
		k.SetMinter(ctx, minter.Address, minter.Allowance)
	}
	for _, pauser := range genesis.Pausers {
		k.SetPauser(ctx, pauser)
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
