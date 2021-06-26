package controllers

import (
	u "cryptocurrencies-api/utils"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
}

var TestController = func(w http.ResponseWriter, r *http.Request) {
	response := u.Message(true, "Test response")
	u.Respond(w, response)
}
