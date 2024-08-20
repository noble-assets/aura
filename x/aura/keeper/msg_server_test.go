package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ondoprotocol/usdy-noble/utils"
	"github.com/ondoprotocol/usdy-noble/utils/mocks"
	"github.com/ondoprotocol/usdy-noble/x/aura/keeper"
	"github.com/ondoprotocol/usdy-noble/x/aura/types"
	"github.com/stretchr/testify/require"
)

var ONE = math.NewInt(1_000_000_000_000_000_000)

func TestBurn(t *testing.T) {
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.AuraKeeperWithBank(t, bank)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set burner in state, with enough allowance for a single burn.
	burner := utils.TestAccount()
	k.SetBurner(ctx, burner.Address, ONE)

	// ACT: Attempt to burn with invalid signer.
	_, err := server.Burn(goCtx, &types.MsgBurn{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidBurner.Error())

	// ACT: Attempt to burn with invalid account address.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account.
	user := utils.TestAccount()

	// ACT: Attempt to burn invalid amount.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to invalid amount.
	require.ErrorContains(t, err, "amount must be positive")

	// ACT: Attempt to burn from user with insufficient funds.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to insufficient funds.
	require.ErrorContains(t, err, "unable to transfer from user to module")

	// ARRANGE: Give user 1 $USDY.
	bank.Balances[user.Address] = sdk.NewCoins(sdk.NewCoin(k.Denom, ONE))

	// ACT: Attempt to burn.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.True(t, bank.Balances[user.Address].IsZero())
	require.True(t, bank.Balances[types.ModuleName].IsZero())
	getburner, err := k.GetBurner(ctx, burner.Address)
	require.NoError(t, err)
	require.True(t, getburner.IsZero())

	// ACT: Attempt another burn with insufficient allowance.
	_, err = server.Burn(goCtx, &types.MsgBurn{
		Signer: burner.Address,
		From:   user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to insufficient allowance.
	require.ErrorContains(t, err, types.ErrInsufficientAllowance.Error())
}

func TestMint(t *testing.T) {
	bank := mocks.BankKeeper{
		Balances:    make(map[string]sdk.Coins),
		Restriction: mocks.NoOpSendRestrictionFn,
	}
	k, ctx := mocks.AuraKeeperWithBank(t, bank)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set minter in state, with enough allowance for a single mint.
	minter := utils.TestAccount()
	err := k.SetMinter(ctx, minter.Address, ONE)
	require.NoError(t, err)

	// ACT: Attempt to mint with invalid signer.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidMinter.Error())

	// ACT: Attempt to mint with invalid account address.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to invalid account address.
	require.ErrorContains(t, err, "unable to decode account address")

	// ARRANGE: Generate a user account and add to blocklist.
	user := utils.TestAccount()
	err = k.SetBlockedAddress(ctx, user.Bytes)
	require.NoError(t, err)

	// ACT: Attempt to mint to blocked address.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to blocked address.
	require.ErrorContains(t, err, "blocked from receiving")

	// ARRANGE: Unblock user account.
	k.DeleteBlockedAddress(ctx, user.Bytes)

	// ACT: Attempt to mint invalid amount.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to invalid amount.
	require.ErrorContains(t, err, "amount must be positive")

	// ACT: Attempt to mint.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	require.Equal(t, ONE, bank.Balances[user.Address].AmountOf(k.Denom))
	require.True(t, bank.Balances[types.ModuleName].IsZero())
	getminter, err := k.GetMinter(ctx, minter.Address)
	require.NoError(t, err)
	require.True(t, getminter.IsZero())

	// ACT: Attempt another mint with insufficient allowance.
	_, err = server.Mint(goCtx, &types.MsgMint{
		Signer: minter.Address,
		To:     user.Address,
		Amount: ONE,
	})
	// ASSERT: The action should've failed due to insufficient allowance.
	require.ErrorContains(t, err, types.ErrInsufficientAllowance.Error())
}

func TestPause(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set pauser in state.
	pauser := utils.TestAccount()
	k.SetPauser(ctx, pauser.Address)

	// ACT: Attempt to pause with invalid signer.
	_, err := server.Pause(goCtx, &types.MsgPause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPauser.Error())
	paused, err := k.GetPaused(ctx)
	require.Error(t, err)
	require.False(t, paused)

	// ACT: Attempt to pause.
	_, err = server.Pause(goCtx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	paused, err = k.GetPaused(ctx)
	require.NoError(t, err)
	require.True(t, paused)

	// ACT: Attempt to pause again.
	_, err = server.Pause(goCtx, &types.MsgPause{
		Signer: pauser.Address,
	})
	// ASSERT: The action should've failed due to module being paused already.
	require.ErrorContains(t, err, "module is already paused")
	paused, err = k.GetPaused(ctx)
	require.NoError(t, err)
	require.True(t, paused)
}

func TestUnpause(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ARRANGE: Set paused state to true.
	err := k.SetPaused(ctx, true)
	require.NoError(t, err)

	// ACT: Attempt to unpause with no owner set.
	_, err = server.Unpause(goCtx, &types.MsgUnpause{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")
	paused, err := k.GetPaused(ctx)
	require.NoError(t, err)
	require.True(t, paused)

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	err = k.SetOwner(ctx, owner.Address)
	require.NoError(t, err)

	// ACT: Attempt to unpause with invalid signer.
	_, err = server.Unpause(goCtx, &types.MsgUnpause{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())
	paused, err = k.GetPaused(ctx)
	require.NoError(t, err)
	require.True(t, paused)

	// ACT: Attempt to unpause.
	_, err = server.Unpause(goCtx, &types.MsgUnpause{
		Signer: owner.Address,
	})
	// ASSERT: The action should've succeeded.
	require.NoError(t, err)
	paused, err = k.GetPaused(ctx)
	require.NoError(t, err)
	require.False(t, paused)

	// ACT: Attempt to unpause again.
	_, err = server.Unpause(goCtx, &types.MsgUnpause{
		Signer: owner.Address,
	})
	// ASSERT: The action should've failed due to module being unpaused already.
	require.ErrorContains(t, err, "module is already unpaused")
	paused, err = k.GetPaused(ctx)
	require.NoError(t, err)
	require.False(t, paused)
}

func TestTransferOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to transfer ownership with no owner set.
	_, err := server.TransferOwnership(goCtx, &types.MsgTransferOwnership{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to transfer ownership with invalid signer.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ACT: Attempt to transfer ownership to same owner.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: owner.Address,
	})
	// ASSERT: The action should've failed due to same owner.
	require.ErrorContains(t, err, types.ErrSameOwner.Error())

	// ARRANGE: Generate a pending owner account.
	pendingOwner := utils.TestAccount()

	// ACT: Attempt to transfer ownership.
	_, err = server.TransferOwnership(goCtx, &types.MsgTransferOwnership{
		Signer:   owner.Address,
		NewOwner: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and set a pending owner in state.
	require.NoError(t, err)
	getPendingOwner, err := k.GetPendingOwner(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, getPendingOwner)
}

func TestAcceptOwnership(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to accept ownership with no pending owner set.
	_, err := server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{})
	// ASSERT: The action should've failed due to no pending owner set.
	require.ErrorContains(t, err, "there is no pending owner")

	// ARRANGE: Set pending owner in state.
	pendingOwner := utils.TestAccount()
	k.SetPendingOwner(ctx, pendingOwner.Address)

	// ACT: Attempt to accept ownership with invalid signer.
	_, err = server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidPendingOwner.Error())

	// ACT: Attempt to accept ownership.
	_, err = server.AcceptOwnership(goCtx, &types.MsgAcceptOwnership{
		Signer: pendingOwner.Address,
	})
	// ASSERT: The action should've succeeded, and updated the owner in state.
	require.NoError(t, err)
	getowner, err := k.GetOwner(ctx)
	require.NoError(t, err)
	require.Equal(t, pendingOwner.Address, getowner)
	getpendingowner, err := k.GetPendingOwner(ctx)
	require.NoError(t, err)
	require.Empty(t, getpendingowner)
}

func TestAddBurner(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add burner with no owner set.
	_, err := server.AddBurner(goCtx, &types.MsgAddBurner{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to add burner with invalid signer.
	_, err = server.AddBurner(goCtx, &types.MsgAddBurner{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two burner accounts, add one to state.
	burner1, burner2 := utils.TestAccount(), utils.TestAccount()
	k.SetBurner(ctx, burner2.Address, ONE)

	// ACT: Attempt to add burner that already exists.
	_, err = server.AddBurner(goCtx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner2.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to existing burner.
	require.ErrorContains(t, err, "is already a burner")

	// ACT: Attempt to add burner with invalid allowance.
	_, err = server.AddBurner(goCtx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner1.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ACT: Attempt to add burner.
	_, err = server.AddBurner(goCtx, &types.MsgAddBurner{
		Signer:    owner.Address,
		Burner:    burner1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set burner in state.
	require.NoError(t, err)
	getburner, err := k.GetBurner(ctx, burner1.Address)
	require.NoError(t, err)
	require.Equal(t, ONE, getburner)
}

func TestRemoveBurner(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove burner with no owner set.
	_, err := server.RemoveBurner(goCtx, &types.MsgRemoveBurner{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to remove burner with invalid signer.
	_, err = server.RemoveBurner(goCtx, &types.MsgRemoveBurner{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a burner account.
	burner := utils.TestAccount()

	// ACT: Attempt to remove burner that does not exist.
	_, err = server.RemoveBurner(goCtx, &types.MsgRemoveBurner{
		Signer: owner.Address,
		Burner: burner.Address,
	})
	// ASSERT: The action should've failed due to non existent burner.
	require.ErrorContains(t, err, "is not a burner")

	// ARRANGE: Set burner in state.
	k.SetBurner(ctx, burner.Address, ONE)

	// ACT: Attempt to remove burner.
	_, err = server.RemoveBurner(goCtx, &types.MsgRemoveBurner{
		Signer: owner.Address,
		Burner: burner.Address,
	})
	// ASSERT: The action should've succeeded, and removed burner in state.
	require.NoError(t, err)
	hasburner, err := k.HasBurner(ctx, burner.Address)
	require.NoError(t, err)
	require.False(t, hasburner)
}

func TestSetBurnerAllowance(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set burner allowance with no owner set.
	_, err := server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to set burner allowance with invalid signer.
	_, err = server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a burner account.
	burner := utils.TestAccount()

	// ACT: Attempt to set burner allowance that does not exist.
	_, err = server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to non existent burner.
	require.ErrorContains(t, err, "is not a burner")

	// ARRANGE: Set burner in state.
	k.SetBurner(ctx, burner.Address, math.ZeroInt())

	// ACT: Attempt to set burner allowance with invalid allowance.
	_, err = server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ACT: Attempt to set burner allowance.
	_, err = server.SetBurnerAllowance(goCtx, &types.MsgSetBurnerAllowance{
		Signer:    owner.Address,
		Burner:    burner.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set burner allowance in state.
	require.NoError(t, err)
	getburner, err := k.GetBurner(ctx, burner.Address)
	require.NoError(t, err)
	require.Equal(t, ONE, getburner)
}

func TestAddMinter(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add minter with no owner set.
	_, err := server.AddMinter(goCtx, &types.MsgAddMinter{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to add minter with invalid signer.
	_, err = server.AddMinter(goCtx, &types.MsgAddMinter{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two minter accounts, add one to state.
	minter1, minter2 := utils.TestAccount(), utils.TestAccount()
	k.SetMinter(ctx, minter2.Address, ONE)

	// ACT: Attempt to add minter that already exists.
	_, err = server.AddMinter(goCtx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter2.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to existing minter.
	require.ErrorContains(t, err, "is already a minter")

	// ACT: Attempt to add minter with invalid allowance.
	_, err = server.AddMinter(goCtx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter1.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ACT: Attempt to add minter.
	_, err = server.AddMinter(goCtx, &types.MsgAddMinter{
		Signer:    owner.Address,
		Minter:    minter1.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set minter in state.
	require.NoError(t, err)
	getminter, err := k.GetMinter(ctx, minter1.Address)
	require.NoError(t, err)
	require.Equal(t, ONE, getminter)
}

func TestRemoveMinter(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove minter with no owner set.
	_, err := server.RemoveMinter(goCtx, &types.MsgRemoveMinter{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to remove minter with invalid signer.
	_, err = server.RemoveMinter(goCtx, &types.MsgRemoveMinter{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a minter account.
	minter := utils.TestAccount()

	// ACT: Attempt to remove minter that does not exist.
	_, err = server.RemoveMinter(goCtx, &types.MsgRemoveMinter{
		Signer: owner.Address,
		Minter: minter.Address,
	})
	// ASSERT: The action should've failed due to non-existent minter.
	require.ErrorContains(t, err, "is not a minter")

	// ARRANGE: Set minter in state.
	k.SetMinter(ctx, minter.Address, ONE)

	// ACT: Attempt to remove minter.
	_, err = server.RemoveMinter(goCtx, &types.MsgRemoveMinter{
		Signer: owner.Address,
		Minter: minter.Address,
	})
	// ASSERT: The action should've succeeded, and removed minter in state.
	require.NoError(t, err)
	hasminter, err := k.HasMinter(ctx, minter.Address)
	require.NoError(t, err)
	require.False(t, hasminter)
}

func TestSetMinterAllowance(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to set minter allowance with no owner set.
	_, err := server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to set minter allowance with invalid signer.
	_, err = server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a minter account.
	minter := utils.TestAccount()

	// ACT: Attempt to set minter allowance that does not exist.
	_, err = server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've failed due to non-existent minter.
	require.ErrorContains(t, err, "is not a minter")

	// ARRANGE: Set minters in state.
	k.SetMinter(ctx, minter.Address, math.ZeroInt())

	// ACT: Attempt to set minter allowance with invalid allowance.
	_, err = server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE.Neg(),
	})
	// ASSERT: The action should've failed due to negative allowance.
	require.ErrorContains(t, err, "allowance cannot be negative")

	// ACT: Attempt to set minter allowance.
	_, err = server.SetMinterAllowance(goCtx, &types.MsgSetMinterAllowance{
		Signer:    owner.Address,
		Minter:    minter.Address,
		Allowance: ONE,
	})
	// ASSERT: The action should've succeeded, and set minter allowance in state.
	require.NoError(t, err)
	getminter, err := k.GetMinter(ctx, minter.Address)
	require.NoError(t, err)
	require.Equal(t, ONE, getminter)
}

func TestAddPauser(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add pauser with no owner set.
	_, err := server.AddPauser(goCtx, &types.MsgAddPauser{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to add pauser with invalid signer.
	_, err = server.AddPauser(goCtx, &types.MsgAddPauser{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two pauser accounts, add one to state.
	pauser1, pauser2 := utils.TestAccount(), utils.TestAccount()
	k.SetPauser(ctx, pauser2.Address)

	// ACT: Attempt to add pauser that already exists.
	_, err = server.AddPauser(goCtx, &types.MsgAddPauser{
		Signer: owner.Address,
		Pauser: pauser2.Address,
	})
	// ASSERT: The action should've failed due to existing pauser.
	require.ErrorContains(t, err, "is already a pauser")

	// ACT: Attempt to add pauser.
	_, err = server.AddPauser(goCtx, &types.MsgAddPauser{
		Signer: owner.Address,
		Pauser: pauser1.Address,
	})
	// ASSERT: The action should've succeeded, and set pauser in state.
	require.NoError(t, err)
	haspauser, err := k.HasPauser(ctx, pauser1.Address)
	require.NoError(t, err)
	require.True(t, haspauser)
}

func TestRemovePauser(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove pauser with no owner set.
	_, err := server.RemovePauser(goCtx, &types.MsgRemovePauser{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to remove pauser with invalid signer.
	_, err = server.RemovePauser(goCtx, &types.MsgRemovePauser{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a pauser account.
	pauser := utils.TestAccount()

	// ACT: Attempt to remove pauser that does not exist.
	_, err = server.RemovePauser(goCtx, &types.MsgRemovePauser{
		Signer: owner.Address,
		Pauser: pauser.Address,
	})
	// ASSERT: The action should've failed due to non-existent pauser.
	require.ErrorContains(t, err, "is not a pauser")

	// ARRANGE: Set pauser in state.
	k.SetPauser(ctx, pauser.Address)

	// ACT: Attempt to remove pauser.
	_, err = server.RemovePauser(goCtx, &types.MsgRemovePauser{
		Signer: owner.Address,
		Pauser: pauser.Address,
	})
	// ASSERT: The action should've succeeded, and removed pauser in state.
	require.NoError(t, err)
	haspauser, err := k.HasPauser(ctx, pauser.Address)
	require.NoError(t, err)
	require.False(t, haspauser)
}

func TestAddBlockedChannel(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to add blocked channel with no owner set.
	_, err := server.AddBlockedChannel(goCtx, &types.MsgAddBlockedChannel{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to add blocked channel with invalid signer.
	_, err = server.AddBlockedChannel(goCtx, &types.MsgAddBlockedChannel{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate two channel, add one to state.
	channel1, channel2 := "channel-0", "channel-1"
	k.SetBlockedChannel(ctx, channel2)

	// ACT: Attempt to add blocked channel that is blocked.
	_, err = server.AddBlockedChannel(goCtx, &types.MsgAddBlockedChannel{
		Signer:  owner.Address,
		Channel: channel2,
	})
	// ASSERT: The action should've failed due to blocked channel.
	require.ErrorContains(t, err, "is already blocked")

	// ACT: Attempt to add blocked channel.
	_, err = server.AddBlockedChannel(goCtx, &types.MsgAddBlockedChannel{
		Signer:  owner.Address,
		Channel: channel1,
	})
	// ASSERT: The action should've succeeded, and set channel in state.
	require.NoError(t, err)
	hasblockedchannel, err := k.HasBlockedChannel(ctx, channel1)
	require.NoError(t, err)
	require.True(t, hasblockedchannel)
}

func TestRemoveBlockedChannel(t *testing.T) {
	k, ctx := mocks.AuraKeeper(t)
	goCtx := sdk.WrapSDKContext(ctx)
	server := keeper.NewMsgServer(k)

	// ACT: Attempt to remove blocked channel with no owner set.
	_, err := server.RemoveBlockedChannel(goCtx, &types.MsgRemoveBlockedChannel{})
	// ASSERT: The action should've failed due to no owner set.
	require.ErrorContains(t, err, "there is no owner")

	// ARRANGE: Set owner in state.
	owner := utils.TestAccount()
	k.SetOwner(ctx, owner.Address)

	// ACT: Attempt to remove blocked channel with invalid signer.
	_, err = server.RemoveBlockedChannel(goCtx, &types.MsgRemoveBlockedChannel{
		Signer: utils.TestAccount().Address,
	})
	// ASSERT: The action should've failed due to invalid signer.
	require.ErrorContains(t, err, types.ErrInvalidOwner.Error())

	// ARRANGE: Generate a channel.
	channel := "channel-0"

	// ACT: Attempt to remove blocked channel that isn't blocked.
	_, err = server.RemoveBlockedChannel(goCtx, &types.MsgRemoveBlockedChannel{
		Signer:  owner.Address,
		Channel: channel,
	})
	// ASSERT: The action should've failed due to allowed channel.
	require.ErrorContains(t, err, "is not blocked")

	// ARRANGE: Set channel in state.
	k.SetBlockedChannel(ctx, channel)

	// ACT: Attempt to remove blocked channel.
	_, err = server.RemoveBlockedChannel(goCtx, &types.MsgRemoveBlockedChannel{
		Signer:  owner.Address,
		Channel: channel,
	})
	// ASSERT: The action should've succeeded, and removed channel in state.
	require.NoError(t, err)
	hasblockedchannel, err := k.HasBlockedChannel(ctx, channel)
	require.NoError(t, err)
	require.False(t, hasblockedchannel)
}
