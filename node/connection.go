package node

import (
	"net"
	"net/rpc"
	"time"
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

func StartConnection(localAddr string, serverAddr string) *Node {
	Log.Println("Initiating Network Connection")

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
	Log.Printf("Sending HeartBeat.. [Rate : %d]\n",HeartBeatRate)
	reply := false
	for {
		rpcConn.Call("Server.HeartBeat", nodeInfo, &reply)
		time.Sleep(HeartBeatRate * time.Millisecond)
	}
}

const (
	ROUTE string = "ROUTE" 
	GETFILE string = "GETFILE"
	END string = "END"
)

type Operation struct {
	Op   string
	Next []byte
}

type Route struct {
	Dst  string
	Next []byte
}

type Data struct {
	PublicKey string
	FInfo FileInfo
	Data []byte
}

type Message struct {
	Op   Operation
	Dst  Route
	Data []byte
}

func (n *Node) Incoming(arg []byte, reply *bool) error {
	Log.Println("RPC - Received RPC Message...")
	var msg Message
	err := decryptStruct(arg,n.privateKey,&msg)
	if err != nil {
		Log.Println(err)
	}

	Log.Printf("RPC - Message Decrypted [OP : %s]\n",msg.Op.Op)
	switch msg.Op.Op (
	case ROUTE :

	case GETFILE:

	case END :

	)
	return nil
}
