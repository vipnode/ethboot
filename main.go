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
	"github.com/ethereum/go-ethereum/p2p/discv5"
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

	network, err := discv5.ListenUDP(
		nodeKey,
		conn,
		conn.LocalAddr().(*net.UDPAddr),
		"<no database>", // TODO: Do we want a database?
		nil,
	)
	if err != nil {
		exit(4, "failed to start discovery server: %s", err)
	}
	log.Printf("listening on enode://%s@%s", discover.PubkeyID(&nodeKey.PublicKey), conn.LocalAddr())

	// How to start a discv4 server, which does not support replacing the
	// bootnodes yet so we're using only v5 for now. Also, light clients require v5?
	/*
		config := discover.Config{
			PrivateKey: nodeKey,
		}
		if _, err := discover.ListenUDP(conn, config); err != nil {
			exit(4, "failed to start discovery server: %s", err)
		}
	*/

	// Must be run with discv5?
	remoteURL := "enode://19b5013d24243a659bda7f1df13933bb05820ab6c3ebf6b5e0854848b97e1f7e308f703466e72486c5bc7fe8ed402eb62f6303418e05d330a5df80738ac974f6@163.172.138.100:30303?discport=30301"
	// Local override
	remoteURL = "enode://4ec357a8409303a0fcced83ed2751ff14ed17c3764c617bf98bbef3f048bbd3e03d6732d234ab7c1c5426a22efa6934997c1f8979482d2360ffeef9fc7cc2c94@127.0.0.1:30303"
	remoteNode, err := discv5.ParseNode(remoteURL)
	if err != nil {
		exit(5, "failed to parse remote node URL: %s", err)
	}

	log.Printf("adding remote node: %s", remoteNode)
	network.SetFallbackNodes([]*discv5.Node{remoteNode})

	// Block?
	select {}
}

// exit prints an error and exits with the given code
func exit(code int, msg string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", a...)
	os.Exit(code)
}
