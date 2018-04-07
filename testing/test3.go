package main

import (
	"./node"
	"fmt"
	"os"
	//"log"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("USAGE : go run app.go localIP:Port serverIP:Port")
	}

	client, err := node.Run(os.Args[1], os.Args[2], true)
	if err != nil {
		fmt.Println("Error ::: Unable to Run Node")
		fmt.Printf("-> %s\n", err)
	}

	fileInfoArr := client.Search("file1")
	fileInfo := fileInfoArr[0]

	fileInfoArr2 := client.Search("file2")
	fileInfo2 := fileInfoArr2[0]

	err = client.GetFile(fileInfo)
	if err != nil {
		fmt.Println("Error ::: Unable to Get File")
		fmt.Printf("-> %s\n", err)
	}

	err = client.GetFile(fileInfo2)
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
