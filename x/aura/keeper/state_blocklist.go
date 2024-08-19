package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//

func (k *Keeper) GetBlocklistOwner(ctx sdk.Context) (string, error) {
	return k.BlocklistOwner.Get(ctx)
}

func (k *Keeper) SetBlocklistOwner(ctx sdk.Context, owner string) error {
	return k.BlocklistOwner.Set(ctx, owner)
}

//

func (k *Keeper) DeleteBlocklistPendingOwner(ctx sdk.Context) error {
	return k.BlocklistPendingOwner.Remove(ctx)
}

func (k *Keeper) GetBlocklistPendingOwner(ctx sdk.Context) (string, error) {
	return k.BlocklistPendingOwner.Get(ctx)
}

func (k *Keeper) SetBlocklistPendingOwner(ctx sdk.Context, pendingOwner string) error {
	return k.BlocklistPendingOwner.Set(ctx, pendingOwner)
}

//

func (k *Keeper) DeleteBlockedAddress(ctx sdk.Context, address []byte) error {
	return k.BlockedAddresses.Remove(ctx, address)
}

func (k *Keeper) GetBlockedAddresses(ctx sdk.Context) (addresses []string, err error) {
	err = k.BlockedAddresses.Walk(ctx, nil, func(address []byte) (bool, error) {
		addresses = append(addresses, sdk.AccAddress(address).String())
		return false, nil
	})
	return addresses, err
}

func (k *Keeper) HasBlockedAddress(ctx sdk.Context, address []byte) (bool, error) {
	return k.BlockedAddresses.Has(ctx, address)
}

func (k *Keeper) SetBlockedAddress(ctx sdk.Context, address []byte) error {
	return k.BlockedAddresses.Set(ctx, address)
}
