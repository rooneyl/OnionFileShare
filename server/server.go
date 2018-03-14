package server

import "fmt"

type Server struct {
	nodes []Node
}

type Node struct {
	addr      string
	publicKey stirng
}

func Run() {
	fmt.Println("server")
}

func (s *Server) HeartBeat(node Node) {
	//RPC Incomming call from node
	//HeartBeat Rate = 1 sec
}
