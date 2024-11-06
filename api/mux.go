package api

import (
	"log"
	"net/http"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

func NewMux(logger *log.Logger, key *hdkeychain.ExtendedKey, net *chaincfg.Params) *http.ServeMux {
	api := &api{
		logger: logger,
		key:    key,
	}
	r := http.NewServeMux()
	r.HandleFunc("POST /", cors(api.Info))
	r.HandleFunc("POST /configure", cors(api.Info))
	r.HandleFunc("POST /listen", cors(api.Listen))
	r.HandleFunc("POST /enumerate", cors(api.Enumerate))
	r.HandleFunc("POST /acquire/{path}", cors(api.Acquire))
	r.HandleFunc("POST /acquire/{path}/{session}", cors(api.Acquire))
	r.HandleFunc("POST /release/{session}", cors(api.Release))
	r.HandleFunc("POST /call/{session}", cors(api.Call))
	r.HandleFunc("POST /post/{session}", cors(api.Call))
	r.HandleFunc("POST /read/{session}", cors(api.Read))
	return r
}

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://connect.trezor.io")
		next.ServeHTTP(w, r)
	}
}
