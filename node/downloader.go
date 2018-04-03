package node

import (
	"errors"
	"fmt"
	"math/rand"
	"net/rpc"
	"sync"
	"time"
)

type Downloader struct {
	fileStatus []bool
	fileNode   []NodeInfo
	file       FileInfo

	mutex      *sync.Mutex
	randomNode []NodeInfo
	nodeAPI    *NodeAPI
}

func StartDownloader(nodeAPI *NodeAPI) *Downloader {
	downloader := Downloader{
		nodeAPI: nodeAPI,
		mutex:   &sync.Mutex{},
	}
	rand.Seed(time.Now().Unix())
	return &downloader
}

func (d *Downloader) getFile(file FileInfo) error {
	length := file.Size/FileChunkSize + 1
	d.fileStatus = make([]bool, length)
	d.fileNode = file.Nodes
	d.file = file
	writeFile(file)

	d.updateRandomNode()
	for i := 0; i < len(d.fileStatus); i++ {
		d.requestChunk(i)
	}

	err := d.downloadStatus()
	if err != nil {
		return err
	}
	doneWriting(file)

	return nil
}

func (d *Downloader) requestChunk(index int) error {
	selectedNode := d.fileNode[rand.Int()%len(d.fileNode)]
	Log.Printf("Downloader - Requesting Chunk[%d] from [%s]\n", index, selectedNode.Addr)

	// Data
	dataMessage := DecryptedData{
		RSA:   d.nodeAPI.node.rsaPublic,
		File:  Chunk{index, len(d.fileStatus), nil},
		Finfo: d.file,
	}
	encryptedData := d.generateEncryptedMessage(selectedNode.PublicKey, dataMessage)

	// Routing
	routingMessage := DecryptedRouting{
		Operation: END,
	}
	encryptedRouting := d.generateEncryptedMessage(d.nodeAPI.node.rsaPublic, routingMessage)
	encryptedRouting, routingInfo := d.layerMessage(encryptedRouting, selectedNode)

	message := Message{encryptedRouting, encryptedData}
	conn, err := rpc.Dial("tcp", routingInfo.Addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	reply := false
	conn.Call("Node.Incoming", message, &reply)
	return nil
}

func (d *Downloader) layerMessage(encryptedMessage EncryptedMessage, selectedNode NodeInfo) (EncryptedMessage, NodeInfo) {
	var routingMessage DecryptedRouting

	routingNode := d.nodeAPI.node.nodeInfo
	for i := 0; i < MinNumRoute; i++ {
		routingMessage = DecryptedRouting{
			Operation:   ROUTING,
			Destination: routingNode.Addr,
			Next:        encryptedMessage,
		}
		routingNode = d.randomNode[rand.Int()%len(d.randomNode)]
		encryptedMessage = d.generateEncryptedMessage(routingNode.PublicKey, routingMessage)
	}

	routingMessage = DecryptedRouting{
		Operation:   GETFILE,
		Destination: routingNode.Addr,
		Next:        encryptedMessage,
	}
	encryptedMessage = d.generateEncryptedMessage(selectedNode.PublicKey, routingMessage)
	routingNode = selectedNode

	for i := 0; i < MinNumRoute; i++ {
		routingMessage = DecryptedRouting{
			Operation:   ROUTING,
			Destination: routingNode.Addr,
			Next:        encryptedMessage,
		}
		routingNode = d.randomNode[rand.Int()%len(d.randomNode)]
		encryptedMessage = d.generateEncryptedMessage(routingNode.PublicKey, routingMessage)
	}

	return encryptedMessage, routingNode
}

func (d *Downloader) generateEncryptedMessage(rsa []byte, struc interface{}) EncryptedMessage {
	aes := generateAESKey()

	encryptedByte, err := encryptData(aes, struc)
	if err != nil {
		Log.Fatal("Downloader - Encryption Failed")
	}

	encryptedESA, err := encryptAESKey(aes, rsa)
	if err != nil {
		Log.Fatal("Downloader - Encryption Failed")
	}

	encryptedMessage := EncryptedMessage{encryptedESA, encryptedByte}

	return encryptedMessage
}

func (d *Downloader) downloadStatus() error {
	for i := 0; i < MaxNumFileRequest; i++ {
		time.Sleep(time.Second * 3)

		d.mutex.Lock()
		complete := true
		counter := 0
		for _, status := range d.fileStatus {
			if status == false {
				complete = false
			} else {
				counter++
			}
		}
		d.mutex.Unlock()

		if complete {
			return nil
		}

		err := d.updateRandomNode()
		if err != nil {
			return err
		}

		err = d.updateFileNode()
		if err != nil {
			return err
		}

		for i, status := range d.fileStatus {
			if !status {
				d.requestChunk(i)
			}
		}
		fmt.Printf("File [%s] downloaded... [%d/%d]\n", d.file.Fname, counter, len(d.fileStatus))
	}

	return errors.New("DownLoader - Failed to Download File [downloadStatus]")
}

func (d *Downloader) receivedChunk(index int) {
	d.mutex.Lock()
	d.fileStatus[index] = true
	d.mutex.Unlock()
}

func (d *Downloader) updateRandomNode() error {
	err := d.nodeAPI.node.connServer.Call("Server.GetNode", MinNumRoute*10, &d.randomNode)
	if err != nil {
		return err
	}

	Log.Printf("Downloader - Updated RandomNodes, numNodes = [%d]", len(d.randomNode))
	return nil
}

func (d *Downloader) updateFileNode() error {
	fileInfo, err := d.nodeAPI.Search(d.file.Fname)
	if err != nil {
		return err
	}

	for _, file := range fileInfo {
		if file.Hash == d.file.Hash {
			d.fileNode = file.Nodes
			return nil
		}
	}
	return errors.New("DownLoader - FileNode No Longer Availiable")
}
