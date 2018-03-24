package main

import (
	"fmt"

	"./node"
)

func main() {
	fmt.Println("FileManager Test")

	fileManager := node.GetFileManager()
	//fname := "file1.zip"
	fname := "test.txt"
	length := 10

	fmt.Println("Search File")
	file, err := fileManager.SearchFile(fname)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("GetChunk")
	chunks := make([]node.Chunk, length)
	for i := 0; i < length; i++ {
		c, err := fileManager.GetChunk(i, length, fname)
		chunks[i] = c
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("ChangeDir")
	err = fileManager.ChangeDir("./Data2")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("GetDir: ", fileManager.GetDir())

	fmt.Println("WriteFile")
	err = fileManager.WriteFile(file)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("WriteChunk")
	for i := 0; i < length; i++ {
		err := fileManager.WriteChunk(file, chunks[i])
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("DoneWriting")
	b, err := fileManager.DoneWriting(file)
	if err != nil {
		fmt.Println(err)
	}

	if !b {
		fmt.Println("FAIL")
	} else {
		fmt.Println("SUCCESS")
	}

}
