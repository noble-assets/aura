package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/types"
	"github.com/ondoprotocol/usdy-noble/x/aura/types/blocklist"
)

type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService

	Denom      string
	bankKeeper types.BankKeeper

	schema collections.Schema

	paused          collections.Item[bool]
	owner           collections.Item[[]byte]
	pendingOwner    collections.Item[[]byte]
	burner          collections.Map[[]byte, []byte]
	minter          collections.Map[[]byte, []byte]
	pauser          collections.KeySet[[]byte]
	blockedChannels collections.KeySet[[]byte]

	// Blocklist
	blocklistOwner        collections.Item[string]
	blocklistPendingOwner collections.Item[string]
	blockedAddresses      collections.KeySet[[]byte]
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	denom string,
	bankKeeper types.BankKeeper,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)
	paused := collections.NewItem(
		sb,
		types.PausedKey,
		"paused",
		collections.BoolValue,
	)
	owner := collections.NewItem(
		sb,
		types.OwnerKey,
		"owner",
		collections.BytesValue,
	)
	pendingOwner := collections.NewItem(
		sb,
		types.PendingOwnerKey,
		"pendingOwner",
		collections.BytesValue,
	)
	burner := collections.NewMap(
		sb,
		types.BurnerPrefix,
		"burner",
		collections.BytesKey,
		collections.BytesValue,
	)
	minter := collections.NewMap(
		sb,
		types.MinterPrefix,
		"minter",
		collections.BytesKey,
		collections.BytesValue,
	)
	pauser := collections.NewKeySet(
		sb,
		types.PauserPrefix,
		"pauser",
		collections.BytesKey,
	)
	blockedChannels := collections.NewKeySet(
		sb,
		types.BlockedChannelPrefix,
		"blockedChannels",
		collections.BytesKey,
	)

	blocklistOwner := collections.NewItem(
		sb,
		blocklist.OwnerKey,
		"blocklistOwner",
		collections.StringValue,
	)
	blocklistPendingOwner := collections.NewItem(
		sb,
		blocklist.PendingOwnerKey,
		"blocklistPendingOwner",
		collections.StringValue,
	)
	blockedAddresses := collections.NewKeySet(
		sb,
		blocklist.BlockedAddressPrefix,
		"blockedAddresses",
		collections.BytesKey,
	)
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:          cdc,
		storeService: storeService,

		Denom:      denom,
		bankKeeper: bankKeeper,

		schema:          schema,
		paused:          paused,
		owner:           owner,
		pendingOwner:    pendingOwner,
		burner:          burner,
		minter:          minter,
		pauser:          pauser,
		blockedChannels: blockedChannels,
		// Blocklist collections
		blocklistOwner:        blocklistOwner,
		blocklistPendingOwner: blocklistPendingOwner,
		blockedAddresses:      blockedAddresses,
	}
}

// SetBankKeeper overwrites the bank keeper used in this module.
func (k *Keeper) SetBankKeeper(bankKeeper types.BankKeeper) {
	k.bankKeeper = bankKeeper
}

// SendRestrictionFn executes necessary checks against all USDY transfers.
func (k *Keeper) SendRestrictionFn(goCtx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) (newToAddr sdk.AccAddress, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if amount := amt.AmountOf(k.Denom); !amount.IsZero() {
		burning := !fromAddr.Equals(types.ModuleAddress) && toAddr.Equals(types.ModuleAddress)
		if burning {
			return toAddr, nil
		}

		paused, err := k.paused.Get(ctx)
		if err != nil {
			return toAddr, err
		}
		if paused {
			return toAddr, fmt.Errorf("%s transfers are paused", k.Denom)
		}

		minting := fromAddr.Equals(types.ModuleAddress) && !toAddr.Equals(types.ModuleAddress)

		if !minting {
			hasBlockedAddress, err := k.HasBlockedAddress(ctx, fromAddr)
			if err != nil {
				return toAddr, err
			}
			if hasBlockedAddress {
				return toAddr, fmt.Errorf("%s is blocked from sending %s", fromAddr, k.Denom)
			}
		}

		hasBlockedAddress, err := k.HasBlockedAddress(ctx, toAddr)
		if err != nil {
			return toAddr, err
		}
		if hasBlockedAddress {
			return toAddr, fmt.Errorf("%s is blocked from receiving %s", toAddr, k.Denom)
		}

		blockedChannels, err := k.GetBlockedChannels(ctx)
		if err != nil {
			return toAddr, err
		}
		for _, channel := range blockedChannels {
			escrow := transfertypes.GetEscrowAddress(transfertypes.PortID, channel)

			if toAddr.Equals(escrow) {
				return toAddr, fmt.Errorf("%s transfers are blocked on %s", k.Denom, channel)
			}
		}
	}

	return toAddr, nil
}
