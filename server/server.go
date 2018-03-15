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
	"time"
)

var (
	logger *log.Logger
	//first string is pubkey, second string is ip address
	nodeList map[string]string
)

type SneakyNode struct {
	ip     string
	pubKey string
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		log.Println("Usage: go run server.go ip:port")
		return
	}

	nodeList = make(map[string]string)

	sneakyNode := new(SneakyNode)
	rpc.Register(sneakyNode)

	// Listen for new connection
	ln, err := net.ListenTCP("tcp", getAddr(":8080"))
	checkError(err)
	// Accept incoming connection
	rpc.Accept(ln)

	return
}

func (s *SneakyNode) Hello(sn *SneakyNode, reply *string) error {
	register(sn)

	// run goroutine in infinite loop
	go HeartBeat(sn)

	*reply = "registered"
	return nil
}

func register(sn *SneakyNode) {
	if nodeExists(sn.pubKey) {
		return
	} else {
		nodeList[sn.pubKey] = sn.ip
		return
	}
}

func nodeExists(pubkey string) bool {
	for k, _ := range nodeList {
		if k == pubkey {
			return true
		}
	}
	return false
}

func HeartBeat(sn *SneakyNode) {
	for {
		_, err := rpc.Dial("tcp", sn.ip)
		if err != nil {
			for k, _ := range nodeList {
				if k == sn.pubKey {
					delete(nodeList, k)
				}
			}

		}

		time.Sleep(8 * time.Second)
	}
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
