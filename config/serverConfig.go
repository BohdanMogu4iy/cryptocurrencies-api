package config

import (
	"os"
)

type ServerConfigStruct struct {
	Port string
}

var ServerConfig *ServerConfigStruct

func init() {
	if PORT, ok := os.LookupEnv("PORT"); ok {
		ServerConfig = &ServerConfigStruct{
			Port: PORT,
		}
	} else {
		ServerConfig = &ServerConfigStruct{
			Port: "8000",
		}
	}
}
