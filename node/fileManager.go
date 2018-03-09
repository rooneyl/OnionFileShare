package node

import (
	"fmt"
)

const Path = "../Data"

type Chunk struct {
	Index  int
	Length int
	Data   []byte
}

type FileInfo struct {
	Fname string
	Size  int
	Hash  string
}

func getChunk(index int, length int, fname string) (Chunk, error) {
	return &Chunk, nil
}

func combineChunk(Chunks []Chunk) error {
	return nil
}

func serachFile(fname string) (FileInfo, error) {
	return &FileInfo, nil
}
