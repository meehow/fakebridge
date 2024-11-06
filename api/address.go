package api

import (
	"fakebridge/encoder"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/meehow/go-ptr"
	"github.com/meehow/go-trezor/pb"
	"google.golang.org/protobuf/proto"
)

func GetAddress(key *hdkeychain.ExtendedKey, msg *encoder.Message) ([]byte, error) {
	req := new(pb.GetAddress)
	err := proto.Unmarshal(msg.Data, req)
	if err != nil {
		return nil, err
	}
	address, err := key.Address(netParams(req.GetCoinName()))
	if err != nil {
		return nil, err
	}
	// TODO: Mac
	return encoder.Encode(&pb.Address{
		Address: ptr.String(address.String()),
	})
}
