package node

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

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
	file, err := os.Open(fname)
	data := make([]byte, length)
	_, err = file.ReadAt(data, int64(index))

	if err == nil {
		chunk.Index = index
		chunk.Length = length
		chunk.Data = data
	}
	return chunk, err
}

func combineChunk(Chunks []Chunk) (combined []byte, err error) {
	sort.Slice(Chunks, func(i, j int) bool {
		return Chunks[i].Index < Chunks[j].Index
	})
	for _, chunk := range Chunks {
		combined = append(combined, chunk.Data...)
	}
	//TODO not sure what to return
	return combined, nil
}

func searchFile(fname string) (FileInfo, error) {
	var fileInfo FileInfo
	content, err := ioutil.ReadFile(fname)
	if err == nil {
		fileInfo.Fname = fname
		fileInfo.Size = len(content)
		fileInfo.Hash = string(content[:])
	}
	return fileInfo, err
}

func changeDir(path string) error {
	//TODO not sure
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	dir := filepath.Join(basepath, path)
	err := os.Chdir(dir)
	return err
}
