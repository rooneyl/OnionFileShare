package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}
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
var localPath = "./server/serverfile/"

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
		checkError(err)
		f.WriteString(os.Args[1] + "\n")
		f.Close()
	}

	server.ServerAddr = os.Args[1]
	server.Nodes = make(map[string]NodeStatus)

	listener, err := net.ListenTCP("tcp", nil)
	add := listener.Addr().String()
	addrSt := strings.Split(add, ":")
	port := addrSt[len(addrSt)-1]
	server.ServerAddr = server.ServerAddr + ":" + port
	// listener, err := net.Listen("tcp", server.ServerAddr)
	if err != nil {
		Log.Fatal("Error - Unable to Establish Connection")
	}

	err = rpc.Register(&server)
	if err != nil {
		Log.Fatal("Error - RPC Register Failed")
	}

	Log.Printf("Running at [%s]", server.ServerAddr+":"+port)
	go rpc.Accept(listener)
	go Sync(server)

	online := 0
	for {
		time.Sleep(time.Second * 2)
		currentTime := time.Now()
		mutex.Lock()
		for addr, node := range server.Nodes {
			if currentTime.After(node.time.Add(time.Second * 5)) {
				delete(server.Nodes, addr)
			}
		}
		if len(server.Nodes) != online {
			online = len(server.Nodes)
			Log.Printf("Number of Node Online [%d]", online)
		}
		mutex.Unlock()
	}
}

func (s *Server) HeartBeat(nodeInfo NodeInfo, reply *bool) error {
	mutex.Lock()
	s.Nodes[nodeInfo.Addr] = NodeStatus{nodeInfo, time.Now()}
	mutex.Unlock()
	return nil
}

func (s *Server) Search(fileName string, reply *[]FileInfo) error {
	Log.Printf("RPC - Search...[%s]\n", fileName)
	fileSource := make(map[string]*FileInfo)
	mutex.Lock()
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
	mutex.Unlock()

	for _, fileInfo := range fileSource {
		*reply = append(*reply, *fileInfo)
	}

	return nil
}

func (s *Server) GetNode(numNode int, nodes *[]NodeInfo) error {
	mutex.Lock()
	for _, node := range s.Nodes {
		*nodes = append(*nodes, node.Node)
		numNode--
		if numNode == 0 {
			break
		}
	}
	mutex.Unlock()

	return nil
}

func (s *Server) GetServers(nonce string, servers *[]string) error {
	*servers = []string{}
	f, err := os.OpenFile(filepath.Join(localPath, "serverList.txt"), os.O_RDONLY, 0644)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		*servers = append(*servers, scanner.Text())
	}
	return err
}

func Sync(s Server) {
	for {
		SyncServers(s)
		time.Sleep(2 * time.Second)
	}
}

func SyncServers(s Server) {

	f, err := os.OpenFile(filepath.Join(localPath, "serverList.txt"), os.O_RDONLY, 0644)
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if scanner.Text() != s.ServerAddr {
			client, err := rpc.Dial("tcp", scanner.Text())
			if err != nil {
				//fmt.Println("server at " + scanner.Text() + " is down")
				break
			}

			var nodes *[]NodeInfo
			err = client.Call("Server.GetNode", 1000, &nodes)
			if err != nil {
				continue
			}

			for _, nodeInfo := range *nodes {
				if _, exist := s.Nodes[nodeInfo.Addr]; exist {
					// node already exists in s.Nodes, do nothing
				} else {
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
		log.Fatal(e)
	}
}
