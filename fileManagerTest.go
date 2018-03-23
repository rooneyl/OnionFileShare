package main

import (
	"fmt"

	"./node"
)

func main() {
	fmt.Println("FileManager Test")

	fileManager := node.GetFileManager()

	fmt.Println("Search File")
	file, err := fileManager.SearchFile("file1.zip")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("GetChunk")
	chunks := make([]node.Chunk, 10)
	for i := 0; i < 10; i++ {
		c, err := fileManager.GetChunk(i, 10, "file1.zip")
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
	for i := 0; i < 10; i++ {
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
		fmt.Println(b)
	}

}
