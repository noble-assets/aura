package keeper_test

import (
	"testing"

	"github.com/ondoprotocol/usdy-noble/utils"
	"github.com/ondoprotocol/usdy-noble/utils/mocks"
	"github.com/stretchr/testify/require"
)

func TestGetBlockedAddresses(t *testing.T) {
	keeper, ctx := mocks.AuraKeeper(t)

	// ACT: Retrieve all blocked addresses with no state.
	addresses, err := keeper.GetBlockedAddresses(ctx)
	// ASSERT: No error returned.
	require.NoError(t, err)
	// ASSERT: No addresses returned.
	require.Empty(t, addresses)

	// ARRANGE: Set blocklist addresses in state.
	user1, user2 := utils.TestAccount(), utils.TestAccount()
	keeper.SetBlockedAddress(ctx, user1.Bytes)
	keeper.SetBlockedAddress(ctx, user2.Bytes)

	// ACT: Retrieve all blocked addresses.
	addresses, err = keeper.GetBlockedAddresses(ctx)
	// ASSERT: No error returned.
	require.NoError(t, err)
	// ASSERT: Addresses returned.
	require.Len(t, addresses, 2)
	require.Contains(t, addresses, user1.Address)
	require.Contains(t, addresses, user2.Address)
}
