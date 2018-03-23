package node

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

type NodeAPI struct {
	localAddr  string
	serverAddr string
	node       *Node
	downloader *Downloader
}

func Run(localAddr string, serverAddr string, debug bool) (NodeAPI, error) {
	Log = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)
	if !debug {
		Log.SetFlags(0)
		Log.SetOutput(ioutil.Discard)
	}

	nodeAPI := NodeAPI{
		localAddr:  localAddr,
		serverAddr: serverAddr,
		node:       StartConnection(localAddr, serverAddr),
	}

	nodeAPI.downloader = &Downloader{nodeAPI.node}

	return nodeAPI, nil
}

func (n *NodeAPI) Search(fileName string) []FileInfo {
	Log.Println([]byte("file1.zip"))
	Log.Println([]byte(fileName))
	var fileInfo []FileInfo
	err := n.node.rpcConn.Call("Server.Search", fileName, &fileInfo)
	if err != nil {
		Log.Fatal("Error ::: Connection with Server Unavailiable")
	}
	Log.Println("XX")
	return fileInfo
}

func (n *NodeAPI) GetFile(file FileInfo) error {
	// for i, node := range file.Nodes {
	// if node.Addr == n.localAddr {
	// file.Nodes = append(file.Nodes[:i], file.Nodes[i+1:]...)
	// }
	// }

	if len(file.Nodes) == 0 {
		return errors.New("Nodes Unavailiable or Already Has The File")
	}

	return n.downloader.getFile(file)
}

func (n *NodeAPI) ChangeDirectory(path string) error {
	return changeDir(path)
}

func (n *NodeAPI) Disconnect() error {
	err := n.node.listener.Close()
	if err != nil {
		return err
	}

	return n.node.rpcConn.Close()
}
