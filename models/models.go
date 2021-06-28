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
		Email        string `schema:"email";json: "email"`
		Password     string `schema:"password";json: "password"`
		RefreshToken string `json:"refreshToken"`
		storage.StandardFields
	}
	TokenSchema struct {
		UserId       interface{} `json:"userId"`
		RefreshToken string `json:"refreshToken"`
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
