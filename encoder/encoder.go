package encoder

import (
	"encoding/binary"
	"errors"

	"github.com/meehow/go-trezor/pb"
	"google.golang.org/protobuf/proto"
)

var ErrMalformedData = errors.New("malformed data")

type Message struct {
	Kind pb.MessageType
	Data []byte
}

func Encode(msg proto.Message) ([]byte, error) {
	name := "MessageType_" + string(msg.ProtoReflect().Descriptor().Name())
	kind := uint16(pb.MessageType_value[name])
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	var header [6]byte
	binary.BigEndian.PutUint16(header[0:2], kind)
	binary.BigEndian.PutUint32(header[2:6], uint32(len(data)))
	return append(header[:], data...), nil
}

func Decode(binbody []byte) (*Message, error) {
	if len(binbody) < 6 {
		return nil, ErrMalformedData
	}
	kind := binary.BigEndian.Uint16(binbody[0:2])
	size := binary.BigEndian.Uint32(binbody[2:6])
	data := binbody[6:]
	if uint32(len(data)) != size {
		return nil, ErrMalformedData
	}
	return &Message{
		Kind: pb.MessageType(kind),
		Data: data,
	}, nil
}
