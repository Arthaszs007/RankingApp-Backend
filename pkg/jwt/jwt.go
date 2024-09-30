package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// defined the key to sign
var jwtKey = []byte("jwt_key")

type Claim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// generate JWT ,will return an error or token,need give username and expired timestamp
func GenerateJWT(username string, expirationTime time.Time) (string, error) {

	claim := &Claim{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	jwtToken, err := token.SignedString(jwtKey)
	if err != nil {

		return "", err
	}
	return jwtToken, nil
}

// pass a string token as param to veirfy ,if pass return true, or false
func VerifyJWT(tokenStr string) (string, error) {
	// verify the signkey in the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return "", err
	} else if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// get claims from the token and return username if pass
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, _ := claims["username"].(string)
		return username, nil
	} else {
		return "", fmt.Errorf("invalid token claims")
	}
}
