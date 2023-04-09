package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenJWT(m map[interface{}]interface{}, signinigKey []byte) (access, refresh string, err error) {
	var (
		accessToken, refreshToken *jwt.Token
		claims                    jwt.MapClaims
	)
	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)

	claims = accessToken.Claims.(jwt.MapClaims)
	rClaims := refreshToken.Claims.(jwt.MapClaims)

	for key, value := range m {
		claims[key.(string)] = value
		rClaims[key.(string)] = value
	}

	claims["iss"] = "user"
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().AddDate(0, 0, 1).Unix()

	rClaims["iss"] = "user"
	rClaims["iat"] = time.Now().Unix()
	rClaims["exp"] = time.Now().AddDate(0, 0, 30).Unix()

	accessTokenString, err := accessToken.SignedString(signinigKey)
	if err != nil {
		err = fmt.Errorf("access_token generating error: %s", err)
		return
	}

	refreshTokenString, err := refreshToken.SignedString(signinigKey)
	if err != nil {
		err = fmt.Errorf("refresh_token generating error: %s", err)
		return
	}

	return accessTokenString, refreshTokenString, nil
}

// ExtractClaims extracts claims from given token
func ExtractClaims(tokenString string, tokenSecretKey string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
