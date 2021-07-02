package models

import (
	"cryptocurrencies-api/storage"
	"sync"
)

var (
	TokenStorage   storage.Storage
	AccountStorage storage.Storage
)

type (
	AccountSchema struct {
		Login        string `schema:"Login";json:"Login"`
		Password     string `schema:"Password";json:"Password"`
		storage.StandardFields
	}

	TokenSchema struct {
		UserId       interface{} `json:"UserId"`
		RefreshToken string `json:"RefreshToken"`
		storage.StandardFields
	}
)

func init() {
	var accountFileMutex sync.Mutex
	var accountFileReadWriteMutex sync.Mutex
	AccountStorage = storage.Storage{
		UnitSchema: AccountSchema{},
		File: storage.ConcurrencyFile{
			FileName:           "accounts.json",
			FileMutex:          &accountFileMutex,
			FileReadWriteMutex: &accountFileReadWriteMutex,
		},
	}

	var tokenFileMutex sync.Mutex
	var tokenFileReadWriteMutex sync.Mutex
	TokenStorage = storage.Storage{
		UnitSchema: TokenSchema{},
		File: storage.ConcurrencyFile{
			FileName:           "tokens.json",
			FileMutex:          &tokenFileMutex,
			FileReadWriteMutex: &tokenFileReadWriteMutex,
		},
	}
}
