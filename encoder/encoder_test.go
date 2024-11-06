package encoder_test

import (
	"encoding/hex"
	"encoding/json"
	"fakebridge/encoder"
	"testing"

	"github.com/meehow/go-trezor/pb"
	"google.golang.org/protobuf/proto"
)

const messageHex = "0011000000c60a097472657a6f722e696f100218082000321866613465623030303030303030303030303030303030303038014a05656e2d5553520a66616b6562726964676560016a14dd4671a5104952ef505d28d1f9e94d1484b4607a800100aa0106536166652033ca01065472657a6f72f00101f00102f00103f00104f00105f00107f00109f0010bf0010cf0010df0010ef0010ff00110f00111f00112f00113800200c80200d00203e2020454324231e80201f00200f8028001800340880301900301980300b00300"

var messageBin, _ = hex.DecodeString(messageHex)

func TestDecode(t *testing.T) {
	msg, err := encoder.Decode(messageBin)
	if err != nil {
		t.Fatal(err)
	}
	features := new(pb.Features)
	err = proto.Unmarshal(msg.Data, features)
	if err != nil {
		t.Fatal(err)
	}
	featuresJson, err := json.MarshalIndent(features, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(featuresJson))

}
