package cosmos

import (
	"encoding/base64"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	Bech32Prefix         = "ethm"
	Bech32PrefixAccAddr  = Bech32Prefix
	Bech32PrefixAccPub   = Bech32Prefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

func init() {
	sdk.GetConfig().SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	sdk.GetConfig().SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	sdk.GetConfig().SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
}

func pubkeyToValCons(pubkey string) (string, error) {
	pubkeyBytes, err := base64.StdEncoding.DecodeString(pubkey)
	if err != nil {
		return "", err
	}
	address := crypto.Address(tmhash.SumTruncated(pubkeyBytes))
	conAddress := sdk.ConsAddress(address)
	return conAddress.String(), nil
}
