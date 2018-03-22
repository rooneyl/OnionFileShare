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

	fileInfo := client.Search("file1.zip")
	for _, s := range fileInfo {
		fmt.Printf("FileName - [%s] , FileSize - [%d]\n", s.Fname, s.Size)
		fmt.Printf("Hash - %s\n", s.Hash)
		fmt.Printf("Length [%d]\n", len(s.Nodes))
		fmt.Println(s)
	}

	for {

	}
	/**
	/ Client interface:
	/ client.Serach(filename string) []FileInfo
	/ client.ChangeDirectory(path string) error
	/ client.GetFile(filename string) error
	/ client.Disconnect() error
	**/
}
