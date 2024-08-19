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

	paused, err := k.GetPaused(ctx)
	return &types.QueryPausedResponse{Paused: paused}, err
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
	pendingOwner, err := k.GetPendingOwner(ctx)
	return &types.QueryOwnerResponse{
		Owner:        owner,
		PendingOwner: pendingOwner,
	}, err
}

func (k queryServer) Burners(goCtx context.Context, req *types.QueryBurners) (*types.QueryBurnersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	burners, err := k.GetBurners(ctx)
	return &types.QueryBurnersResponse{Burners: burners}, err
}

func (k queryServer) Minters(goCtx context.Context, req *types.QueryMinters) (*types.QueryMintersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	minters, err := k.GetMinters(ctx)
	return &types.QueryMintersResponse{Minters: minters}, err
}

func (k queryServer) Pausers(goCtx context.Context, req *types.QueryPausers) (*types.QueryPausersResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pausers, err := k.GetPausers(ctx)
	return &types.QueryPausersResponse{Pausers: pausers}, err
}

func (k queryServer) BlockedChannels(goCtx context.Context, req *types.QueryBlockedChannels) (*types.QueryBlockedChannelsResponse, error) {
	if req == nil {
		return nil, errors.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	blockedChannels, err := k.GetBlockedChannels(ctx)
	return &types.QueryBlockedChannelsResponse{BlockedChannels: blockedChannels}, err
}
