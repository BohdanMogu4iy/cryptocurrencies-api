package main

import (
	"cryptocurrencies-api/app"
	"cryptocurrencies-api/models"
)

var (
	accountStorage = models.AccountStorage
	tokenStorage = models.TokenStorage
)

func main() {
	accountStorage.InitStorage()
	tokenStorage.InitStorage()
	app.RunServer()
}
