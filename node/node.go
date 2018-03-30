package node

import (
	"net"
	"net/rpc"
	"time"
)

type NodeInfo struct {
	Addr      string
	PublicKey []byte
}

type Node struct {
	nodeInfo       NodeInfo
	listener       net.Listener
	rpcConn        *rpc.Client
	privateKey     []byte
	dataPublicKey  []byte
	dataPrivateKey []byte
	AESKey		   []byte
	fileStatus     map[string][]byte
}

func StartConnection(localAddr string, serverAddr string) *Node {
	Log.Println("Initiating Network Connection")

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		Log.Fatalf("Error ::: Listening at [%s] Failed\n", localAddr)
	}
	Log.Printf("Running Node at %s\n", localAddr)

	connServer, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		Log.Fatalf("Error ::: Connecting to Server Failed at [%s] Failed\n", serverAddr)
	}
	Log.Printf("Connected to Server at %s\n", serverAddr)

	pubKey, priKey := GenerateKeys()
	dataPubKey, dataPriKey := GenerateKeys()

	node := Node{
		nodeInfo:       NodeInfo{localAddr, pubKey},
		listener:       listener,
		rpcConn:        connServer,
		privateKey:     priKey,
		dataPublicKey:  dataPubKey,
		dataPrivateKey: dataPriKey,
		fileStatus:     make(map[string][]byte),
	}

	go func(rpcConn *rpc.Client, nodeInfo NodeInfo) {
		Log.Printf("Sending HeartBeat.. [Rate : %d]\n", HeartBeatRate)
		reply := false
		for {
			rpcConn.Call("Server.HeartBeat", nodeInfo, &reply)
			time.Sleep(HeartBeatRate * time.Millisecond)
		}
	}(node.rpcConn, node.nodeInfo)

	rpc.Register(&node)
	go rpc.Accept(node.listener)

	return &node
}

const (
	ROUTE   string = "ROUTE"
	GETFILE string = "GETFILE"
	END     string = "END"
	SEARCH  string = "SEARCH"
	RESULT  string = "RESULT"
)

type OpBox struct {
	AESKey []byte
	OperationData []byte
}

type Operation struct {
	Op string
	Next OpBox
}

//type Operation struct {
//	Op   string
//	Next []byte
//}

type RouteBox struct {
	AESKey []byte
	RouteData []byte
}

type Route struct {
	Dst string
	Next RouteBox
}

//type Route struct {
//	Dst  string
//	Next []byte
//}

type DataBox struct {
	AESKey 	  []byte
	DataData   []byte
}

type Data struct {
	PublicKey []byte
	FInfo     FileInfo
	Data      Chunk
}

//type Data struct {
//	PublicKey []byte
//	FInfo     FileInfo
//	Data      Chunk
//}

type Message struct {
	Op   OpBox
	Dst  RouteBox
	Data DataBox
}

type OpMsg struct {
	AESKey []byte
	Op     []byte
}

type DstMsg struct {
	AESKey []byte
	Dst    []byte
}

type DataMsg struct {
	AESKey []byte
	Data   []byte
}

//type Message struct {
//	Op   []byte
//	Dst  []byte
//	Data []byte
//}

func (n *Node) Incoming(arg Message, reply *bool) error {
	Log.Println("RPC - Received RPC Message...")

	var operation Operation
	var route Route

	//decrypt the Destination
	errRoute := DecryptStruct(arg.Dst.AESKey, arg.Dst.RouteData, n.privateKey, &route)
	//errRoute := DecryptStruct(arg.Dst, n.privateKey, &route)

	//decrypt the Operation
	errOp := DecryptStruct(arg.Op.AESKey, arg.Op.OperationData, n.privateKey, &operation)
	//errOp := DecryptStruct(arg.Op, n.privateKey, &operation)


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
		//decrypt the data
		errData := DecryptStruct(arg.Data.AESKey, arg.Data.DataData, n.privateKey, &data)
		//errData := DecryptStruct(arg.Data, n.privateKey, &data)
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
		aesKey, encData, err := EncryptStruct(data, data.PublicKey)
		//encData, err := EncryptStruct(data, data.PublicKey)
		if err != nil {
			Log.Println("RPC - Encrypting Data Failed")
			return nil
		}

		arg.Data.DataData = encData
		arg.Data.AESKey = aesKey
		//arg.Data = encData
		//arg.AESKey = easKey

	case END:
		var data Data
		errData := DecryptStruct(arg.Data.AESKey, arg.Data.DataData, n.dataPrivateKey, &data)
		//errData := DecryptStruct(arg.Data, n.dataPrivateKey, &data)
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
	//msg := Message{operation.Next, route.Next, arg.Data}
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

func (n *Node) Search(fileName string, reply *FileInfo) error {
	Log.Printf("RPC - Search...[%s]\n", fileName)
	fileInfo, err := searchFile(fileName)
	if err != nil {
		Log.Println(err)
	}
	*reply = fileInfo
	return err
}
