package config

import (
	"os"
)

type serverConfigStruct struct {
	Port string
	Version string
}

var ServerConfig *serverConfigStruct

func init() {
	ServerConfig = &serverConfigStruct{}
	if PORT, ok := os.LookupEnv("PORT"); ok {
		ServerConfig.Port = PORT
	} else {ServerConfig.Port = "8000"}
	if VERSION, ok := os.LookupEnv("VERSION"); ok {
		ServerConfig.Version = VERSION
	}else {ServerConfig.Version = "v1"}
}
