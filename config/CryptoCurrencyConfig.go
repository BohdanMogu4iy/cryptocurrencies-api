package config

type CryptocurrencyConfigStruct struct {
	Currency string
}

var CryptocurrencyConfig *CryptocurrencyConfigStruct

func init() {
	CryptocurrencyConfig = &CryptocurrencyConfigStruct{Currency: "USD"}
}
