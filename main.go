package main

import (
	"fakebridge/api"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
)

const defaultMnemonic = "all all all all all all all all all all all all"

func main() {
	mnemonic := os.Getenv("MNEMONIC")
	password := os.Getenv("PASSWORD")
	if mnemonic == "" {
		mnemonic = defaultMnemonic
		fmt.Printf("Using default mnemonic %q. Define custom one in env variable MNEMONIC\n", defaultMnemonic)
	}
	testnet := false
	index := 0
	port := 21325
	flag.BoolVar(&testnet, "testnet", testnet, "generate testnet addresses")
	flag.IntVar(&index, "index", index, "last element of derivation path")
	flag.IntVar(&port, "port", port, "tcp/ip port to listen on")
	flag.Parse()
	derivationPath := []uint32{44 + 0x80000000, 0 + 0x80000000, 0 + 0x80000000, 0}
	net := &chaincfg.MainNetParams
	if testnet {
		net = &chaincfg.TestNet3Params
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		fmt.Println("Invalid BIP39 mnemonic, falling back to Electrum mnemonic")
		seed, derivationPath, err = electrumSeed(mnemonic, password)
		if err != nil {
			log.Fatal(err)
		}
	}
	master, err := hdkeychain.NewMaster(seed, net)
	if err != nil {
		log.Fatal(err)
	}
	key := master
	derivationPath = append(derivationPath, uint32(index))
	for _, p := range derivationPath {
		key, err = key.Derive(p)
		if err != nil {
			log.Fatal(err)
		}
	}
	address, err := key.Address(net)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Address: %s\n", address)
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	mux := api.NewMux(logger, key, net)
	err = http.ListenAndServe(fmt.Sprintf("localhost:%d", port), mux)
	if err != nil {
		log.Fatal(err)
	}
}
