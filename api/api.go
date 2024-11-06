package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fakebridge/encoder"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/meehow/go-ptr"
	"github.com/meehow/go-trezor/pb"
)

var (
	sessionChan = make(chan string, 1)
	sessionID   = ""
)

const (
	timeout = 30
	version = "2.0.34"
	githash = "unknown"
)

type api struct {
	logger *log.Logger
	key    *hdkeychain.ExtendedKey
	cancel context.CancelFunc
}

type Device struct {
	Path         string  `json:"path"`
	Vendor       int     `json:"vendor"`
	Product      int     `json:"product"`
	Debug        bool    `json:"debug"` // has debug enabled?
	Session      *string `json:"session"`
	DebugSession *string `json:"debugSession"`
}

func (a *api) Info(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]string{
		"version": version,
		"githash": githash,
	})
	a.checkJSONError(w, err)
}

func (a *api) Listen(w http.ResponseWriter, r *http.Request) {
	select {
	case sessionID = <-sessionChan:
	case <-time.After(time.Second * timeout):
	}
	a.Enumerate(w, r)
}

func (a *api) Enumerate(w http.ResponseWriter, r *http.Request) {
	var session *string
	if sessionID != "" {
		session = &sessionID
	}
	err := json.NewEncoder(w).Encode([]Device{{
		Path:    "1",
		Vendor:  4617,
		Product: 21441,
		Session: session,
	}})
	a.checkJSONError(w, err)
}

func (a *api) Acquire(w http.ResponseWriter, r *http.Request) {
	session := strconv.Itoa(rand.Int())
	err := json.NewEncoder(w).Encode(map[string]string{
		"session": session,
	})
	a.checkJSONError(w, err)
	sessionChan <- session
}

func (a *api) Release(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	err := json.NewEncoder(w).Encode("session")
	a.checkJSONError(w, err)
	sessionChan <- ""
}

func (a *api) Call(w http.ResponseWriter, r *http.Request) {
	binbody, err := io.ReadAll(hex.NewDecoder(r.Body))
	if err != nil {
		a.respondError(w, err)
		return
	}
	defer r.Body.Close()
	resp, err := a.call(binbody)
	if err != nil {
		a.logger.Println(err)
		resp, err = encoder.Encode(&pb.Failure{
			Code:    pb.Failure_Failure_UnexpectedMessage.Enum(),
			Message: ptr.String(err.Error()),
		})
		if err != nil {
			a.respondError(w, err)
			return
		}
	}
	_, err = w.Write([]byte(hex.EncodeToString(resp)))
	if err != nil {
		a.respondError(w, err)
	}
}

func (a *api) call(binbody []byte) ([]byte, error) {
	msg, err := encoder.Decode(binbody)
	if err != nil {
		return nil, err
	}
	log.Printf("---> %s %x", msg.Kind, binbody)
	switch msg.Kind {
	case pb.MessageType_MessageType_Initialize:
		return GetFeatures()
	case pb.MessageType_MessageType_GetFeatures:
		return GetFeatures()
	case pb.MessageType_MessageType_GetAddress:
		return GetAddress(a.key, msg)
	case pb.MessageType_MessageType_SignMessage:
		return SignMessage(a.key, msg)
	case pb.MessageType_MessageType_Cancel:
		if a.cancel != nil {
			a.cancel()
		}
		return nil, nil
	default:
		return nil, fmt.Errorf("%q not implemented", msg.Kind)
	}
}

func (a *api) Read(w http.ResponseWriter, r *http.Request) {}

func (a *api) checkJSONError(w http.ResponseWriter, err error) {
	if err != nil {
		a.respondError(w, err)
	}
}

func (a *api) respondError(w http.ResponseWriter, err error) {
	a.logger.Print("Returning error: " + err.Error())
	w.WriteHeader(http.StatusBadRequest)
	err = json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
	if err != nil {
		a.logger.Print("Error while writing error: " + err.Error())
	}
}

func netParams(coinName string) *chaincfg.Params {
	if coinName == "Bitcoin" {
		return &chaincfg.MainNetParams
	}
	return &chaincfg.TestNet3Params
}
