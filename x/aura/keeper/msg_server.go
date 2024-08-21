package keeper

import (
	"context"
	"errors"
	"fmt"

	sdkerrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/types"
)

var _ types.MsgServer = &msgServer{}

type msgServer struct {
	*Keeper
}

func NewMsgServer(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	hasBurner, err := k.HasBurner(ctx, msg.Signer)
	if err != nil || !hasBurner {
		return nil, types.ErrInvalidBurner
	}
	allowance, err := k.GetBurner(ctx, msg.Signer)
	if err != nil || allowance.LT(msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientAllowance, "burner %s has an allowance of %s", msg.Signer, allowance.String())
	}

	from, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to decode account address %s", msg.From)
	}

	if !msg.Amount.IsPositive() {
		return nil, errors.New("amount must be positive")
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.Denom, msg.Amount))
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to transfer from user to module")
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to burn from module")
	}

	err = k.SetBurner(ctx, msg.Signer, allowance.Sub(msg.Amount))

	// NOTE: The bank module emits an event for us.
	return &types.MsgBurnResponse{}, err
}

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	hasMinter, err := k.HasMinter(ctx, msg.Signer)
	if err != nil || !hasMinter {
		return nil, types.ErrInvalidMinter
	}
	allowance, err := k.GetMinter(ctx, msg.Signer)
	if err != nil || allowance.LT(msg.Amount) {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientAllowance, "minter %s has an allowance of %s", msg.Signer, allowance.String())
	}

	to, err := sdk.AccAddressFromBech32(msg.To)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to decode account address %s", msg.To)
	}

	if !msg.Amount.IsPositive() {
		return nil, errors.New("amount must be positive")
	}

	coins := sdk.NewCoins(sdk.NewCoin(k.Denom, msg.Amount))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to mint to module")
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to transfer from module to user")
	}

	err = k.SetMinter(ctx, msg.Signer, allowance.Sub(msg.Amount))

	// NOTE: The bank module emits an event for us.
	return &types.MsgMintResponse{}, err
}

func (k msgServer) Pause(goCtx context.Context, msg *types.MsgPause) (*types.MsgPauseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	hasPauser, err := k.HasPauser(ctx, msg.Signer)
	if err != nil || !hasPauser {
		return nil, types.ErrInvalidPauser
	}

	// Ignore if the paused state is not set, in that case it is considered unpaused
	paused, _ := k.GetPaused(ctx)
	if paused {
		return nil, errors.New("module is already paused")
	}

	err = k.SetPaused(ctx, true)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to pause module")
	}

	return &types.MsgPauseResponse{}, ctx.EventManager().EmitTypedEvent(&types.Paused{
		Account: msg.Signer,
	})
}

func (k msgServer) Unpause(goCtx context.Context, msg *types.MsgUnpause) (*types.MsgUnpauseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	paused, err := k.GetPaused(ctx)
	if err != nil || !paused {
		return nil, errors.New("module is already unpaused")
	}

	err = k.SetPaused(ctx, false)

	return &types.MsgUnpauseResponse{}, ctx.EventManager().EmitTypedEvent(&types.Unpaused{
		Account: msg.Signer,
	})
}

func (k msgServer) TransferOwnership(goCtx context.Context, msg *types.MsgTransferOwnership) (*types.MsgTransferOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	if msg.NewOwner == owner {
		return nil, types.ErrSameOwner
	}

	err = k.SetPendingOwner(ctx, msg.NewOwner)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to set pending owner")
	}
	return &types.MsgTransferOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&types.OwnershipTransferStarted{
		PreviousOwner: owner,
		NewOwner:      msg.NewOwner,
	})
}

func (k msgServer) AcceptOwnership(goCtx context.Context, msg *types.MsgAcceptOwnership) (*types.MsgAcceptOwnershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pendingOwner, err := k.GetPendingOwner(ctx)
	if err != nil || pendingOwner == "" {
		return nil, types.ErrNoPendingOwner
	}
	if msg.Signer != pendingOwner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidPendingOwner, "expected %s, got %s", pendingOwner, msg.Signer)
	}

	// Ignore the error if the owner is not set, as this is  just used for event
	owner, _ := k.GetOwner(ctx)

	err = k.SetOwner(ctx, msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to set owner")
	}
	err = k.DeletePendingOwner(ctx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to delete pending owner")
	}

	return &types.MsgAcceptOwnershipResponse{}, ctx.EventManager().EmitTypedEvent(&types.OwnershipTransferred{
		PreviousOwner: owner,
		NewOwner:      msg.Signer,
	})
}

func (k msgServer) AddBurner(goCtx context.Context, msg *types.MsgAddBurner) (*types.MsgAddBurnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasBurner, err := k.HasBurner(ctx, msg.Burner)
	if err != nil || hasBurner {
		return nil, fmt.Errorf("%s is already a burner", msg.Burner)
	}

	if msg.Allowance.IsNegative() {
		return nil, errors.New("allowance cannot be negative")
	}

	err = k.SetBurner(ctx, msg.Burner, msg.Allowance)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to set burner")
	}

	return &types.MsgAddBurnerResponse{}, ctx.EventManager().EmitTypedEvent(&types.BurnerAdded{
		Address:   msg.Burner,
		Allowance: msg.Allowance,
	})
}

func (k msgServer) RemoveBurner(goCtx context.Context, msg *types.MsgRemoveBurner) (*types.MsgRemoveBurnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasBurner, err := k.HasBurner(ctx, msg.Burner)
	if err != nil || !hasBurner {
		return nil, fmt.Errorf("%s is not a burner", msg.Burner)
	}

	err = k.DeleteBurner(ctx, msg.Burner)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to delete burner")
	}

	return &types.MsgRemoveBurnerResponse{}, ctx.EventManager().EmitTypedEvent(&types.BurnerRemoved{
		Address: msg.Burner,
	})
}

func (k msgServer) SetBurnerAllowance(goCtx context.Context, msg *types.MsgSetBurnerAllowance) (*types.MsgSetBurnerAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasBurner, err := k.HasBurner(ctx, msg.Burner)
	if err != nil || !hasBurner {
		return nil, fmt.Errorf("%s is not a burner", msg.Burner)
	}

	if msg.Allowance.IsNegative() {
		return nil, errors.New("allowance cannot be negative")
	}

	allowance, err := k.GetBurner(ctx, msg.Burner)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to get burner %s allowance", msg.Burner)
	}
	err = k.SetBurner(ctx, msg.Burner, msg.Allowance)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to set burner")
	}

	return &types.MsgSetBurnerAllowanceResponse{}, ctx.EventManager().EmitTypedEvent(&types.BurnerUpdated{
		Address:           msg.Burner,
		PreviousAllowance: allowance,
		NewAllowance:      msg.Allowance,
	})
}

func (k msgServer) AddMinter(goCtx context.Context, msg *types.MsgAddMinter) (*types.MsgAddMinterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasMinter, err := k.HasMinter(ctx, msg.Minter)
	if err != nil || hasMinter {
		return nil, fmt.Errorf("%s is already a minter", msg.Minter)
	}

	if msg.Allowance.IsNegative() {
		return nil, errors.New("allowance cannot be negative")
	}

	err = k.SetMinter(ctx, msg.Minter, msg.Allowance)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to set minter")
	}

	return &types.MsgAddMinterResponse{}, ctx.EventManager().EmitTypedEvent(&types.MinterAdded{
		Address:   msg.Minter,
		Allowance: msg.Allowance,
	})
}

func (k msgServer) RemoveMinter(goCtx context.Context, msg *types.MsgRemoveMinter) (*types.MsgRemoveMinterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasMinter, err := k.HasMinter(ctx, msg.Minter)
	if err != nil || !hasMinter {
		return nil, fmt.Errorf("%s is not a minter", msg.Minter)
	}

	err = k.DeleteMinter(ctx, msg.Minter)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to delete minter")
	}

	return &types.MsgRemoveMinterResponse{}, ctx.EventManager().EmitTypedEvent(&types.MinterRemoved{
		Address: msg.Minter,
	})
}

func (k msgServer) SetMinterAllowance(goCtx context.Context, msg *types.MsgSetMinterAllowance) (*types.MsgSetMinterAllowanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasMinter, err := k.HasMinter(ctx, msg.Minter)
	if err != nil || !hasMinter {
		return nil, fmt.Errorf("%s is not a minter", msg.Minter)
	}

	if msg.Allowance.IsNegative() {
		return nil, errors.New("allowance cannot be negative")
	}

	allowance, err := k.GetMinter(ctx, msg.Minter)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to get minter %s allowance", msg.Minter)
	}
	err = k.SetMinter(ctx, msg.Minter, msg.Allowance)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to set minter")
	}

	return &types.MsgSetMinterAllowanceResponse{}, ctx.EventManager().EmitTypedEvent(&types.MinterUpdated{
		Address:           msg.Minter,
		PreviousAllowance: allowance,
		NewAllowance:      msg.Allowance,
	})
}

func (k msgServer) AddPauser(goCtx context.Context, msg *types.MsgAddPauser) (*types.MsgAddPauserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasPauser, err := k.HasPauser(ctx, msg.Pauser)
	if err != nil || hasPauser {
		return nil, fmt.Errorf("%s is already a pauser", msg.Pauser)
	}

	err = k.SetPauser(ctx, msg.Pauser)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to set pauser")
	}

	return &types.MsgAddPauserResponse{}, ctx.EventManager().EmitTypedEvent(&types.PauserAdded{
		Address: msg.Pauser,
	})
}

func (k msgServer) RemovePauser(goCtx context.Context, msg *types.MsgRemovePauser) (*types.MsgRemovePauserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasPauser, err := k.HasPauser(ctx, msg.Pauser)
	if err != nil || !hasPauser {
		return nil, fmt.Errorf("%s is not a pauser", msg.Pauser)
	}

	err = k.DeletePauser(ctx, msg.Pauser)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to delete pauser")
	}

	return &types.MsgRemovePauserResponse{}, ctx.EventManager().EmitTypedEvent(&types.PauserRemoved{
		Address: msg.Pauser,
	})
}

func (k msgServer) AddBlockedChannel(goCtx context.Context, msg *types.MsgAddBlockedChannel) (*types.MsgAddBlockedChannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasBlockedChannel, err := k.HasBlockedChannel(ctx, msg.Channel)
	if err != nil || hasBlockedChannel {
		return nil, fmt.Errorf("%s is already blocked", msg.Channel)
	}

	err = k.SetBlockedChannel(ctx, msg.Channel)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to set blocked channel")
	}

	return &types.MsgAddBlockedChannelResponse{}, ctx.EventManager().EmitTypedEvent(&types.BlockedChannelAdded{
		Channel: msg.Channel,
	})
}

func (k msgServer) RemoveBlockedChannel(goCtx context.Context, msg *types.MsgRemoveBlockedChannel) (*types.MsgRemoveBlockedChannelResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil || owner == "" {
		return nil, types.ErrNoOwner
	}
	if msg.Signer != owner {
		return nil, sdkerrors.Wrapf(types.ErrInvalidOwner, "expected %s, got %s", owner, msg.Signer)
	}

	hasBlockedChannel, err := k.HasBlockedChannel(ctx, msg.Channel)
	if err != nil || !hasBlockedChannel {
		return nil, fmt.Errorf("%s is not blocked", msg.Channel)
	}

	err = k.DeleteBlockedChannel(ctx, msg.Channel)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to delete blocked channel")
	}

	return &types.MsgRemoveBlockedChannelResponse{}, ctx.EventManager().EmitTypedEvent(&types.BlockedChannelRemoved{
		Channel: msg.Channel,
	})
}
