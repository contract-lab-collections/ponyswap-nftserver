package bctools

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

/*
	msg => rawMsg => msgHash => signature
	fmt.Sprintf("xxx%s", msg) => rawMsg
	crypto.Keccak256Hash([]byte(rawMsg)) => msgHash
	crypto.Sign(msgHash.Bytes(), privateKey) => signature

	data => hexData => dataHash => signature
	hexutil.Decode(hexdata) => dataHash
	crypto.Sign(dataHash, privateKey) => signature
*/

const EthSignStr = "Ethereum Signed Message:"

func EthSign(privKey string, msg string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", err
	}

	rawMsg := fmt.Sprintf("\x19%s\n%d%s", EthSignStr, len(msg), msg)
	msgHash := crypto.Keccak256Hash([]byte(rawMsg))
	sig, err := crypto.Sign(msgHash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}

func EthVerifySignAddress(address string, msg string, sign string) bool {
	sign = strings.TrimPrefix(sign, "0x")
	if len(sign) != 130 {
		return false
	}

	rawMsg := fmt.Sprintf("\u0019%s\n%d%s", EthSignStr, len(msg), msg)
	dataHash := crypto.Keccak256Hash([]byte(rawMsg))

	signature, err := hex.DecodeString(sign)
	if err != nil {
		return false
	}

	if signature[64] >= 27 {
		signature[64] -= 27
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(dataHash.Bytes(), signature)
	if err != nil {
		return false
	}

	addr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	return strings.EqualFold(addr.Hex(), address)
}

func EthVerifySignature(pubKey string, rawMsg string, sign string) bool {
	dataHash := crypto.Keccak256Hash([]byte(rawMsg))

	sign = strings.TrimPrefix(sign, "0x")
	signature, err := hex.DecodeString(sign)
	if err != nil {
		return false
	}

	pubkey, err := hex.DecodeString(pubKey)
	if err != nil {
		return false
	}
	return crypto.VerifySignature(pubkey, dataHash[:], signature[:len(signature)-1])
}

/////////////////////////////////////////////////////////////////////////////////////

// The signature content is already hashed data
func EthSignHexBytes(privKey string, hexdata string) ([]byte, error) {
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	dataHash, err := hexutil.Decode(hexdata)
	if err != nil {
		return nil, err
	}

	signature, err := crypto.Sign(dataHash, privateKey)
	if err != nil {
		return nil, err
	}

	if signature[64] < 27 {
		signature[64] += 27
	}

	return signature, err
}

func VerifySignatureByHexdata(account string, hexdata string, sign string) bool {
	sign = strings.TrimPrefix(sign, "0x")
	if len(sign) != 130 {
		return false
	}

	dataHash, err := hex.DecodeString(strings.TrimPrefix(hexdata, "0x"))
	if err != nil {
		return false
	}

	signature, err := hex.DecodeString(sign)
	if err != nil {
		return false
	}

	if signature[64] >= 27 {
		signature[64] -= 27
	}

	sigPublicKeyECDSA, err := crypto.SigToPub(dataHash, signature)
	if err != nil {
		return false
	}

	addr := crypto.PubkeyToAddress(*sigPublicKeyECDSA)
	return strings.EqualFold(addr.Hex(), account)
}
