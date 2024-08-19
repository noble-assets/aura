package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorstypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ondoprotocol/usdy-noble/x/aura/types/blocklist"
)

var _ blocklist.QueryServer = &blocklistQueryServer{}

type blocklistQueryServer struct {
	*Keeper
}

func NewBlocklistQueryServer(keeper *Keeper) blocklist.QueryServer {
	return &blocklistQueryServer{Keeper: keeper}
}

func (k blocklistQueryServer) Owner(goCtx context.Context, req *blocklist.QueryOwner) (*blocklist.QueryOwnerResponse, error) {
	if req == nil {
		return nil, errorstypes.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	blocklistOwner, err := k.GetBlocklistOwner(ctx)
	if err != nil {
		return nil, err
	}
	blocklistPendingOwner, err := k.GetBlocklistPendingOwner(ctx)
	return &blocklist.QueryOwnerResponse{
		Owner:        blocklistOwner,
		PendingOwner: blocklistPendingOwner,
	}, err
}

func (k blocklistQueryServer) Addresses(goCtx context.Context, req *blocklist.QueryAddresses) (*blocklist.QueryAddressesResponse, error) {
	if req == nil {
		return nil, errorstypes.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	blocklistStore := prefix.NewStore(store, blocklist.BlockedAddressPrefix)
	var addresses []string
	pagination, err := query.Paginate(blocklistStore, req.Pagination, func(key []byte, _ []byte) error {
		addresses = append(addresses, sdk.AccAddress(key).String())
		return nil
	})

	return &blocklist.QueryAddressesResponse{
		Addresses:  addresses,
		Pagination: pagination,
	}, err
}

func (k blocklistQueryServer) Address(goCtx context.Context, req *blocklist.QueryAddress) (*blocklist.QueryAddressResponse, error) {
	if req == nil {
		return nil, errorstypes.ErrInvalidRequest
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to decode address %s", req.Address)
	}

	blocked, err := k.HasBlockedAddress(ctx, address)
	return &blocklist.QueryAddressResponse{Blocked: blocked}, err
}
