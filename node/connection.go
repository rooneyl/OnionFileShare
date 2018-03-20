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
	nodeInfo       NodeInfo
	listener       net.Listener
	rpcConn        *rpc.Client
	privateKey     string
	dataPublicKey  string
	dataPrivateKey string
	fileStatus     map[string][]byte
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
	dataPubKey, dataPriKey := generateKeys()

	node := Node{
		nodeInfo:       NodeInfo{localAddr, pubKey},
		listener:       listener,
		rpcConn:        connServer,
		privateKey:     priKey,
		dataPublicKey:  dataPubKey,
		dataPrivateKey: dataPriKey,
	}

	go func(rpcConn *rpc.Client, nodeInfo NodeInfo) {
		Log.Printf("Sending HeartBeat.. [Rate : %d]\n", HeartBeatRate)
		reply := false
		for {
			rpcConn.Call("Server.HeartBeat", nodeInfo, &reply)
			time.Sleep(HeartBeatRate * time.Millisecond)
		}
	}(node.rpcConn, node.nodeInfo)

	return &node
}

const (
	ROUTE   string = "ROUTE"
	GETFILE string = "GETFILE"
	END     string = "END"
	SEARCH  string = "SEARCH"
	RESULT  string = "RESULT"
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
	FInfo     FileInfo
	Data      Chunk
}

type Message struct {
	Op   []byte
	Dst  []byte
	Data []byte
}

func (n *Node) Incoming(arg Message, reply *bool) error {
	Log.Println("RPC - Received RPC Message...")

	var route Route
	var operation Operation

	errRoute := decryptStruct(arg.Dst, n.privateKey, &route)
	errOp := decryptStruct(arg.Op, n.privateKey, &operation)

	if errOp != nil || errRoute != nil {
		Log.Println("RPC - Decrypting Message Failed")
		return nil
	}

	Log.Printf("RPC - Message Decrypted [OP : %s] [DST : %s]\n", operation.Op, route.Dst)
	if operation.Op != ROUTE && operation.Op != END && operation.Op != GETFILE {
		Log.Println("RPC - Invalid OP")
		return nil
	}

	switch operation.Op {

	case GETFILE:
		var data Data
		errData := decryptStruct(arg.Data, n.privateKey, &data)
		if errData != nil {
			Log.Println("RPC - Decrypting Data Failed")
			return nil
		}

		chunk, err := getChunk(data.Data.Index, data.Data.Length, data.FInfo.Fname)
		if err != nil {
			Log.Println("RPC - GetFile: Unable to Get Chunk")
			return nil
		}

		data.Data = chunk
		encData, err := encryptStruct(data, data.PublicKey)
		if err != nil {
			Log.Println("RPC - Encrypting Data Failed")
			return nil
		}

		arg.Data = encData

	case END:
		var data Data
		errData := decryptStruct(arg.Data, n.dataPrivateKey, &data)
		if errData != nil {
			Log.Println("RPC - Encrypting Data Failed")
			return nil
		}

		errWrite := writeChunk(data.FInfo, data.Data)
		if errWrite != nil {
			Log.Println("FileManager - Writing Chunk Failed")
			return nil
		}

		n.fileStatus[data.FInfo.Fname][data.Data.Index] = 1

		return nil
	default:

	}

	msg := Message{operation.Next, route.Next, arg.Data}
	rpcConn, err := rpc.Dial("tcp", route.Dst)
	defer rpcConn.Close()
	if err != nil {
		Log.Println("RPC - Dial Failed")
		return nil
	}

	err = rpcConn.Call("Node.Incoming", msg, &reply)
	if err != nil {
		Log.Println("RPC - Sending Message to Next Route Failed")
		return nil
	}

	return nil
}
