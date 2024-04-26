package keeper_test

import (
	"testing"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/noble-assets/ondo/utils"
	"github.com/noble-assets/ondo/utils/mocks"
	"github.com/noble-assets/ondo/x/usdy/keeper"
	"github.com/noble-assets/ondo/x/usdy/types/blocklist"
	"github.com/stretchr/testify/require"
)

func TestOwnerQuery(t *testing.T) {
	k, ctx := mocks.USDYKeeper(t)
	server := keeper.NewBlocklistQueryServer(k)

	// ACT: Attempt to query blocklist owner with invalid request.
	_, err := server.Owner(ctx, nil)
	// ASSERT: The query should've failed due to invalid request.
	require.ErrorContains(t, err, errors.ErrInvalidRequest.Error())

	// ACT: Attempt to query blocklist owner with no state.
	_, err = server.Owner(ctx, &blocklist.QueryOwner{})
	// ASSERT: The query should've failed.
	require.ErrorContains(t, err, collections.ErrNotFound.Error())

	// ARRANGE: Set blocklist owner in state.
	owner := utils.TestAccount()
	require.NoError(t, k.Owner.Set(ctx, owner.Address))

	// ACT: Attempt to query blocklist owner with state.
	res, err := server.Owner(ctx, &blocklist.QueryOwner{})
	// ASSERT: The query should've succeeded, with empty pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Empty(t, res.PendingOwner)

	// ARRANGE: Set blocklist pending owner in state.
	pendingOwner := utils.TestAccount()
	require.NoError(t, k.PendingOwner.Set(ctx, pendingOwner.Address))

	// ACT: Attempt to query blocklist owner with state.
	res, err = server.Owner(ctx, &blocklist.QueryOwner{})
	// ASSERT: The query should've succeeded, with pending owner.
	require.NoError(t, err)
	require.Equal(t, owner.Address, res.Owner)
	require.Equal(t, pendingOwner.Address, res.PendingOwner)
}