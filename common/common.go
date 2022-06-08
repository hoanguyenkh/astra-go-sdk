package common

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	signingTypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"golang.org/x/crypto/cryptobyte"
	cryptobyteasn1 "golang.org/x/crypto/cryptobyte/asn1"
	"math/big"
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
	tx, err := txConfig.TxJSONDecoder()([]byte(txJSON))
	if err != nil {
		return nil, err
	}

	//convert to []byte
	txBytes, err := txConfig.TxEncoder()(tx)
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
	ethAddr := ethCommon.HexToAddress(ethAddress)
	baseAddr := types.AccAddress(ethAddr.Bytes())
	return baseAddr.String(), nil
}

func CosmosAddressToEthAddress(cosmosAddress string) (string, error) {
	baseAddr, err := types.AccAddressFromBech32(cosmosAddress)
	if err != nil {
		return "", err
	}

	ethAddress := ethCommon.BytesToAddress(baseAddr.Bytes())
	return ethAddress.String(), nil
}

func IsBlocked(code uint32) bool {
	if code == CodeTypeOK {
		return true
	}

	return false
}

func GenPrivateKeySign() (string, string) {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	pubkey := elliptic.MarshalCompressed(crypto.S256(), key.X, key.Y)

	privkey := make([]byte, 32)
	blob := key.D.Bytes()
	copy(privkey[32-len(blob):], blob)

	privkeyStr := hex.EncodeToString(privkey)

	pubkeyStr := base64.StdEncoding.EncodeToString(pubkey)

	return privkeyStr, pubkeyStr
}

func SignatureData(privateKey string, msg string) (string, error) {
	privKey, err := crypto.HexToECDSA(privateKey)

	hash := sha256.Sum256([]byte(msg))
	sig, err := ecdsa.SignASN1(rand.Reader, privKey, hash[:])
	if err != nil {
		panic(err)
	}

	//signEncode := hex.EncodeToString(sig)

	signEncode := base64.StdEncoding.EncodeToString(sig)

	return signEncode, nil
}

func VerifySignature(publicKey string, signature string, msg string) (bool, error) {
	if len(publicKey) <= 0 {
		return false, errors.New("public key is empty")
	}

	if len(signature) <= 0 {
		return false, errors.New("signature is empty")
	}

	if len(msg) <= 0 {
		return false, errors.New("data is empty")
	}

	publicKeyByte, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false, errors.Wrap(err, "DecodeString")
	}

	pkKey, err := crypto.DecompressPubkey(publicKeyByte)
	if err != nil {
		return false, errors.Wrap(err, "DecompressPubkey")
	}

	signatureDecode, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, errors.Wrap(err, "DecodeString")
	}

	var (
		r, s  = &big.Int{}, &big.Int{}
		inner cryptobyte.String
	)

	input := cryptobyte.String(signatureDecode)
	if !input.ReadASN1(&inner, cryptobyteasn1.SEQUENCE) ||
		!input.Empty() ||
		!inner.ReadASN1Integer(r) ||
		!inner.ReadASN1Integer(s) ||
		!inner.Empty() {
		return false, errors.New("ReadASN1 error")
	}

	hash := sha256.Sum256([]byte(msg))
	valid := ecdsa.Verify(pkKey, hash[:], r, s)
	return valid, nil
}

func VerifyHdPath(hdPath string) (bool, error) {
	_, err := hdwallet.ParseDerivationPath(hdPath)
	if err != nil {
		return false, errors.Wrap(err, "ParseDerivationPath")
	}

	return true, nil
}
