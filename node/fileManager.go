package node

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

var Path = "../Data"

type Chunk struct {
	Fname  string
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
	f, err := os.Open(filepath.Join(Path, fname))
	if err != nil {
		return chunk, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return chunk, err
	}
	size := fi.Size() / int64(length)
	data := make([]byte, size)
	_, err = f.ReadAt(data, int64(index))
	if err == nil {
		chunk.Fname = fname
		chunk.Index = index
		chunk.Length = length
		chunk.Data = data
	}
	return chunk, err
}

func combineChunks(Chunks []Chunk) (fname string, combined []byte, b bool) {
	if len(Chunks) == 0 {
		return "", nil, false
	}
	if len(Chunks) > 1 {
		sort.Slice(Chunks, func(i, j int) bool {
			return Chunks[i].Index < Chunks[j].Index
		})
	}

	fname = Chunks[0].Fname
	for i, chunk := range Chunks {
		if chunk.Fname != fname || chunk.Index != i || Chunks[0].Length != len(Chunks) {
			return "", nil, false
		}
		combined = append(combined, chunk.Data...)
	}
	return fname, combined, true
}

func createFile(fname string, data []byte) error {
	name := filepath.Join(Path, fname)
	err := ioutil.WriteFile(name, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func searchFile(fname string) (FileInfo, error) {
	var fileInfo FileInfo
	f, err := os.Open(filepath.Join(Path, fname))
	if err != nil {
		return fileInfo, err
	}
	defer f.Close()
	h := md5.New()
	n, err := io.Copy(h, f)
	if err == nil {
		fileInfo.Fname = fname
		fileInfo.Size = int(n)
		fileInfo.Hash = hex.EncodeToString(h.Sum(nil))
	}
	return fileInfo, err
}

func changeDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	Path = path
	return nil
}
