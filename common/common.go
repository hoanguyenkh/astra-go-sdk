package common

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	signingTypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func DecodePublicKey(rpcClient client.Context, pkJSON string) (cryptoTypes.PubKey, error) {
	var pk cryptoTypes.PubKey
	err := rpcClient.Codec.UnmarshalInterfaceJSON([]byte(pkJSON), &pk)
	if err != nil {
		return nil, errors.Wrap(err, "UnmarshalInterfaceJSON`")
	}
	return pk, nil
}

func IsMulSign(pk cryptoTypes.PubKey) (bool, error) {
	lpk, ok := pk.(*multisig.LegacyAminoPubKey)
	if !ok {
		return false, nil
	}

	if lpk.Threshold > 1 {
		return true, nil
	}

	return false, nil
}

func IsTxSigner(user types.AccAddress, signers []types.AccAddress) bool {
	for _, s := range signers {
		if bytes.Equal(user.Bytes(), s.Bytes()) {
			return true
		}
	}

	return false
}

func TxBuilderJsonDecoder(txConfig client.TxConfig, txJSON string) ([]byte, error) {
	txJSONDecoder, err := txConfig.TxJSONDecoder()([]byte(txJSON))
	if err != nil {
		return nil, err
	}

	txBytes, err := txConfig.TxEncoder()(txJSONDecoder)
	if err != nil {
		panic(err)
	}

	return txBytes, nil
}

func TxBuilderJsonEncoder(txConfig client.TxConfig, tx client.TxBuilder) (string, error) {
	txJSONEncoder, err := txConfig.TxJSONEncoder()(tx.GetTx())
	if err != nil {
		return "", err
	}

	return string(txJSONEncoder), nil
}

func TxBuilderSignatureJsonEncoder(txConfig client.TxConfig, tx client.TxBuilder) (string, error) {
	sigs, err := tx.GetTx().GetSignaturesV2()
	if err != nil {
		return "", err
	}

	json, err := txConfig.MarshalSignatureJSON(sigs)

	return string(json), nil
}

func TxBuilderSignatureJsonDecoder(txConfig client.TxConfig, txJson string) ([]signingTypes.SignatureV2, error) {
	return txConfig.UnmarshalSignatureJSON([]byte(txJson))
}

func TxHash(txBytes []byte) string {
	return fmt.Sprintf("%X", tmhash.Sum(txBytes))
}

func IsAddressValid(address string) (bool, error) {
	receiver, err := types.AccAddressFromBech32(address)
	if err != nil {
		return false, err
	}

	return receiver.String() == address, nil
}

func EthAddressToCosmosAddress(ethAddress string) (string, error) {
	baseAddr, err := types.AccAddressFromHex(ethAddress)
	if err != nil {
		return "", err
	}

	return baseAddr.String(), nil
}
