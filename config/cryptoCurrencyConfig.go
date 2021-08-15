package config

type cryptocurrencyConfigStruct struct {
	Currency string
}

var CryptocurrencyConfig *cryptocurrencyConfigStruct

func init() {
	CryptocurrencyConfig = &cryptocurrencyConfigStruct{Currency: "USD"}
}
