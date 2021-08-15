package storage

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type ConcurrencyFile struct {
	FileName  string
	FileMutex *sync.Mutex
	FileReadWriteMutex *sync.Mutex
}

func (f *ConcurrencyFile) LockFile()  {
	f.FileMutex.Lock()
}

func (f *ConcurrencyFile) UnlockFile()  {
	f.FileMutex.Unlock()
}

func (f *ConcurrencyFile) Open() (*os.File, error) {
	f.FileReadWriteMutex.Lock()
	return os.Open(f.FileName)
}

func (f *ConcurrencyFile) Close(file *os.File) error {
	defer f.FileReadWriteMutex.Unlock()
	return file.Close()
}

func ThreadSafe(f *ConcurrencyFile, someFunc func(*os.File, []byte) ([]byte, error), data []byte) ([]byte, error) {
	file, err := f.Open()
	defer func(f *ConcurrencyFile, file *os.File) {
		err := f.Close(file)
		if err != nil {
			log.Fatal(err)
		}
	}(f, file)

	if err != nil {
		return nil, err
	}
	return someFunc(file, data)
}

func write(file *os.File, data []byte) ([]byte, error) {
	return nil, ioutil.WriteFile(file.Name(), data, 0644)
}

func (f *ConcurrencyFile) SafeWrite(data []byte) error{
	_, err := ThreadSafe(f, write, data)
	return err
}

func read(file *os.File, _ []byte) ([]byte, error) {
	return ioutil.ReadAll(file)
}

func (f *ConcurrencyFile) SafeRead() ([]byte, error){
	return ThreadSafe(f, read, nil)
}




