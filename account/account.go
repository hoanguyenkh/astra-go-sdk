package account

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/pkg/errors"
	ethermintHd "github.com/tharsis/ethermint/crypto/hd"
	ethermintTypes "github.com/tharsis/ethermint/types"
)

type Account struct {
	coinType uint32
}

func NewAccount(coinType uint32) *Account {
	return &Account{coinType: coinType}
}

//Create new an Account

func (a *Account) CreateAccount() (*PrivateKeySerialized, error) {
	mnemonicEntropySize := 256
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return nil, errors.Wrap(err, "NewEntropy")
	}

	//bip39 mnemonic
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return nil, errors.Wrap(err, "NewMnemonic")
	}

	privKey, err := a.ImportAccount(mnemonic)
	if err != nil {
		return nil, errors.Wrap(err, "importAccount")
	}

	return privKey, nil
}

func (a *Account) CreateMulSignAccount(totalSign, multisigThreshold int) ([]*PrivateKeySerialized, string, string, error) {
	var listPrivate []*PrivateKeySerialized
	pks := make([]cryptoTypes.PubKey, totalSign)
	for i := 0; i < totalSign; i++ {
		k, err := a.CreateAccount()
		if err != nil {
			continue
		}

		listPrivate = append(listPrivate, k)
		pks[i] = k.PublicKey()
	}

	pk := multisig.NewLegacyAminoPubKey(multisigThreshold, pks)

	addr := types.AccAddress(pk.Address())
	apk, err := codecTypes.NewAnyWithValue(pk)
	if err != nil {
		return nil, "", "", errors.Wrap(err, "NewKeyOutput")
	}
	bz, err := codec.ProtoMarshalJSON(apk, nil)

	return listPrivate, addr.String(), string(bz), nil
}

//Import an Account

func (a *Account) ImportAccount(mnemonic string) (*PrivateKeySerialized, error) {
	if a.coinType == 60 {
		derivedPriv, err := ethermintHd.EthSecp256k1.Derive()(
			mnemonic,
			keyring.DefaultBIP39Passphrase,
			ethermintTypes.BIP44HDPath)

		if err != nil {
			return nil, errors.Wrap(err, "Derive")
		}

		privateKey := ethermintHd.EthSecp256k1.Generate()(derivedPriv)
		return NewPrivateKeySerialized(mnemonic, privateKey), nil

	}

	//cosmos
	derivedPriv, err := hd.Secp256k1.Derive()(
		mnemonic,
		keyring.DefaultBIP39Passphrase,
		types.FullFundraiserPath,
	)

	if err != nil {
		return nil, errors.Wrap(err, "Derive")
	}

	privateKey := hd.Secp256k1.Generate()(derivedPriv)
	return NewPrivateKeySerialized(mnemonic, privateKey), nil
}
