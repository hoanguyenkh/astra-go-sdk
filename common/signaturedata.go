package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/cryptobyte"
	cryptobyteasn1 "golang.org/x/crypto/cryptobyte/asn1"
	"math/big"
)

func GenPrivateKeySign() (string, string, error) {
	key, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	pubkey := elliptic.MarshalCompressed(crypto.S256(), key.X, key.Y)

	privkey := make([]byte, 32)
	blob := key.D.Bytes()
	copy(privkey[32-len(blob):], blob)

	privkeyStr := hex.EncodeToString(privkey)

	pubkeyStr := base64.StdEncoding.EncodeToString(pubkey)

	return privkeyStr, pubkeyStr, nil
}

func SignatureData(privateKey string, msg string) (string, error) {
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(msg))
	sig, err := ecdsa.SignASN1(rand.Reader, privKey, hash[:])
	if err != nil {
		return "", err
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
