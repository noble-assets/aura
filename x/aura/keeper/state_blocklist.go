package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//

func (k *Keeper) GetBlocklistOwner(ctx sdk.Context) (string, error) {
	return k.blocklistOwner.Get(ctx)
}

func (k *Keeper) SetBlocklistOwner(ctx sdk.Context, owner string) error {
	return k.blocklistOwner.Set(ctx, owner)
}

//

func (k *Keeper) DeleteBlocklistPendingOwner(ctx sdk.Context) error {
	return k.blocklistPendingOwner.Remove(ctx)
}

func (k *Keeper) GetBlocklistPendingOwner(ctx sdk.Context) (string, error) {
	return k.blocklistPendingOwner.Get(ctx)
}

func (k *Keeper) SetBlocklistPendingOwner(ctx sdk.Context, pendingOwner string) error {
	return k.blocklistPendingOwner.Set(ctx, pendingOwner)
}

//

func (k *Keeper) DeleteBlockedAddress(ctx sdk.Context, address []byte) error {
	return k.blockedAddresses.Remove(ctx, address)
}

func (k *Keeper) GetBlockedAddresses(ctx sdk.Context) (addresses []string, err error) {
	err = k.blockedAddresses.Walk(ctx, nil, func(address []byte) (bool, error) {
		addresses = append(addresses, sdk.AccAddress(address).String())
		return false, nil
	})
	return addresses, err
}

func (k *Keeper) HasBlockedAddress(ctx sdk.Context, address []byte) (bool, error) {
	return k.blockedAddresses.Has(ctx, address)
}

func (k *Keeper) SetBlockedAddress(ctx sdk.Context, address []byte) error {
	return k.blockedAddresses.Set(ctx, address)
}
