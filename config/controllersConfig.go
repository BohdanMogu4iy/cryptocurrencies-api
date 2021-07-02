package config

import u "cryptocurrencies-api/utils"

type ControllersConfigStruct struct {
	Messages map[string]map[string]interface{}
}

var ControllersConfig *ControllersConfigStruct

func init() {
	ControllersConfig = &ControllersConfigStruct{
		Messages: map[string]map[string]interface{}{
			"MissingToken":                    u.Message(false, "Missing authentication Token"),
			"MalformedToken":                  u.Message(false, "Malformed authentication Token"),
			"InvalidToken":                    u.Message(false, "Invalid authentication Token"),
			"ValidationErrorSignatureInvalid": u.Message(true, "Invalid authentication Token signature"),
			"ValidationErrorClaimsInvalid":    u.Message(true, "Invalid authentication Token claims"),
			"ExpiredOrNotActiveToken":         u.Message(false, "Authentication Token is either expired or not active yet"),
			"NotRelevantToken":                u.Message(false, "Not relevant authentication Token"),
			"InternalServerError":             u.Message(false, "Internal Server Error"),
			"BadRequest":                      u.Message(false, "Bad request"),
			"AccountExists":                   u.Message(false, "Account already exists"),
			"AccountCreated":                  u.Message(true, "Account has been created"),
			"InvalidEmail":                    u.Message(true, "Invalid email. Check it"),
			"InvalidEmailOrPassword":          u.Message(true, "Invalid email or password. Please try again or create new account"),
			"InvalidPassword":                 u.Message(true, "Invalid login credentials. Please try again"),
			"AccountHasBeenCreated":           u.Message(true, "Account has been created"),
			"AOK":                             u.Message(true, "AOK, have a nice day!"),
		},
	}
}
