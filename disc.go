package main

import (
	"crypto/ecdsa"
	"log"
	"net"

	"github.com/vipnode/ethboot/forked/discv5"
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
	nodes := []string{
		// vipnode
		//"enode://19b5013d24243a659bda7f1df13933bb05820ab6c3ebf6b5e0854848b97e1f7e308f703466e72486c5bc7fe8ed402eb62f6303418e05d330a5df80738ac974f6@163.172.138.100:30303?discport=30301",
		// local
		//"enode://4ec357a8409303a0fcced83ed2751ff14ed17c3764c617bf98bbef3f048bbd3e03d6732d234ab7c1c5426a22efa6934997c1f8979482d2360ffeef9fc7cc2c94@127.0.0.1:30303?discport=30301",
		// infura server
		"enode://5177285d3cfa92945c1b515e476717d65a809799cd6138a3331154f8607b7073d851704c9f3837ee23db853baf5df746b7af6541ac53fce7523f97dbf278ca8c@18.207.134.211:30303",
	}
	remoteNodes := []*discv5.Node{}
	for _, remoteURL := range nodes {
		remoteNode, err := discv5.ParseNode(remoteURL)
		if err != nil {
			exit(5, "failed to parse remote node URL: %s", err)
		}
		remoteNodes = append(remoteNodes, remoteNode)
	}

	log.Printf("adding remote nodes: %s", remoteNodes)
	network.SetFallbackNodes(remoteNodes)
	return network, err
}

/*
func startNodiscover(nodeKey *ecdsa.PrivateKey, conn *net.UDPConn) error {
	config := nodiscover.Config{
		PrivateKey: nodeKey,
	}
	if _, err := nodiscover.ListenUDP(conn, config); err != nil {
		return err
	}
	return nil
}
*/
