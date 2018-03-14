package main

import (
	"./Node"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("USAGE : go run app.go localIP:Port serverIP:Port")
	}

	client, err := node.Run(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println("Error ::: Unable to Run Node")
		fmt.Printf("-> %s\n", err)
	}

	_ = client

	/**
	/ Client interface:
	/ client.serach(filename string)
	/ client.changeDirectory(path string)
	/ client.getFile(filename string)
	/ client.disconnect()
	**/
}
