package controllers

import (
	"cryptocurrencies-api/config"
	"cryptocurrencies-api/models"
	u "cryptocurrencies-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/schema"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var decoder = schema.NewDecoder()

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.AccountSchema{}

	err := decoder.Decode(account, r.URL.Query())
	if err != nil {
		config.ControllersConfig.Responses["BadRequest"](w)
		return
	}

	selected, err := models.AccountStorage.SelectValues([]interface{}{account}, []string{"Login"})
	if err != nil{
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}
	if len(selected) > 0 {
		config.ControllersConfig.Responses["UserExists"](w)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	insertedUser, err := models.AccountStorage.InsertValues([]interface{}{account})
	if len(insertedUser) == 0 {
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}

	accessToken, err := GenerateToken(account.Id, config.JwtConfig.AccessTokenExpiresMinutes)
	if err != nil {
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}

	refreshToken, err := GenerateToken(account.Id, config.JwtConfig.RefreshTokenExpiresMinutes)
	if err != nil {
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}

	tokenUser := &models.TokenSchema{
		UserId: account.Id,
		RefreshToken: refreshToken,
	}

	insertedToken, err := models.TokenStorage.InsertValues([]interface{}{tokenUser})
	if len(insertedToken) == 0 {
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}

	response := config.ControllersConfig.Messages["AccountCreated"]
	response["accessToken"] = accessToken
	response["refreshToken"] = refreshToken

	u.Respond(w, response)
}

var LoginAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.AccountSchema{}

	err := decoder.Decode(account, r.URL.Query())
	if err != nil {
		config.ControllersConfig.Responses["BadRequest"](w)
		return
	}

	selected, err := models.AccountStorage.SelectValues([]interface{}{account}, []string{"Login"})
	if err != nil{
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}
	if len(selected) == 0 {
		config.ControllersConfig.Responses["InvalidCredentials"](w)
		return
	}

	password := u.GetField(selected, "Password").(string)

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(account.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		config.ControllersConfig.Responses["InvalidCredentials"](w)
		return
	}

	account.Id = r.Context().Value("UserId")

	accessToken, err := GenerateToken(account.Id, config.JwtConfig.AccessTokenExpiresMinutes)
	if err != nil {
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}

	userToken := models.TokenSchema{UserId: r.Context().Value("UserId")}
	if values, err := models.TokenStorage.SelectValues([]interface{}{&userToken}, []string{"userId"}); err != nil {
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	} else{
		refreshToken := u.GetField(values, "RefreshToken")
		response := config.ControllersConfig.Messages["AOK"]
		response["accessToken"] = accessToken
		response["refreshToken"] = refreshToken
		u.Respond(w, response)
	}
}

var RefreshToken = func(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value("UserId")

	accessToken, err := GenerateToken(userId, config.JwtConfig.AccessTokenExpiresMinutes)
	if err != nil {
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}

	refreshToken, err :=  GenerateToken(userId, config.JwtConfig.RefreshTokenExpiresMinutes)
	if err != nil {
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}

	userToken := models.TokenSchema{UserId: r.Context().Value("UserId")}
	if _, err := models.TokenStorage.UpdateValues([]interface{}{&userToken}); err != nil {
		config.ControllersConfig.Responses["InternalServerError"](w)
		return
	}

	response := config.ControllersConfig.Messages["AOK"]
	response["accessToken"] = accessToken
	response["refreshToken"] = refreshToken

	u.Respond(w, response)
}

func GenerateToken(userId interface{}, expireMinutes uint) (string, error) {
	exp := time.Now().Add(time.Minute * time.Duration(expireMinutes)).Unix()

	tokenClaims := &config.TokenClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	return jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenClaims).SignedString([]byte(config.JwtConfig.Secret))
}
