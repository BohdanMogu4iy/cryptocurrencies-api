package main

import (
	"cryptocurrencies-api/app"
	"cryptocurrencies-api/models"
	"github.com/joho/godotenv"
)

var (
	accountStorage = models.AccountStorage
	tokenStorage = models.TokenStorage
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return 
	}
	accountStorage.InitStorage()
	tokenStorage.InitStorage()
	app.RunServer()
}
