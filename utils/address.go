package utils

import (
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Account struct {
	Address string
	Bytes   []byte
}

func TestAccount() Account {
	bytes := secp256k1.GenPrivKey().PubKey().Address().Bytes()
	address, _ := sdk.Bech32ifyAddressBytes("noble", bytes)

	return Account{
		Address: address,
		Bytes:   bytes,
	}
}
