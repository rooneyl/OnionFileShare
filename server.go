package main

import (
	//"crypto/ecdsa"
	//"crypto/elliptic"
	//"encoding/gob"
	//"encoding/json"
	//"errors"
	//"flag"
	//"fmt"
	//"io/ioutil"
	"log"
	//"math/rand"
	"net"
	"net/rpc"
	"os"
	//"sort"
	//"sync"
	//"time"
)

type (
	IP     string
	PubKey string
)

var (
	logger *log.Logger
	//replace first string with pubkey and second string with ip
	nodeList map[PubKey]IP
)

type SneakyNode struct {
	ip     IP
	pubKey PubKey
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		log.Println("Usage: go run server.go ip:port")
		return
	}

	nodeList = make(map[PubKey]IP)

	sneakyNode := new(SneakyNode)
	rpc.Register(sneakyNode)

	// Listen for new connection
	ln, err := net.ListenTCP("tcp", getAddr(":8080"))
	checkError(err)
	// Accept incoming connection
	rpc.Accept(ln)

	return
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getAddr(ip string) *net.TCPAddr {
	addr, err := net.ResolveTCPAddr("tcp", ip)
	checkError(err)
	log.Println("Listening on", addr)
	return addr
}
