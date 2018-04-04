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

	var nodeAPI NodeAPI
	nodeAPI.localAddr = localAddr
	nodeAPI.serverAddr = serverAddr
	nodeAPI.downloader = StartDownloader(&nodeAPI)
	nodeAPI.node = StartConnection(localAddr, serverAddr, &nodeAPI)

	return nodeAPI, nil
}

func (n *NodeAPI) Search(fileName string) ([]FileInfo, error) {
	var fileInfo []FileInfo
	err := n.node.connServer.Call("Server.Search", fileName, &fileInfo)
	if err != nil {
		n.node.ReconnectS()
		if err != nil {
			log.Println(err)
		}
	}
	if err != nil {
		return nil, err
		Log.Fatal("Error ::: Connection with Server Unavailiable")
	}
	return fileInfo, err
}

func (n *NodeAPI) GetFile(file FileInfo) error {
	if len(file.Nodes) == 0 {
		return errors.New("Nodes Unavailiable or Already Has The File")
	}

	return n.downloader.getFile(file)
}

func (n *NodeAPI) ChangeDirectory(path string) error {
	return changeDir(path)
}

func (n *NodeAPI) GetPath() string {
	return getPath()

}

func (n *NodeAPI) ListDirs() ([]string, error) {
	return listDirs()

}

func (n *NodeAPI) ListFiles() ([]string, error) {
	return listFiles()

}

func (n *NodeAPI) Disconnect() error {
	err := n.node.listener.Close()
	if err != nil {
		return err
	}

	return n.node.connServer.Close()
}
