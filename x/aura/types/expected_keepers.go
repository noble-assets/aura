package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
)

type BankKeeper interface {
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type TransferKeeper interface {
	Transfer(ctx context.Context, msg *transfertypes.MsgTransfer) (*transfertypes.MsgTransferResponse, error)
}
