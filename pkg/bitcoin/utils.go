package bitcoin

import "github.com/btcsuite/btcd/wire"

func IsWitnessAddress(vin *wire.TxIn) bool {
	return len(vin.SignatureScript) == 0 && len(vin.Witness) > 1
}
func IsSegWitAddress(vin *wire.TxIn) bool {
	return len(vin.SignatureScript) != 0 && len(vin.Witness) > 1
}
func IsScriptSigAddress(vin *wire.TxIn) bool {
	return len(vin.SignatureScript) != 0 && len(vin.Witness) == 0
}
func IsTaprootAddress(vin *wire.TxIn) bool {
	return len(vin.SignatureScript) != 0 && len(vin.Witness) == 1
}
