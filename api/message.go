package api

import (
	"bytes"
	"fakebridge/encoder"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/meehow/go-ptr"
	"github.com/meehow/go-trezor/pb"
	"google.golang.org/protobuf/proto"
)

const MessageSignatureHeader = "Bitcoin Signed Message:\n"

func SignMessage(key *hdkeychain.ExtendedKey, msg *encoder.Message) ([]byte, error) {
	req := new(pb.SignMessage)
	err := proto.Unmarshal(msg.Data, req)
	if err != nil {
		return nil, err
	}
	address, err := key.Address(netParams(req.GetCoinName()))
	if err != nil {
		return nil, err
	}
	priv, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	wire.WriteVarString(&buf, 0, MessageSignatureHeader)
	wire.WriteVarBytes(&buf, 0, req.Message)
	messageHash := chainhash.DoubleHashB(buf.Bytes())
	signature, err := ecdsa.SignCompact(priv, messageHash, true)
	if err != nil {
		return nil, err
	}
	return encoder.Encode(&pb.MessageSignature{
		Address:   ptr.String(address.String()),
		Signature: signature,
	})
}
