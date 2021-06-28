package storage

import (
	u "cryptocurrencies-api/utils"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

type StandardFields struct {
	Id int64 `json:"id"`
}

type StandardFileSchema struct {
	Items []interface{} `json:"items"`
	LastId int64 `json:"last_id"`
}

type Storage struct {
	UnitSchema interface{}
	File       ConcurrencyFile
}

func (s *Storage) fillFile(file *os.File) error  {
	s.File.FileReadWriteMutex.Lock()
	defer s.File.FileReadWriteMutex.Unlock()

	standardFileJson := &StandardFileSchema{Items: make([]interface{}, 0), LastId: 0}
	if b, err := json.Marshal(standardFileJson); err != nil {
		return err
	}else if _, err := file.Write(b); err != nil {
		return err
	}
	return nil
}

func (s *Storage) truncateFile(file *os.File) error  {
	s.File.FileReadWriteMutex.Lock()
	defer s.File.FileReadWriteMutex.Unlock()

	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	return nil
}

func (s *Storage) InitStorage(){
	s.File.LockFile()
	defer s.File.UnlockFile()

	file, err := os.OpenFile(s.File.FileName, os.O_RDWR, 0644)
	switch{
	case os.IsNotExist(err):
		if _, err := os.Create(s.File.FileName); err != nil{
			log.Fatal(err)
		}
		if err := s.fillFile(file); err != nil {
			log.Fatal(err)
		}
	case err != nil:
		log.Fatal(err)
	default:
		defer func(file *os.File) {
			if err := file.Close(); err != nil {
				log.Fatal(err)
			}
		}(file)
		if b, err := ioutil.ReadAll(file); err != nil{
			log.Fatal(err)
		}else {
			var FileData StandardFileSchema
			if err := u.UnmarshalJSON(b, &FileData); err != nil{
				if err := s.truncateFile(file); err != nil {
					log.Fatal(err)
				}
				if err := s.fillFile(file); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func (s *Storage) InsertValues(values []interface{}) error {
	s.File.LockFile()
	defer s.File.UnlockFile()

	oldByteValue, err := s.File.SafeRead()
	if err != nil {
		return err
	}

	var FileData StandardFileSchema
	if err := json.Unmarshal(oldByteValue, &FileData); err != nil {
		return err
	}

	for _, value := range values {
		reflect.ValueOf(value).Elem().FieldByName("StandardFields").Set(reflect.ValueOf(StandardFields{FileData.LastId+1}))
		FileData.Items = append(FileData.Items, value)
		FileData.LastId++
	}

	newByteValue, _ := json.MarshalIndent(FileData, " ", "\t")
	if err := s.File.SafeWrite(newByteValue); err != nil {
		return err
	}
	return nil
}

func (s Storage) UpdateValues(values []interface{}) error {
	s.File.LockFile()
	defer s.File.UnlockFile()

	oldByteValue, err := s.File.SafeRead()
	if err != nil {
		return err
	}

	var FileData StandardFileSchema
	if err := json.Unmarshal(oldByteValue, &FileData); err != nil {
		return err
	}

	for _, value := range values {
		valueJson := s.UnitSchema
		if valueByte, err := json.MarshalIndent(value, " ", "\t"); err ==  nil{
			if err := json.Unmarshal(valueByte, &valueJson); err != nil {
				return err
			}
		}else {return err}

		for i, item := range FileData.Items{
			if  valueJson.(map[string]interface{})["id"] == item.(map[string]interface{})["id"]{
				FileData.Items[i] = valueJson
			}
		}
	}

	newByteValue, _ := json.MarshalIndent(FileData, " ", "\t")
	if err := s.File.SafeWrite(newByteValue); err != nil {
		return err
	}
	return nil
}

func (s Storage) DeleteValues(values []interface{}) (int, error) {
	s.File.LockFile()
	defer s.File.UnlockFile()

	oldByteValue, err := s.File.SafeRead()
	if err != nil {
		return 0, err
	}

	var FileData StandardFileSchema
	if err := json.Unmarshal(oldByteValue, &FileData); err != nil {
		return 0, err
	}

	deletedCount := 0

	for _, value := range values {
		valueJson := s.UnitSchema
		if valueByte, err := json.MarshalIndent(value, " ", "\t"); err ==  nil{
			if err := json.Unmarshal(valueByte, &valueJson); err != nil {
				return 0, err
			}
		}else {return 0, err}

		for i, item := range FileData.Items{
			if  valueJson.(map[string]interface{})["id"] == item.(map[string]interface{})["id"]{
				FileData.Items = append(FileData.Items[:i], FileData.Items[i+1:]...)
				deletedCount++
			}
		}
	}

	newByteValue, _ := json.MarshalIndent(FileData, " ", "\t")
	if err := s.File.SafeWrite(newByteValue); err != nil {
		return 0, err
	}
	return deletedCount, nil
}

func (s Storage) SelectValues(values []interface{}) ([]interface{}, error) {
	s.File.LockFile()
	defer s.File.UnlockFile()

	oldByteValue, err := s.File.SafeRead()
	if err != nil {
		return nil, err
	}

	var FileData StandardFileSchema
	if err := json.Unmarshal(oldByteValue, &FileData); err != nil {
		return nil, err
	}

	selectedValues := make([]interface{}, 0)

	for _, value := range values {
		valueJson := s.UnitSchema
		if valueByte, err := json.MarshalIndent(value, " ", "\t"); err ==  nil{
			if err := json.Unmarshal(valueByte, &valueJson); err != nil {
				return nil, err
			}
		}else {return nil, err}

		for i, item := range FileData.Items{
			if  valueJson.(map[string]interface{})["id"] == item.(map[string]interface{})["id"]{
				selectedValues = append(selectedValues, FileData.Items[i])
			}
		}
	}

	newByteValue, _ := json.MarshalIndent(FileData, " ", "\t")
	if err := s.File.SafeWrite(newByteValue); err != nil {
		return nil, err
	}
	return selectedValues, nil
}
