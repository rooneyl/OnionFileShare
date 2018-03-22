package node

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path"
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
	tmpf, err := ioutil.TempFile(Path, fname)
	if err != nil {
		return err
	}
	defer tmpf.Close()
	b := make([]byte, finfo.Size)
	_, err = tmpf.Write(b)
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
	size := fi.Size() / int64(length)
	data := make([]byte, size)
	_, err = f.ReadAt(data, int64(index))
	if err == nil {
		chunk.Index = index
		chunk.Length = length
		chunk.Data = data
	}
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
	_, err = tmpf.WriteAt(chunk.Data, int64(chunk.Index))
	return err
}

func hashFile(f *os.File) (string, int, error) {
	h := md5.New()
	n, err := io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil)), int(n), err
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
		fileInfo.Size = size
		fileInfo.Hash = hash
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
