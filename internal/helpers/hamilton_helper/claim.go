package hamilton_helper

import (
	"encoding/base64"
	"encoding/json"
	"github.com/manicminer/hamilton/auth"
	"golang.org/x/oauth2"
	"strings"
)

func ParseClaims(token *oauth2.Token) (claims auth.Claims, err error) {
	if token == nil {
		return
	}
	jwt := strings.Split(token.AccessToken, ".")
	payload, err := base64.RawURLEncoding.DecodeString(jwt[1])
	if err != nil {
		return
	}
	err = json.Unmarshal(payload, &claims)
	return
}
