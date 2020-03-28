package api

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type simpleLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Bind takes the HTTP request object and binds the body to the struct
func (data *simpleLogin) Bind(r *http.Request) error {
	return nil
}

type jwtValidation struct {
	JWT string `json:"jwt"`
}

// Bind takes the HTTP request object and binds the body to the struct
func (data *jwtValidation) Bind(r *http.Request) error {
	return nil
}

// VerifyLoginRoute verifies the submitted username and password
func VerifyLoginRoute(w http.ResponseWriter, r *http.Request) {
	input := simpleLogin{}
	render.Bind(r, &input)

	// simply check the passed in data with the environment
	if input.Username != Config.Username || input.Password != Config.Password {
		// they didn't match, so we send an error
		SendError(&w, r, http.StatusUnauthorized, "credential_failure", "invalid credentials", &map[string]interface{}{
			"username": input.Username,
		})
		return
	}

	// we will generate a JWT to send back that can be checked against the /validate end point
	jwt, err := createJwt(&JWTUser{
		Username: input.Username,
	})
	if err != nil {
		// there's really no reason this should fail
		SendError(&w, r, http.StatusInternalServerError, "jwt_creation_error", "could not create that jwt", &map[string]interface{}{
			"username": input.Username,
			"error":    err.Error(),
		})
		return

	}
	Send(w, http.StatusOK, map[string]string{
		"jwt": jwt,
	})
	return
}

// ValidateJWTRoute validates the JWT for a user
func ValidateJWTRoute(w http.ResponseWriter, r *http.Request) {
	input := jwtValidation{}
	render.Bind(r, &input)

	if input.JWT == "" {
		SendError(&w, r, http.StatusUnauthorized, "missing_jwt", "you must provide the JWT in the body of the request", &map[string]interface{}{})
		return
	}

	// parse it
	user, err := parseJwt(input.JWT)
	if err != nil {
		SendError(&w, r, http.StatusUnauthorized, "jwt_parse_failure", "invalid jwt", &map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	expires, err := time.Parse("2006-01-02T15:04:05Z", user.Expires)
	if err != nil {
		SendError(&w, r, http.StatusUnauthorized, "jwt_parse_failure", "invalid jwt", &map[string]interface{}{
			"expires": user.Expires,
			"error":   err.Error(),
		})
		return
	}

	if expires.Before(time.Now()) {
		SendError(&w, r, http.StatusUnauthorized, "jwt_expired", "jwt expired", &map[string]interface{}{
			"expires": user.Expires,
		})
		return
	}

	Send(w, http.StatusOK, user)
	return
}
