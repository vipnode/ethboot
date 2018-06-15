package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	ethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

func main() {
	var (
		addr    = flag.String("listen", ":30301", "address to listen on")
		keyPath = flag.String("nodekey", "nodekey", "node private key path")
	)

	ethlog.Root().SetHandler(
		ethlog.LvlFilterHandler(
			ethlog.Lvl(9),
			ethlog.StreamHandler(os.Stdout, ethlog.TerminalFormat(true))))

	flag.Parse()

	var nodeKey *ecdsa.PrivateKey
	var err error

	if _, err = os.Stat(*keyPath); err == nil {
		if nodeKey, err = crypto.LoadECDSA(*keyPath); err != nil {
			exit(1, "failed to load node key: %s", err)
		}
	} else if os.IsNotExist(err) {
		log.Printf("Generating a fresh key: %s", keyPath)

		if nodeKey, err = crypto.GenerateKey(); err != nil {
			exit(2, "failed to generate node key: %s", err)
		}
		if err = crypto.SaveECDSA(*keyPath, nodeKey); err != nil {
			exit(2, "failed to save node key: %s", err)
		}
	}

	udpAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		exit(3, "failed to resolve UDP address: %s", err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		exit(3, "failed to open UDP connection: %s", err)
	}

	log.Printf("listening on enode://%s@%s", discover.PubkeyID(&nodeKey.PublicKey), conn.LocalAddr())

	if _, err := startDiscv5(nodeKey, conn); err != nil {
		exit(4, "failed to start discovery server: %s", err)
	}

	// Block?
	select {}
}

// exit prints an error and exits with the given code
func exit(code int, msg string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", a...)
	os.Exit(code)
}
