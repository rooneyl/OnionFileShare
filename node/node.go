package node

import (
	"log"
	"os"
)

type Node struct {
	localAddr  string
	serverAddr string
	PublicKey  string
	PrivateKey string
}

func Run(localAddr string, serverAddr string) (Node, error) {
	Log = log.New(os.Stdout, "INFO: ", log.Ltime|log.Lshortfile)
	Log.Println("Running Node")

	Node := Node{
		localAddr:  localAddr,
		serverAddr: serverAddr,
	}

	return Node, nil
}
