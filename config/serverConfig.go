package config

import (
	"log"
	"os"
	"path"
)

type serverConfigStruct struct {
	Port    string
	Version string
	DataDir string
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
	if basePath, ok := os.LookupEnv("BASE_DIR"); ok {
		ServerConfig.DataDir =  path.Join(basePath, "data")
	}else {log.Fatal("BASE_DIR ENVIRONMENT VARIABLES is absent")}
}
