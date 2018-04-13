package main

import (
	"./node"
	"fmt"

	"os"
	"time"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run demo-randomNode.go publicIP serverIP:Port")
		return
	}
	client, err := node.Run(os.Args[1], os.Args[2], false)
	if err != nil {
		fmt.Println(err)
	}
	err = client.ChangeDirectory("./Empty")
	if err != nil {
		fmt.Println(err)
	}

	for {
		time.Sleep(time.Minute)
	}
}
