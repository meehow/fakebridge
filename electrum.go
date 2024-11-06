package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

func electrumSeed(mnemonic, password string) ([]byte, []uint32, error) {
	seed := pbkdf2.Key([]byte(mnemonic), []byte("electrum"+password), 2048, 64, sha512.New)
	version := seedVersion(mnemonic)
	switch version {
	case "01": // Standard
		return seed, []uint32{0}, nil
	case "100": // Segwit
		return nil, []uint32{0, 0}, nil
	case "101": // 2FA
		return nil, nil, errors.New("2FA seed not implemented")
	default:
		return nil, nil, fmt.Errorf("electrum seed version 0x%s not implemented", version)
	}
}

// https://electrum.readthedocs.io/en/latest/seedphrase.html
func seedVersion(mnemonic string) string {
	h := hmac.New(sha512.New, []byte("Seed version"))
	h.Write([]byte(mnemonic))
	sum := h.Sum(nil)
	length := sum[0]>>4 + 2
	return hex.EncodeToString(sum)[:length]
}
