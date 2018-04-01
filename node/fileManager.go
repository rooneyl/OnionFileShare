package node

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strings"
)

var Path = "./data"

type Chunk struct {
	Index  int
	Length int
	Data   []byte
}

type FileInfo struct {
	Fname string
	Size  int
	Hash  string
	Nodes []NodeInfo
}

func writeFile(finfo FileInfo) error {
	fname := finfo.Fname + ".tmp"
	b := make([]byte, finfo.Size)
	err := ioutil.WriteFile(path.Join(Path, fname), b, 0644)
	return err
}

func doneWriting(finfo FileInfo) (bool, error) {
	fname := finfo.Fname + ".tmp"
	tmp := path.Join(Path, fname)
	tmpf, err := os.Open(tmp)
	if err != nil {
		return false, err
	}
	defer tmpf.Close()

	hash, _, err := hashFile(tmpf)
	if err != nil {
		return false, err
	}
	ok, err := checkHash(hash, finfo, tmp)
	if err != nil {
		return false, err
	}
	tmpf.Close()
	err = os.RemoveAll(tmp)
	if ok {
		return true, err
	}
	return false, err
}

func checkHash(hash string, finfo FileInfo, tmp string) (bool, error) {
	if hash == finfo.Hash {
		data, _ := ioutil.ReadFile(tmp)
		err := ioutil.WriteFile(path.Join(Path, finfo.Fname), data, 0644)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func getChunk(index int, length int, fname string) (Chunk, error) {
	var chunk Chunk
	f, err := os.Open(path.Join(Path, fname))
	if err != nil {
		return chunk, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return chunk, err
	}
	size := int64(math.Ceil(float64(fi.Size()) / float64(length)))
	offset := int64(index) * size
	if (offset + size) > fi.Size() {
		size = fi.Size() - offset
	}
	data := make([]byte, size)
	n, err := f.ReadAt(data, offset)
	fmt.Println("chunk", index, "( size", size, ") read", n, "bytes :", offset, "to", (offset + int64(n)))
	if err == nil {
		chunk.Index = index
		chunk.Length = length
		chunk.Data = data
	}

	//fmt.Println("chunk", index, ":", chunk.Data)
	return chunk, err
}

func writeChunk(finfo FileInfo, chunk Chunk) error {
	name := finfo.Fname + ".tmp"
	tmp := path.Join(Path, name)
	tmpf, err := os.OpenFile(tmp, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer tmpf.Close()

	size := int64(math.Ceil(float64(finfo.Size) / float64(chunk.Length)))
	offset := int64(chunk.Index) * size
	n, err := tmpf.WriteAt(chunk.Data, offset)
	fmt.Println("chunk", chunk.Index, "wrote", n, "bytes:", offset, "to", (offset + int64(n)))
	//fmt.Println("chunk", chunk.Index, ":", chunk.Data)
	return err
}

func hashFile(f *os.File) (string, int64, error) {
	h := md5.New()
	n, err := io.Copy(h, f)
	hash := hex.EncodeToString(h.Sum(nil))
	fmt.Println(f.Name(), " ", hash)
	return hash, n, err
}

func searchFile(fname string) (FileInfo, error) {
	var fileInfo FileInfo
	f, err := os.Open(path.Join(Path, fname))
	if err != nil {
		return fileInfo, err
	}
	defer f.Close()

	hash, size, err := hashFile(f)
	if err == nil {
		fileInfo.Fname = fname
		fileInfo.Size = int(size)
		fileInfo.Hash = hash
		fmt.Println("File size", fileInfo.Size, "bytes")
	}
	return fileInfo, err
}

func changeDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	Path = strings.ToLower(path)
	return nil
}

func getDir() string {
	return Path
}

func getDirs() ([]string, error) {
	var dirs []string
	files, err := ioutil.ReadDir(Path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		//fmt.Println(file.Name())
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	return dirs, err
}
