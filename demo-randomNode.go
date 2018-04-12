package main

import (
	"./node"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run demo-randomNode.go publicIP serverIP:Port")
		return
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	portNumber := r.Int()%40000 + 10000
	client, err := node.Run(os.Args[1]+":"+strconv.Itoa(portNumber), os.Args[2], false)
	if err != nil {
		fmt.Println(err)
	}
	err = client.ChangeDirectory("./Empty")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Node - [%s]\n", os.Args[1]+":"+strconv.Itoa(portNumber))

	for {
		time.Sleep(time.Minute)
	}
}
