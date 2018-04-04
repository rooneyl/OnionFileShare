package node

import (
	"fmt"
	"net"
	"net/rpc"
	"time"
)

type NodeInfo struct {
	Addr      string
	PublicKey []byte
}

type Node struct {
	nodeInfo   NodeInfo
	listener   net.Listener
	connServer *rpc.Client

	rsaPublic  []byte
	rsaPrivate []byte

	nodeAPI   *NodeAPI
	servers   []string
	servernum int
}

var servers = []string{}

func StartConnection(localAddr string, serverAddr string, nodeAPI *NodeAPI) *Node {
	Log.Println("Node - Initiating Network Connection")

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		Log.Fatalf("Node - Listening Fail [%s]\n", localAddr)
	}
	Log.Printf("Node - Running Node [%s]\n", localAddr)

	connServer, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		Log.Fatalf("Node - Connecting Server Fail [%s]\n", serverAddr)
	}
	Log.Printf("Node - Server Connected [%s]\n", serverAddr)

	publicKey, privateKey := generateRSAKey()

	node := Node{
		nodeInfo:   NodeInfo{localAddr, publicKey},
		listener:   listener,
		connServer: connServer,
		rsaPublic:  publicKey,
		rsaPrivate: privateKey,
		nodeAPI:    nodeAPI,
		servers:    []string{},
		servernum:  0,
	}
	node.connServer.Call("Server.GetServers", "", &servers)
	node.servers = servers

	go func(connServer *rpc.Client, nodeInfo NodeInfo) {
		reply := false
		for {
			err = connServer.Call("Server.HeartBeat", nodeInfo, &reply)
			if err != nil {
				node.ReconnectS()
				if err != nil {
					fmt.Println(err)
				}
			}
			time.Sleep(HeartBeatRate * time.Millisecond)
		}
	}(node.connServer, node.nodeInfo)

	rpc.Register(&node)
	go rpc.Accept(node.listener)

	return &node
}

const (
	ROUTING string = "ROUTING"
	GETFILE string = "GETFILE"
	END     string = "END"
)

type Message struct {
	Routing EncryptedMessage
	Data    EncryptedMessage
}

type EncryptedMessage struct {
	ESA  []byte
	Data []byte
}

type DecryptedRouting struct {
	Operation   string
	Destination string
	Next        EncryptedMessage
}

type DecryptedData struct {
	RSA   []byte
	File  Chunk
	Finfo FileInfo
}

func (n *Node) ReconnectS() (err error) {
	for n.servernum < len(servers) {
		n.connServer, err = rpc.Dial("tcp", servers[n.servernum])
		if err == nil {
			break
		}
		n.servernum++
	}
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) Incoming(msg Message, reply *bool) error {
	Log.Println("Node - Received Incoming Message")
	var routingMessage DecryptedRouting
	err := decrypting(msg.Routing, n.rsaPrivate, &routingMessage)
	if err != nil {
		return nil
	}

	Log.Printf("Node - Received OP = [%s]", routingMessage.Operation)

	switch routingMessage.Operation {
	case GETFILE:
		getFile(n, routingMessage, msg.Data)

	case ROUTING:
		routing(n, routingMessage, msg.Data)

	case END:
		end(n, msg.Data)
	}

	return nil
}

func (n *Node) Search(fileName string, reply *FileInfo) error {
	fileInfo, err := searchFile(fileName)
	*reply = fileInfo
	return err
}

func getFile(node *Node, routing DecryptedRouting, data EncryptedMessage) {
	Log.Println("Node - Processing GetFile")

	var dataMessage DecryptedData
	err := decrypting(data, node.rsaPrivate, &dataMessage)
	if err != nil {
		return
	}

	chunk, err := getChunk(dataMessage.File.Index, dataMessage.File.Length, dataMessage.Finfo.Fname)
	if err != nil {
		Log.Println("Node - Failed [getChunk]")
		return
	}
	dataMessage.File = chunk

	aesKey := generateAESKey()
	encryptedData, _ := encryptData(aesKey, dataMessage)
	encryptedAES, _ := encryptAESKey(aesKey, dataMessage.RSA)
	encryptedMessage := EncryptedMessage{encryptedAES, encryptedData}

	sendMessage(routing.Destination, routing.Next, encryptedMessage)
}

func routing(node *Node, routing DecryptedRouting, data EncryptedMessage) {
	Log.Printf("Node - Processing Routing to [%s]", routing.Destination)
	sendMessage(routing.Destination, routing.Next, data)
}

func end(node *Node, data EncryptedMessage) {
	Log.Println("Node - Processing END")

	var dataMessage DecryptedData
	err := decrypting(data, node.rsaPrivate, &dataMessage)
	if err != nil {
		return
	}

	err = writeChunk(dataMessage.Finfo, dataMessage.File)
	if err != nil {
		Log.Println("Node - Failed [writeChunk]")
		Log.Println(err)
		return
	}

	node.nodeAPI.downloader.receivedChunk(dataMessage.File.Index)
}

func decrypting(encryptedMessage EncryptedMessage, rsaPrivate []byte, v interface{}) error {
	aesKey, err := decryptAESKey(encryptedMessage.ESA, rsaPrivate)
	if err != nil {
		Log.Fatal("Node - Decrypting AES Failed")
		return err
	}

	err = decryptData(aesKey, encryptedMessage.Data, v)
	if err != nil {
		Log.Println("Node - Decrypting Data Failed")
		return err
	}

	return nil
}

func sendMessage(dst string, routing EncryptedMessage, data EncryptedMessage) {
	Log.Printf("Node - Sending Message to [%s]", dst)
	conn, err := rpc.Dial("tcp", dst)
	if err != nil {
		Log.Println("Node - Sending Message Failed")
		return
	}
	defer conn.Close()

	reply := false
	conn.Call("Node.Incoming", Message{routing, data}, &reply)
}
