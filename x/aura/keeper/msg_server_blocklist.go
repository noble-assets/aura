package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/types/blocklist"
)

var _ blocklist.MsgServer = &blocklistMsgServer{}

type blocklistMsgServer struct {
	*Keeper
}

func NewBlocklistMsgServer(keeper *Keeper) blocklist.MsgServer {
	return &blocklistMsgServer{Keeper: keeper}
}

func (k blocklistMsgServer) TransferOwnership(goCtx context.Context, msg *blocklist.MsgTransferOwnership) (*blocklist.MsgTransferOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetBlocklistOwner(ctx)
	if err != nil || owner == "" {
		return nil, blocklist.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, errors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if msg.NewOwner == owner {
		return nil, blocklist.ErrSameOwner
	}

	err = k.SetBlocklistPendingOwner(ctx, msg.NewOwner)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set blocklist pending owner")
	}

	return &blocklist.MsgTransferOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&blocklist.OwnershipTransferStarted{
		PreviousOwner: owner,
		NewOwner:      msg.NewOwner,
	})
}

func (k blocklistMsgServer) AcceptOwnership(goCtx context.Context, msg *blocklist.MsgAcceptOwnership) (*blocklist.MsgAcceptOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pendingOwner, err := k.GetBlocklistPendingOwner(ctx)
	if err != nil || pendingOwner == "" {
		return nil, blocklist.ErrNoPendingOwner
	}
	if msg.Signer != pendingOwner {
		return nil, errors.Wrapf(blocklist.ErrInvalidPendingOwner, "expected %s, got %s", pendingOwner, msg.Signer)
	}

	// We dont need to check the error as this value is used just for events
	owner, _ := k.GetBlocklistOwner(ctx)

	err = k.SetBlocklistOwner(ctx, msg.Signer)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set blocklist owner")
	}
	err = k.DeleteBlocklistPendingOwner(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to delete blocklist pending owner")
	}

	return &blocklist.MsgAcceptOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&blocklist.OwnershipTransferred{
		PreviousOwner: owner,
		NewOwner:      msg.Signer,
	})
}

func (k blocklistMsgServer) AddToBlocklist(goCtx context.Context, msg *blocklist.MsgAddToBlocklist) (*blocklist.MsgAddToBlocklistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetBlocklistOwner(ctx)
	if err != nil || owner == "" {
		return nil, blocklist.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, errors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	for _, account := range msg.Accounts {
		address, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to decode account address %s", account)
		}

		err = k.SetBlockedAddress(ctx, address)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to set blocked address %s", account)
		}
	}

	return &blocklist.MsgAddToBlocklistResponse{}, ctx.EventManager().EmitTypedEvent(&blocklist.BlockedAddressesAdded{
		Accounts: msg.Accounts,
	})
}

func (k blocklistMsgServer) RemoveFromBlocklist(goCtx context.Context, msg *blocklist.MsgRemoveFromBlocklist) (*blocklist.MsgRemoveFromBlocklistResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetBlocklistOwner(ctx)
	if err != nil || owner == "" {
		return nil, blocklist.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, errors.Wrapf(blocklist.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	for _, account := range msg.Accounts {
		address, err := sdk.AccAddressFromBech32(account)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to decode account address %s", account)
		}

		err = k.DeleteBlockedAddress(ctx, address)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to delete blocked address %s", account)
		}
	}

	return &blocklist.MsgRemoveFromBlocklistResponse{}, ctx.EventManager().EmitTypedEvent(&blocklist.BlockedAddressesRemoved{
		Accounts: msg.Accounts,
	})
}
