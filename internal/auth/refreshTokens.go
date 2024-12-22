package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	randomDat := make([]byte, 32)
	_, err := rand.Read(randomDat)
	if err != nil {
		return "", fmt.Errorf("Something went wrong w/ generating refresh token")
	}

	hexRandomDat := hex.EncodeToString(randomDat)
	return hexRandomDat, nil
}
