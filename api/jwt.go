package api

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTUser is a user decrypted from a JWT token
type JWTUser struct {
	Username string `json:"username"`
	Expires  string `json:"expires"`
}

type jwtClaims struct {
	User JWTUser `json:"user"`
	jwt.StandardClaims
}

// createJwt creates a new jwt, sets an expiration, and creates the token
func createJwt(payload *JWTUser) (string, error) {
	// if the Expires isn't set, we need to set it to the expiration from the config
	// the only time it may be set is during test
	// generally, if you find yourself setting this by hand, you're doing it wrong
	if payload.Expires == "" {
		payload.Expires = time.Now().Add(Config.TokenExpiresMinutes).Format("2006-01-02T15:04:05Z")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": payload,
	})
	tokenString, err := token.SignedString([]byte(Config.TokenSalt))

	return tokenString, err
}

// parseJwt attempts to parse the JWT
func parseJwt(jwtString string) (JWTUser, error) {
	token, err := jwt.ParseWithClaims(jwtString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(Config.TokenSalt), nil
	})
	if err != nil {
		return JWTUser{}, errors.New("could not parse jwt")
	}

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		u := claims.User
		return u, nil
	}
	return JWTUser{}, errors.New("could not parse jwt")
}
