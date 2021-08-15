package main

import (
	"cryptocurrencies-api/app"
	"cryptocurrencies-api/models"
	"github.com/joho/godotenv"
	"os"
)

var (
	accountStorage = models.AccountStorage
	tokenStorage   = models.TokenStorage
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		return
	}

	err = os.Setenv("BASE_DIR",path)
	if err != nil {
		return
	}

	err = godotenv.Load()
	if err != nil {
		return
	}
	accountStorage.InitStorage()
	tokenStorage.InitStorage()
	app.RunServer()
}
