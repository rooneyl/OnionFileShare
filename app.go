package main

import (
	"./node"
	"fmt"
	"log"
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

	// fileInfo := client.Search("file1.zip")
	// for _, s := range fileInfo {
	//     fmt.Printf("FileName - [%s] , FileSize - [%d]\n", s.Fname, s.Size)
	//     fmt.Printf("Hash - %s\n", s.Hash)
	//     fmt.Printf("Length [%d]\n", len(s.Nodes))
	//     fmt.Println(s)
	// }

	fileInfoArr := client.Search("file1")
	fileInfo := fileInfoArr[0]

	err = client.ChangeDirectory("data")
	if err != nil {
		fmt.Println("Error ::: Unable to Change Directory")
		fmt.Printf("-> %s\n", err)
	}

	err = client.GetFile(fileInfo)
	if err != nil {
		fmt.Println("Error ::: Unable to Get File")
		fmt.Printf("-> %s\n", err)
	}

	err = client.Disconnect()
	if err != nil {
		fmt.Println("Error ::: Unable to Disconnect")
		fmt.Printf("-> %s\n", err)
	}

	/**
	/ Client interface:
	/ client.Serach(filename string) []FileInfo
	/ client.ChangeDirectory(path string) error
	/ client.GetFile(filename string) error
	/ client.Disconnect() error
	**/
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
