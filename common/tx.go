package common

import (
	"context"
	"fmt"
	"github.com/AstraProtocol/astra-go-sdk/account"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"

	keyMultisig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authSigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/pkg/errors"
	emvTypes "github.com/tharsis/ethermint/x/evm/types"
)

type Tx struct {
	txf        tx.Factory
	privateKey *account.PrivateKeySerialized
}

func NewTx(config client.TxConfig, privateKey *account.PrivateKeySerialized, gasLimit uint64, gasPrice string) *Tx {
	txf := tx.Factory{}.
		WithTxConfig(config).
		WithGasPrices(gasPrice).
		WithGas(gasLimit)
	//.SetTimeoutHeight(txf.TimeoutHeight())

	return &Tx{txf: txf, privateKey: privateKey}
}

func (t *Tx) BuildUnsignedTx(msgs types.Msg) (client.TxBuilder, error) {
	return t.txf.BuildUnsignedTx(msgs)
}

func (t *Tx) PrintUnsignedTx(rpcClient client.Context, msgs types.Msg) (string, error) {
	unsignedTx, err := t.BuildUnsignedTx(msgs)
	if err != nil {
		return "", errors.Wrap(err, "BuildUnsignedTx")
	}

	json, err := rpcClient.TxConfig.TxJSONEncoder()(unsignedTx.GetTx())
	if err != nil {
		return "", errors.Wrap(err, "TxJSONEncoder")
	}

	return string(json), nil
}

func (t *Tx) prepareSignTx(rpcClient client.Context, coinType uint32, from types.AccAddress) error {

	if err := rpcClient.AccountRetriever.EnsureExists(rpcClient, from); err != nil {
		return errors.Wrap(err, "EnsureExists")
	}

	initNum, initSeq := t.txf.AccountNumber(), t.txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		var accNum, accSeq uint64
		var err error

		if coinType == 60 {
			queryClient := emvTypes.NewQueryClient(rpcClient)
			cosmosAccount, err := queryClient.CosmosAccount(context.Background(), &emvTypes.QueryCosmosAccountRequest{Address: from.String()})
			if err != nil {
				return errors.Wrap(err, "CosmosAccount")
			}

			accNum = cosmosAccount.AccountNumber
			accSeq = cosmosAccount.Sequence

		} else {
			accNum, accSeq, err = rpcClient.AccountRetriever.GetAccountNumberSequence(rpcClient, from)
			if err != nil {
				return errors.Wrap(err, "GetAccountNumberSequence")
			}
		}

		t.txf.WithAccountNumber(accNum)
		t.txf.WithSequence(accSeq)
	}

	return nil
}

func (t *Tx) SignTx(rpcClient client.Context, txBuilder client.TxBuilder, pubKey cryptoTypes.PubKey) error {
	from := t.privateKey.Address()

	err := t.prepareSignTx(rpcClient, t.privateKey.CoinType(), from)
	if err != nil {
		return errors.Wrap(err, "prepareSignTx")
	}

	signMode := rpcClient.TxConfig.SignModeHandler().DefaultMode()

	isMulSign, err := IsMulSign(pubKey)
	if isMulSign {
		signMode = signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}

	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: t.txf.Sequence(),
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return errors.Wrap(err, "SetSignatures")
	}

	// Construct the SignatureV2 struct
	signerData := authSigning.SignerData{
		ChainID:       rpcClient.ChainID,
		AccountNumber: t.txf.AccountNumber(),
		Sequence:      t.txf.Sequence(),
	}

	signWithPrivKey, err := tx.SignWithPrivKey(
		signMode,
		signerData,
		txBuilder,
		t.privateKey.PrivateKey(),
		rpcClient.TxConfig,
		t.txf.Sequence())

	if err != nil {
		return errors.Wrap(err, "SignWithPrivKey")
	}

	err = txBuilder.SetSignatures(signWithPrivKey)
	if err != nil {
		return errors.Wrap(err, "SetSignatures")
	}

	return nil
}

func (t *Tx) MulSignTx(rpcClient client.Context, txBuilder client.TxBuilder, pubKey cryptoTypes.PubKey, coinType uint32, sigs []signing.SignatureV2) error {
	from := types.AccAddress(pubKey.Address())
	err := t.prepareSignTx(rpcClient, coinType, from)
	if err != nil {
		return errors.Wrap(err, "prepareSignTx")
	}

	multisigPub, ok := pubKey.(*keyMultisig.LegacyAminoPubKey)
	if !ok {
		return errors.Wrap(errors.New("set type error"), "LegacyAminoPubKey")
	}

	multisigSig := multisig.NewMultisig(len(multisigPub.PubKeys))

	for i := 0; i < len(sigs); i++ {
		signingData := authSigning.SignerData{
			ChainID:       t.txf.ChainID(),
			AccountNumber: t.txf.AccountNumber(),
			Sequence:      t.txf.Sequence(),
		}

		for _, sig := range sigs {
			err = authSigning.VerifySignature(sig.PubKey, signingData, sig.Data, rpcClient.TxConfig.SignModeHandler(), txBuilder.GetTx())
			if err != nil {
				addr, _ := types.AccAddressFromHex(sig.PubKey.Address().String())
				return fmt.Errorf("couldn't verify signature for address %s", addr)
			}

			if err := multisig.AddSignatureV2(multisigSig, sig, multisigPub.GetPubKeys()); err != nil {
				return errors.Wrap(err, "AddSignatureV2")
			}
		}
	}

	sigV2 := signing.SignatureV2{
		PubKey:   multisigPub,
		Data:     multisigSig,
		Sequence: t.txf.Sequence(),
	}

	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return errors.Wrap(err, "SetSignatures")
	}

	return nil
}
