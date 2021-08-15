package config

import (
	u "cryptocurrencies-api/utils"
	"net/http"
)

type controllersConfigStruct struct {
	Messages  map[string]map[string]interface{}
	Responses map[string]func(w http.ResponseWriter)
}

var ControllersConfig *controllersConfigStruct

func init() {
	ControllersConfig = &controllersConfigStruct{
		Messages: map[string]map[string]interface{}{
			"MissingToken":            u.Message(false, "Authorization field is missing"),
			"InvalidToken":            u.Message(false, "Authentication token is invalid"),
			"NotRelevantToken":        u.Message(false, "Authentication token is not relevant. Login to refresh tokens."),
			"ExpiredOrNotActiveToken": u.Message(false, "Authentication token is either expired or not active yet"),
			"InternalServerError":     u.Message(false, "Internal Server Error"),
			"BadRequest":              u.Message(false, "Bad request"),
			"UserExists":              u.Message(false, "User already exists"),
			"UserCreated":             u.Message(true, "User has been created"),
			"InvalidCredentials":      u.Message(true, "Invalid credentials. Please try again or create new account"),
			"AOK":                     u.Message(true, "AOK, have a nice day!"),
		},
	}
	ControllersConfig.Responses = map[string]func(w http.ResponseWriter){
		"MissingToken": func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusBadRequest)
			u.Respond(w, ControllersConfig.Messages["MissingToken"])
		},
		"InvalidToken": func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, ControllersConfig.Messages["InvalidToken"])
		},
		"NotRelevantToken": func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, ControllersConfig.Messages["NotRelevantToken"])
		},
		"ExpiredOrNotActiveToken": func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, ControllersConfig.Messages["ExpiredOrNotActiveToken"])
		},
		"UserExists": func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, ControllersConfig.Messages["UserExists"])
		},
		"InvalidCredentials": func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusUnauthorized)
			u.Respond(w, ControllersConfig.Messages["InvalidCredentials"])
		},
		"BadRequest": func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusBadRequest)
			u.Respond(w, ControllersConfig.Messages["BadRequest"])
		},
		"InternalServerError": func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusInternalServerError)
			u.Respond(w, ControllersConfig.Messages["InternalServerError"])
		},
	}
}
