package node

import (
	"net"
	"net/rpc"
	"time"
)

const (
	ROUTE = "itoa"
	GETFILE
	END
)

type NodeInfo struct {
	addr      string
	publicKey string
}

type Node struct {
	nodeInfo   NodeInfo
	listener   net.Listener
	rpcConn    *rpc.Client
	privateKey string
}

type Operation struct {
	Op   string
	Next []byte
}

type Route struct {
	Dst  string
	Next []byte
}

type Message struct {
	Op   Operation
	Dst  Route
	Data []byte
}

func StartConnection(localAddr string, serverAddr string) *Node {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		Log.Fatal(err)
	}
	Log.Printf("Running Node at %s\n", localAddr)

	connServer, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		Log.Fatal(err)
	}
	Log.Printf("Connected to Server at %s\n", serverAddr)

	pubKey, priKey := generateKeys()

	node := Node{
		nodeInfo:   NodeInfo{localAddr, pubKey},
		listener:   listener,
		rpcConn:    connServer,
		privateKey: priKey,
	}

	go sendHB(node.rpcConn, node.nodeInfo)

	return &node
}

func sendHB(rpcConn *rpc.Client, nodeInfo NodeInfo) {
	reply := false
	for {
		rpcConn.Call("Server.HeartBeat", nodeInfo, &reply)
		time.Sleep(600 * time.Millisecond)
	}
}

func (n *Node) Incoming(arg Message, reply *bool) error {
	return nil
}
