package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ondoprotocol/usdy-noble/x/aura/types"
)

var _ types.QueryServer = &queryServer{}

type queryServer struct {
	*Keeper
}

func NewQueryServer(keeper *Keeper) types.QueryServer {
	return &queryServer{Keeper: keeper}
}

func (k queryServer) Denom(_ context.Context, req *types.QueryDenom) (*types.QueryDenomResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	return &types.QueryDenomResponse{Denom: k.Keeper.Denom}, nil
}

func (k queryServer) Paused(goCtx context.Context, req *types.QueryPaused) (*types.QueryPausedResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Ignore if the paused state is not set, in that case it is considered unpaused
	paused, _ := k.GetPaused(ctx)
	return &types.QueryPausedResponse{Paused: paused}, nil
}

func (k queryServer) Owner(goCtx context.Context, req *types.QueryOwner) (*types.QueryOwnerResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	owner, err := k.GetOwner(ctx)
	if err != nil {
		return nil, err
	}
	// Ignore the error in case pending owner is not available
	pendingOwner, _ := k.GetPendingOwner(ctx)
	return &types.QueryOwnerResponse{
		Owner:        owner,
		PendingOwner: pendingOwner,
	}, nil
}

func (k queryServer) Burners(goCtx context.Context, req *types.QueryBurners) (*types.QueryBurnersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	burners, err := k.GetBurners(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryBurnersResponse{Burners: burners}, nil
}

func (k queryServer) Minters(goCtx context.Context, req *types.QueryMinters) (*types.QueryMintersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	minters, err := k.GetMinters(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryMintersResponse{Minters: minters}, nil
}

func (k queryServer) Pausers(goCtx context.Context, req *types.QueryPausers) (*types.QueryPausersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pausers, err := k.GetPausers(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryPausersResponse{Pausers: pausers}, nil
}

func (k queryServer) BlockedChannels(goCtx context.Context, req *types.QueryBlockedChannels) (*types.QueryBlockedChannelsResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	blockedChannels, err := k.GetBlockedChannels(ctx)
	if err != nil {
		return nil, err
	}
	return &types.QueryBlockedChannelsResponse{BlockedChannels: blockedChannels}, nil
}
