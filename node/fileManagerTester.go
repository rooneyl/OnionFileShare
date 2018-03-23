package node

type fileManager int

func (f *fileManager) WriteFile(finfo FileInfo) error {
	return writeFile(finfo)
}
func (f *fileManager) DoneWriting(finfo FileInfo) (bool, error) {
	return doneWriting(finfo)
}
func (x *fileManager) GetChunk(i int, l int, f string) (Chunk, error) {
	return getChunk(i, l, f)
}
func (x *fileManager) WriteChunk(f FileInfo, chunk Chunk) error {
	return writeChunk(f, chunk)
}
func (x *fileManager) SearchFile(f string) (FileInfo, error) {
	return searchFile(f)
}

func (x *fileManager) ChangeDir(f string) error {
	return changeDir(f)
}

func GetFileManager() fileManager {
	var x fileManager
	return x
}
