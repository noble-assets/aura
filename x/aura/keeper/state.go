package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/types"
)

//

func (k *Keeper) GetPaused(ctx sdk.Context) (bool, error) {
	return k.paused.Get(ctx)
}

func (k *Keeper) SetPaused(ctx sdk.Context, paused bool) error {
	return k.paused.Set(ctx, paused)
}

//

func (k *Keeper) GetOwner(ctx sdk.Context) (string, error) {
	ownerBz, err := k.owner.Get(ctx)
	return string(ownerBz), err
}

func (k *Keeper) SetOwner(ctx sdk.Context, owner string) error {
	return k.owner.Set(ctx, []byte(owner))
}

//

func (k *Keeper) DeletePendingOwner(ctx sdk.Context) error {
	return k.pendingOwner.Remove(ctx)
}

func (k *Keeper) GetPendingOwner(ctx sdk.Context) (string, error) {
	ownerBz, err := k.pendingOwner.Get(ctx)
	return string(ownerBz), err
}

func (k *Keeper) SetPendingOwner(ctx sdk.Context, pendingOwner string) error {
	return k.pendingOwner.Set(ctx, []byte(pendingOwner))
}

//

func (k *Keeper) DeleteBurner(ctx sdk.Context, burner string) error {
	return k.burner.Remove(ctx, []byte(burner))
}

func (k *Keeper) GetBurner(ctx sdk.Context, burner string) (allowance math.Int, err error) {
	bz, err := k.burner.Get(ctx, []byte(burner))
	if err != nil {
		return
	}
	err = allowance.Unmarshal(bz)
	return
}

func (k *Keeper) GetBurners(ctx sdk.Context) (burners []types.Burner, err error) {
	err = k.burner.Walk(ctx, nil, func(burnerAddr []byte, allowanceBz []byte) (bool, error) {
		var allowance math.Int
		err = allowance.Unmarshal(allowanceBz)

		burners = append(burners, types.Burner{
			Address:   string(burnerAddr),
			Allowance: allowance,
		})
		return false, err
	})

	return
}

func (k *Keeper) HasBurner(ctx sdk.Context, burner string) (bool, error) {
	return k.burner.Has(ctx, []byte(burner))
}

func (k *Keeper) SetBurner(ctx sdk.Context, burner string, allowance math.Int) error {
	bz, err := allowance.Marshal()
	if err != nil {
		return err
	}
	return k.burner.Set(ctx, []byte(burner), bz)
}

//

func (k *Keeper) DeleteMinter(ctx sdk.Context, minter string) error {
	return k.minter.Remove(ctx, []byte(minter))
}

func (k *Keeper) GetMinter(ctx sdk.Context, minter string) (allowance math.Int, err error) {
	allowanceBz, err := k.minter.Get(ctx, []byte(minter))
	if err != nil {
		return
	}
	err = allowance.Unmarshal(allowanceBz)
	return
}

func (k *Keeper) GetMinters(ctx sdk.Context) (minters []types.Minter, err error) {
	err = k.minter.Walk(ctx, nil, func(minterAddr []byte, allowanceBz []byte) (bool, error) {
		var allowance math.Int
		err = allowance.Unmarshal(allowanceBz)

		minters = append(minters, types.Minter{
			Address:   string(minterAddr),
			Allowance: allowance,
		})
		return false, err
	})

	return
}

func (k *Keeper) HasMinter(ctx sdk.Context, minter string) (bool, error) {
	return k.minter.Has(ctx, []byte(minter))
}

func (k *Keeper) SetMinter(ctx sdk.Context, minter string, allowance math.Int) error {
	bz, _ := allowance.Marshal()
	return k.minter.Set(ctx, []byte(minter), bz)
}

//

func (k *Keeper) DeletePauser(ctx sdk.Context, pauser string) error {
	return k.pauser.Remove(ctx, []byte(pauser))
}

func (k *Keeper) GetPausers(ctx sdk.Context) (pausers []string, err error) {
	err = k.pauser.Walk(ctx, nil, func(pauser []byte) (bool, error) {
		pausers = append(pausers, string(pauser))
		return false, nil
	})

	return
}

func (k *Keeper) HasPauser(ctx sdk.Context, pauser string) (bool, error) {
	return k.pauser.Has(ctx, []byte(pauser))
}

func (k *Keeper) SetPauser(ctx sdk.Context, pauser string) error {
	return k.pauser.Set(ctx, []byte(pauser))
}

//

func (k *Keeper) DeleteBlockedChannel(ctx sdk.Context, channel string) error {
	return k.blockedChannels.Remove(ctx, []byte(channel))
}

func (k *Keeper) GetBlockedChannels(ctx sdk.Context) (channels []string, err error) {
	err = k.blockedChannels.Walk(ctx, nil, func(channel []byte) (bool, error) {
		channels = append(channels, string(channel))
		return false, nil
	})
	return
}

func (k *Keeper) HasBlockedChannel(ctx sdk.Context, channel string) (bool, error) {
	return k.blockedChannels.Has(ctx, []byte(channel))
}

func (k *Keeper) SetBlockedChannel(ctx sdk.Context, channel string) error {
	return k.blockedChannels.Set(ctx, []byte(channel))
}
