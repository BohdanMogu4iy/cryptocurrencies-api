package config

import (
	"os"
)

type ServerConfigStruct struct {
	Port string
	Version string
}

var ServerConfig *ServerConfigStruct

func init() {
	ServerConfig = &ServerConfigStruct{}
	if PORT, ok := os.LookupEnv("PORT"); ok {
		ServerConfig.Port = PORT
	} else {ServerConfig.Port = "8000"}
	if VERSION, ok := os.LookupEnv("VERSION"); ok {
		ServerConfig.Version = VERSION
	}else {ServerConfig.Version = "v1"}
}
