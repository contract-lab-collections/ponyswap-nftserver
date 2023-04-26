package bctools

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func Eip712V4Sign(privKey string, domainSeparator, structHash string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return "", err
	}

	// msg := fmt.Sprintf("\u0019\u0001%s%s", domainSeparator, structHash)
	msg := fmt.Sprintf("\x19\x01%s%s", domainSeparator, structHash)
	dataHash := crypto.Keccak256Hash([]byte(msg))
	sig, err := crypto.Sign(dataHash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}
