package db

import (
	"os"
	"path/filepath"
	"sync"
)

var fwmux sync.Mutex
var frmux sync.Mutex

//storageManager represent place in the disk for storing the data
type storageManager interface {
	write(data []byte, offset int64) int64
	read(buffLength int64, offset int64) ([]byte, error)
	isfull() bool
	close()
	loadFile() (int64, map[string]int64)
}

//FileStorage saves datas to a file
type fileStorageManager struct {
	file *os.File
}

//This is thrown when file size exceeds allowed limit of 1GB ..
type storageFullError struct {
}

func (e *storageFullError) Error() string {
	return "Storage is full !!"
}

func newFileStorageManager(filePath string) (fileStorageManager, int) {

	var mode int
	var fs fileStorageManager

	//If not path provided, create data file in current location
	if filePath == "" {
		filePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	//check if 'data' file exists in the path, if not create the file
	//if the data file is already present, Open the file for read/write operation
	if !fileExists(filePath) {
		res, err := os.Create(filepath.Join(filePath, dataFile))
		handleError(err)
		fs = fileStorageManager{file: res}
		mode = os.O_CREATE

	} else {
		res, err := os.OpenFile(filepath.Join(filePath, dataFile), os.O_RDWR, 0600)
		handleError(err)
		fs = fileStorageManager{file: res}
		mode = os.O_RDWR
	}

	return fs, mode
}

//Close Closes the file
func (fs fileStorageManager) close() {
	fs.file.Close()
}

func (fs fileStorageManager) write(data []byte, offset int64) int64 {

	fwmux.Lock()
	bufflength, err := fs.file.WriteAt(data, offset)
	fwmux.Unlock()
	handleError(err)

	return int64(bufflength)
}

func (fs fileStorageManager) read(buffLength int64, offset int64) ([]byte, error) {

	buf := make([]byte, buffLength)
	frmux.Lock()
	_, err := fs.file.ReadAt(buf, offset)
	frmux.Unlock()

	handleError(err)
	return buf, err
}

func (fs fileStorageManager) isfull() bool {
	stat, err := fs.file.Stat()
	handleError(err)
	g := (((stat.Size() / 1024.0) / 1024.0) / 1024.0)
	return g >= maxstoragesize
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filepath.Join(filePath, dataFile))

	return !os.IsNotExist(err)
}

func (fs fileStorageManager) loadFile() (int64, map[string]int64) {

	cache := make(map[string]int64, 0)
	return fs.loadFileRec(rootPosition, cache)
}

func (fs fileStorageManager) loadFileRec(offset int64, cache map[string]int64) (int64, map[string]int64) {
	var isLeaf bool
	res, err := fs.read(capsize, offset)

	if err != nil && err.Error() == "EOF" {
		isLeaf = true
	}

	node := decode(res)

	if node.Key != "" {
		cache[node.Key] = node.Pos
	}

	if isLeaf {
		return node.NextPos, cache
	}

	return fs.loadFileRec(node.NextPos, cache)
}
