package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"time"
)

var Log *log.Logger = log.New(os.Stdout, "SERVER ::: ", log.Ltime|log.Lshortfile)

type NodeInfo struct {
	Addr      string
	PublicKey []byte
}

type NodeStatus struct {
	Node NodeInfo
	time time.Time
}

type FileInfo struct {
	Fname string
	Size  int
	Hash  string
	Nodes []NodeInfo
}

type Server struct {
	ServerAddr string
	Nodes      map[string]NodeStatus
}

//var filesize = int64(340)
var localPath = "./serverfile/"

func main() {
	gob.Register(&FileInfo{})
	gob.Register(&NodeInfo{})
	var server Server

	if len(os.Args) != 2 {
		Log.Fatal("Usage - go run server.go ip:port")
	}

	// create serverlist and write credentials
	if _, err := os.Stat(filepath.Join(localPath, "serverList.txt")); err == nil {
		// file exists, add credentials
		err = nil
		f, err := os.OpenFile(filepath.Join(localPath, "serverList.txt"), os.O_APPEND|os.O_RDWR, 0644)
		scanner := bufio.NewScanner(f)
		registered := false
		for scanner.Scan() {
			if (scanner.Text()) == os.Args[1] {
				registered = true
			}
		}
		checkError(err)
		if !registered {
			_, err = f.WriteString(os.Args[1] + "\n")
		}
		checkError(err)
		err = f.Close()

		checkError(err)
	} else {
		// Create the file and add credentials

		f, err := os.Create(filepath.Join(localPath, "serverList.txt"))
		//f.Truncate(filesize)
		checkError(err)
		f.WriteString(os.Args[1] + "\n")
		f.Close()
	}

	server.ServerAddr = os.Args[1]
	server.Nodes = make(map[string]NodeStatus)

	listener, err := net.Listen("tcp", server.ServerAddr)
	if err != nil {
		Log.Fatal("Error - Unable to Establish Connection")
	}

	err = rpc.Register(&server)
	if err != nil {
		Log.Fatal("Error - RPC Register Failed")
	}

	Log.Printf("Running at [%s]", server.ServerAddr)
	go rpc.Accept(listener)

	for {
		time.Sleep(time.Second * 5)
		currentTime := time.Now()
		// fmt.Println("Availiable Node ->")
		for addr, node := range server.Nodes {
			if currentTime.After(node.time.Add(time.Second * 5)) {
				delete(server.Nodes, addr)
			} else {
				// fmt.Printf("Node [%s] - time [%s]\n", node.Node.Addr, node.time.String())
			}
		}
	}
	fmt.Println()
}

func (s *Server) HeartBeat(nodeInfo NodeInfo, reply *bool) error {
	s.Nodes[nodeInfo.Addr] = NodeStatus{nodeInfo, time.Now()}
	return nil
}

func (s *Server) Search(fileName string, reply *[]FileInfo) error {
	Log.Printf("RPC - Search...[%s]\n", fileName)
	fileSource := make(map[string]*FileInfo)
	for addr, nodeStatus := range s.Nodes {
		client, err := rpc.Dial("tcp", addr)
		defer client.Close()
		if err != nil {
			continue
		}

		var fileInfo FileInfo
		err = client.Call("Node.Search", fileName, &fileInfo)
		if err != nil {
			continue
		}

		if fileSource[fileInfo.Hash] == nil {
			fileInfo.Nodes = append(fileInfo.Nodes, nodeStatus.Node)
			fileSource[fileInfo.Hash] = &fileInfo
		} else {
			appendNode := append(fileSource[fileInfo.Hash].Nodes, nodeStatus.Node)
			fileSource[fileInfo.Hash].Nodes = appendNode
		}
	}

	for _, fileInfo := range fileSource {
		*reply = append(*reply, *fileInfo)
	}

	return nil
}

func (s *Server) GetNode(numNode int, nodes *[]NodeInfo) error {
	for _, node := range s.Nodes {
		*nodes = append(*nodes, node.Node)
		numNode--
		if numNode == 0 {
			break
		}
	}

	return nil
}

func SyncServers(s Server) {

	f, err := os.OpenFile(filepath.Join(localPath, "serverList.txt"), os.O_RDONLY, 0644)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if scanner.Text() != s.ServerAddr {
			client, err := rpc.Dial("tcp", scanner.Text())
			defer client.Close()
			if err != nil {
				continue
			}

			var nodes *[]NodeInfo
			err = client.Call("Server.GetNode", 1000, &nodes)
			if err != nil {
				continue
			}

			for _, nodeInfo := range *nodes {
				if _, exist := s.Nodes[nodeInfo.Addr]; exist {
				} else {
					//TODO INITIALIZE TTL PROPERLY
					s.Nodes[nodeInfo.Addr] = NodeStatus{Node: nodeInfo, time: time.Now()}
				}
			}
		}
	}

	err = f.Close()
	checkError(err)

}

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
