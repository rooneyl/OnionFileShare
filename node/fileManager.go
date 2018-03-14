package node

import ()

var Path = "../Data"

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
	var chunk Chunk
	return chunk, nil
}

func combineChunk(Chunks []Chunk) error {
	return nil
}

func serachFile(fname string) (FileInfo, error) {
	var fileInfo FileInfo
	return fileInfo, nil
}

func changeDir(path string) error {
	return nil
}
