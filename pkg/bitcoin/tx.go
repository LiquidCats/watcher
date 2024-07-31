package bitcoin

import (
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entity"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type TxTransformer func(vin *wire.TxIn, params *chaincfg.Params) (entity.Address, error)

func TransformWitnessIntoAddress(vin *wire.TxIn, params *chaincfg.Params) (entity.Address, error) {
	pubKeyHash := btcutil.Hash160(vin.Witness[len(vin.Witness)-1])
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, params)
	if err != nil {
		return "", err
	}

	return entity.Address(addr.EncodeAddress()), nil
}

func TransformScriptSigIntoAddress(vin *wire.TxIn, params *chaincfg.Params) (entity.Address, error) {
	pubKeyHash := extractPubKeyHashFromSignatureScript(vin)

	addr, err := btcutil.NewAddressPubKey(pubKeyHash, params)
	if err != nil {
		return "", err
	}

	return entity.Address(addr.EncodeAddress()), nil
}

func TransformSegWithIntoAddress(vin *wire.TxIn, params *chaincfg.Params) (entity.Address, error) {
	pubKeyHash := extractPubKeyHashFromSignatureScript(vin)

	addr, err := btcutil.NewAddressScriptHash(pubKeyHash, params)
	if err != nil {
		return "", err
	}

	return entity.Address(addr.EncodeAddress()), nil
}

func TransformTaprootIntoAddress(vin *wire.TxIn, params *chaincfg.Params) (entity.Address, error) {
	// Generate scriptPubKey for Taproot (P2TR)
	taprootPubKeyHash := chainhash.DoubleHashB(vin.Witness[0])
	scriptPubKey, err := txscript.NewScriptBuilder().AddOp(txscript.OP_1).AddData(taprootPubKeyHash).Script()
	addr, err := btcutil.NewAddressTaproot(chainhash.HashB(scriptPubKey), params)
	if err != nil {
		return "", err
	}

	return entity.Address(addr.EncodeAddress()), nil
}

func extractPubKeyHashFromSignatureScript(vin *wire.TxIn) []byte {
	pubKey, _ := txscript.PushedData(vin.SignatureScript)

	return pubKey[len(pubKey)-1]
}
