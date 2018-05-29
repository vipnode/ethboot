package main

import (
	"crypto/ecdsa"
	"log"
	"net"

	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/vipnode/ethboot/nodiscover"
)

func startDiscv4() {
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

}

func startDiscv5(nodeKey *ecdsa.PrivateKey, conn *net.UDPConn) (*discv5.Network, error) {
	network, err := discv5.ListenUDP(
		nodeKey,
		conn,
		conn.LocalAddr().(*net.UDPAddr),
		"<no database>", // TODO: Do we want a database?
		nil,
	)
	if err != nil {
		return nil, err
	}

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
	return network, err
}

func startNodiscover(nodeKey *ecdsa.PrivateKey, conn *net.UDPConn) error {
	config := nodiscover.Config{
		PrivateKey: nodeKey,
	}
	if _, err := nodiscover.ListenUDP(conn, config); err != nil {
		return err
	}
	return nil
}
