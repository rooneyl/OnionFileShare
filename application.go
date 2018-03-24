package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"./node"
)

var client node.NodeAPI
var scanner *bufio.Scanner

func main() {
	if len(os.Args) != 3 {
		fmt.Println("USAGE : go run application.go localIP:Port serverIP:Port")
		return
	}

	nodeapi, err := node.Run(os.Args[1], os.Args[2], true)
	client = nodeapi
	scanner = bufio.NewScanner(os.Stdin)
	if err != nil {
		fmt.Printf("Error : Unable to Run Node")
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("::: Menu :::")
		fmt.Println("[1] Search/Download File")
		fmt.Println("[2] Change Directory")
		fmt.Println("[3] Exit")

		switch getInput() {
		case "1":
			search()
		case "2":
			changeDir()
		case "3":
			client.Disconnect()
			return

		default:
			fmt.Println("Invalid Input. Try Again")
		}
		fmt.Println()
	}
}

func getInput() string {
	fmt.Print("->")
	scanner.Scan()
	return scanner.Text()
}

func search() {
	fmt.Println("Enter Name of File or Press '0' to Return")
	fname := getInput()
	if fname == "0" {
		return
	}

	//fmt.Printf("%#v\n", client)
	fileInfo, err := client.Search(fname)
	if err != nil {
		fmt.Println("search error: ", err)
		return
	}
	if fileInfo == nil || len(fileInfo) == 0 {
		fmt.Printf("File [%s] Does Not Exists on the Network\n", fname)
		return
	}

	fmt.Println("Select File to Download or Press '0' to Return")
	for i, file := range fileInfo {
		fmt.Printf("File [%d] - Size : [%d]\n", i+1, file.Size)
	}

	for {
		selection, err := strconv.Atoi(getInput())
		if err != nil {
			fmt.Println("Invalid Input. Try Again")
			continue
		}

		if selection > len(fileInfo) || selection < 0 {
			fmt.Println("Invalid Input. Try Again")
			continue
		}

		//TODO
		err = client.GetFile(fileInfo[selection])
	}
}

func changeDir() {
	fmt.Println("Enter Desired Path or Press '0' to Return")
	for {
		path := getInput()
		if path == "0" {
			return
		}

		err := client.ChangeDirectory(path)
		if err != nil {
			fmt.Println("Invalid Path. Try Again")
			continue
		} else {
			fmt.Printf("Path Change Successful. New Path : [%s]\n", path)
			return
		}
	}
}
