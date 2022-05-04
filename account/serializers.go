package account

import (
	"encoding/hex"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	ethermintTypes "github.com/tharsis/ethermint/types"
)

type PrivateKeySerialized struct {
	mnemonic   string
	privateKey cryptoTypes.PrivKey
}

func NewPrivateKeySerialized(mnemonic string, privateKey cryptoTypes.PrivKey) *PrivateKeySerialized {
	return &PrivateKeySerialized{mnemonic: mnemonic, privateKey: privateKey}
}

func (p *PrivateKeySerialized) String() (string, error) {
	pub := p.privateKey.PubKey()

	addr := types.AccAddress(pub.Address())
	validatorAddr := types.ValAddress(pub.Address())
	hexAddr := common.BytesToAddress(pub.Address().Bytes())

	apk, err := codecTypes.NewAnyWithValue(pub)
	if err != nil {
		return "", errors.Wrap(err, "NewKeyOutput")
	}
	bz, err := codec.ProtoMarshalJSON(apk, nil)

	rs := map[string]string{
		"privateKey":   hex.EncodeToString(p.privateKey.Bytes()),
		"mnemonic":     p.mnemonic,
		"publicKey":    string(bz),
		"validatorKey": validatorAddr.String(),
		"address":      addr.String(),
		"hexAddress":   hexAddr.String(),
		"type":         p.privateKey.Type(),
	}

	b, _ := json.MarshalIndent(rs, "", " ")

	return string(b), nil
}

func (p *PrivateKeySerialized) PrivateKey() cryptoTypes.PrivKey {
	return p.privateKey
}

func (p *PrivateKeySerialized) PublicKey() cryptoTypes.PubKey {
	return p.privateKey.PubKey()
}

func (p *PrivateKeySerialized) PublicKeyJson() (string, error) {
	pub := p.privateKey.PubKey()
	apk, err := codecTypes.NewAnyWithValue(pub)
	if err != nil {
		return "", errors.Wrap(err, "NewAnyWithValue")
	}

	bz, err := codec.ProtoMarshalJSON(apk, nil)
	if err != nil {
		return "", errors.Wrap(err, "ProtoMarshalJSON")
	}

	return string(bz), nil
}

func (p *PrivateKeySerialized) AccAddress() types.AccAddress {
	pub := p.privateKey.PubKey()
	addr := types.AccAddress(pub.Address())

	return addr
}

func (p *PrivateKeySerialized) ValidatorAddress() types.ValAddress {
	pub := p.privateKey.PubKey()
	addr := types.ValAddress(pub.Address())

	return addr
}

func (p *PrivateKeySerialized) HexAddress() common.Address {
	pub := p.privateKey.PubKey()
	addr := common.BytesToAddress(pub.Address().Bytes())

	return addr
}

func (p *PrivateKeySerialized) Type() string {
	return p.privateKey.Type()
}

func (p *PrivateKeySerialized) CoinType() uint32 {
	if p.privateKey.Type() == "secp256k1" || p.privateKey.Type() == "ed25519" {
		return types.CoinType
	}

	return ethermintTypes.Bip44CoinType
}
