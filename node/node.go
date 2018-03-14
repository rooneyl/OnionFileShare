package node

import (
	"errors"
	"log"
	"os"
)

type NodeAPI struct {
	localAddr  string
	serverAddr string
	node       *Node
}

func Run(localAddr string, serverAddr string) (NodeAPI, error) {
	Log = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)

	nodeAPI := NodeAPI{
		localAddr:  localAddr,
		serverAddr: serverAddr,
		node:       StartConnection(localAddr, serverAddr),
	}

	return nodeAPI, nil
}

func (n *NodeAPI) Search(fileName string) []FileInfo {
	return nil
}

func (n *NodeAPI) ChangeDirectory(path string) error {
	return changeDir(path)
}

func (n *NodeAPI) GetFile(fileName string) error {
	if len(n.Search(fileName)) == 0 {
		return errors.New("File Unavailiable Online")
	}

	return nil
}

func (n *NodeAPI) Disconnect() error {
	return nil
}
