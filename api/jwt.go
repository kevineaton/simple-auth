package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT is a user decrypted from a JWT token
type JWT struct {
	Username string `json:"username"`
	Expires  string `json:"expires"`
}

type jwtClaims struct {
	User JWT `json:"user"`
}

// createJwt creates a new jwt, sets an expiration, and creates the token
func createJwt(payload *JWT) (string, error) {
	// if the Expires isn't set, we need to set it to the expiration from the config
	// the only time it may be set is during test
	// generally, if you find yourself setting this by hand, you're doing it wrong
	if payload.Expires == "" {
		payload.Expires = time.Now().Add(Config.TokenExpiresMinutes).Format("2006-01-02T15:04:05Z")
	}
	asMap := payload.toMapClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, asMap)
	tokenString, err := token.SignedString([]byte(Config.TokenSalt))
	return tokenString, err
}

// parseJwt attempts to parse the JWT
func parseJwt(jwtString string) (*JWT, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(Config.TokenSalt), nil
	})
	if err != nil {
		return &JWT{}, errors.New("could not parse jwt")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		u := fromMapClaims(claims)
		return u, nil
	}
	return &JWT{}, errors.New("could not parse jwt")
}

func (data *JWT) toMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"username": data.Username,
		"exp":      jwt.NewNumericDate(time.Now().Add(Config.TokenExpiresMinutes)),
		"expires":  data.Expires,
	}
}

func fromMapClaims(data jwt.MapClaims) *JWT {
	u := JWT{}
	for k, v := range data {
		switch k {
		case "username":
			if s, ok := v.(string); ok {
				u.Username = s
			}
		case "expires":
			if s, ok := v.(string); ok {
				u.Expires = s
			}
		}
	}
	return &u
}
