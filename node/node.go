package node

import (
	"log"
	"os"
)

type NodeAPI struct {
	localAddr  string
	serverAddr string
	node       *Node
	downloader *Downloader
}

func Run(localAddr string, serverAddr string) (NodeAPI, error) {
	Log = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)

	nodeAPI := NodeAPI{
		localAddr:  localAddr,
		serverAddr: serverAddr,
		node:       StartConnection(localAddr, serverAddr),
	}

	nodeAPI.downloader = &Downloader{nodeAPI.node}

	return nodeAPI, nil
}

func (n *NodeAPI) Search(fileName string) []FileInfo {
	var fileInfo []FileInfo
	err := n.node.rpcConn.Call("Server.Search", fileName, &fileInfo)
	if err != nil {
		Log.Fatal("Error ::: Connection with Server Unavailiable")
	}
	return fileInfo
}

func (n *NodeAPI) GetFile(file FileInfo) error {
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
