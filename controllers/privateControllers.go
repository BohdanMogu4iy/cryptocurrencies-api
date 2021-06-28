package controllers

import (
	"cryptocurrencies-api/config"
	u "cryptocurrencies-api/utils"
	"encoding/json"
	"fmt"
	"github.com/fatih/structs"
	"log"
	"net/http"
)

var TestController = func(w http.ResponseWriter, r *http.Request) {
	response := u.Message(true, "Test response")
	u.Respond(w, response)
}

type Currency struct {
	Code string `json:"code"`
	Rate string `json:"rate"`
	Description string `json:"description"`
	RateFloat float64 `json:"rate_float"`
}

type BpiStruct struct {
	USD Currency `json:"USD"`
	CNY Currency `json:"CNY"`
}

type UpdatedRate struct {
	Updated string `json:"updated"`
	UpdatedISO string `json:"updatedISO"`
}

type CryptocurrencyRate struct{
	Time UpdatedRate `json:"time"`
	Disclaimer string `json:"disclaimer"`
	Bpi BpiStruct `json:"bpi"`
}

type CryptocurrencyRateRequest struct {
	Currency string
}


var BtcRate = func(w http.ResponseWriter, r *http.Request) {

	req := &CryptocurrencyRateRequest{}

	err := decoder.Decode(req, r.URL.Query())
	if err != nil {
		req.Currency = config.CryptocurrencyConfig.Currency
	}

	url := fmt.Sprintf("https://api.coindesk.com/v1/bpi/currentPrice/%s.json", req.Currency)

	request, err := http.Get(url)
	if err != nil {
		log.Panic(err)
		return
	}

	btcRate := &CryptocurrencyRate{}
	err = json.NewDecoder(request.Body).Decode(&btcRate)
	if err != nil {
		log.Panic(err)
		return
	}

	u.Respond(w, structs.Map(btcRate))
}
