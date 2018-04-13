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
	client.ChangeDirectory("./AppData")
	scanner = bufio.NewScanner(os.Stdin)
	if err != nil {
		fmt.Printf("Error : Unable to Run Node")
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("::: Menu :::")
		fmt.Println("[1] Search/Download File")
		fmt.Println("[2] My Files")
		fmt.Println("[3] Change Directory")
		fmt.Println("[4] Exit")

		switch getInput() {
		case "1":
			search()
		case "2":
			displayFiles()
		case "3":
			changeDir()
		case "4":
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
	fmt.Println("Enter Name of File or '0' to Return")
	fname := getInput()
	if fname == "0" {
		return
	}

	//fmt.Printf("%#v\n", client)
	fileInfos, err := client.Search(fname)
	if err != nil {
		fmt.Println("Search Error: ", err)
		return
	}
	if fileInfos == nil || len(fileInfos) == 0 {
		fmt.Printf("File [%s] Does Not Exist on the Network\n", fname)
		return
	}

	fmt.Println("Select File to Download or '0' to Return")
	for i, file := range fileInfos {
		fmt.Printf("[%d] File [%s] - Size [%d]\n", i+1, file.Fname, file.Size)
	}

	for {
		selection, err := strconv.Atoi(getInput())
		if err != nil {
			fmt.Println("Invalid Input. Try Again")
			continue
		}

		if selection == 0 {
			return
		}

		if selection > len(fileInfos) || selection < 0 {
			fmt.Println("Invalid Input. Try Again")
			continue
		}

		err = client.GetFile(fileInfos[selection-1])
		if err != nil {
			fmt.Println("Could not download this file: ", err)
			fmt.Println("Refreshing File List")
			fileInfos, err = client.Search(fname)
			if err != nil {
				fmt.Println("Search Error: ", err)
				return
			}
			if fileInfos == nil || len(fileInfos) == 0 {
				fmt.Printf("File [%s] Does Not Exist on the Network\n", fname)
				return
			}
			fmt.Println("Select File to Download or '0' to Return")
			for i, file := range fileInfos {
				fmt.Printf("[%d] File [%s] - Size [%d]\n", i+1, file.Fname, file.Size)
			}
			continue
		}
		fmt.Printf("File [%s] downloaded into current path: [%s]\n", fname, client.GetPath())
		return
	}
}

func changeDir() {
	fmt.Println()
	fmt.Printf("Current Path [%s]:\n", client.GetPath())
	dirs, err := client.ListDirs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, dir := range dirs {
		fmt.Println(" - ", dir)
	}

	fmt.Println("Enter Desired Path or '0' to Return")
	for {
		path := getInput()
		if path == "0" {
			fmt.Println()
			return
		}

		err := client.ChangeDirectory(path)
		if err != nil {
			fmt.Println("Invalid Path. Try Again")
			continue
		} else {
			fmt.Printf("Path Change Successful. New Path : [%s]\n", path)
			fmt.Println()
			return
		}
	}
}

func displayFiles() {
	fmt.Println()
	fmt.Printf("Current Path [%s]:\n", client.GetPath())
	fnames, err := client.ListFiles()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range fnames {
		fmt.Println(" - ", f)
	}

	fmt.Println("Enter '1' to Change Directory or '0' to Return")
	for {
		input := getInput()
		if input == "0" {
			return
		}
		if input == "1" {
			changeDir()
			fmt.Printf("Current Path [%s]:\n", client.GetPath())
			fnames, err := client.ListFiles()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, f := range fnames {
				fmt.Println(" - ", f)
			}
			fmt.Println("Enter '1' to Change Directory or '0' to Return")
		}
	}
}
