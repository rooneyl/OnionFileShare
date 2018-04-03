package main

import (
	"./node"
	"time"
)

func main() {
	client, _ := node.Run("127.0.0.1:20000", "127.0.0.1:10000", true)
	client.ChangeDirectory("./data2")
	for {
		time.Sleep(time.Minute)
	}
}
