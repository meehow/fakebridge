package api

import (
	_ "embed"
	"encoding/json"
	"fakebridge/encoder"
	"fmt"

	"github.com/meehow/go-trezor/pb"
)

var (
	//go:embed features.json
	featuresJSON []byte
)

func GetFeatures() ([]byte, error) {
	features := new(pb.Features)
	err := json.Unmarshal(featuresJSON, features)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}
	return encoder.Encode(features)
}
