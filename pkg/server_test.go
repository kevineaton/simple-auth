package simpleauth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicValidationRoutes(t *testing.T) {
	username := "sample_user"
	password := "sample_password"
	originalUser := os.Getenv("SA_USERNAME")
	originalPassword := os.Getenv("SA_PASSWORD")
	originalSalt := os.Getenv("SA_TOKEN_SALT")
	os.Setenv("SA_USERNAME", username)
	os.Setenv("SA_PASSWORD", password)
	os.Setenv("SA_TOKEN_SALT", "salt!")
	defer os.Setenv("SA_USERNAME", originalUser)
	defer os.Setenv("SA_PASSWORD", originalPassword)
	defer os.Setenv("SA_TOKEN_SALT", originalSalt)

	Setup()
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)

	enc.Encode(map[string]string{})
	code, _, _ := TestAPICall(http.MethodPost, "/verify", b, VerifyLoginRoute)
	assert.Equal(t, http.StatusUnauthorized, code)

	b.Reset()
	enc.Encode(map[string]string{
		"username": "fake",
		"password": "won't work",
	})
	code, _, _ = TestAPICall(http.MethodPost, "/verify", b, VerifyLoginRoute)
	assert.Equal(t, http.StatusUnauthorized, code)

	b.Reset()
	enc.Encode(map[string]string{
		"username": username,
		"password": password,
	})

	code, res, _ := TestAPICall(http.MethodPost, "/verify", b, VerifyLoginRoute)
	require.Equal(t, http.StatusOK, code)
	_, body, _ := UnmarshalTestMap(res)
	jwt, jwtOK := body["jwt"].(string)
	assert.True(t, jwtOK)
	assert.NotEqual(t, "", jwt)

	// check the JWT is valid
	user, err := parseJwt(jwt)
	assert.Nil(t, err)
	assert.Equal(t, username, user.Username)

	// validate it with the server
	b.Reset()
	enc.Encode(map[string]string{
		"jwt": jwt,
	})

	code, _, _ = TestAPICall(http.MethodPost, "/validate", b, ValidateJWTRoute)
	assert.Equal(t, http.StatusOK, code)

	// try some garbage JWT calls to the validate

	// start with an expired JWT
	expiredJWT, _ := createJwt(&JWTUser{
		Expires:  time.Now().AddDate(0, -2, 0).Format("2006-01-02T15:04:05Z"),
		Username: username,
	})

	b.Reset()
	enc.Encode(map[string]string{
		"jwt": expiredJWT,
	})

	code, _, _ = TestAPICall(http.MethodPost, "/validate", b, ValidateJWTRoute)
	assert.Equal(t, http.StatusUnauthorized, code)

	// an invalid expires
	invalidExpiresJWT, _ := createJwt(&JWTUser{
		Expires:  "November 28th",
		Username: username,
	})

	b.Reset()
	enc.Encode(map[string]string{
		"jwt": invalidExpiresJWT,
	})

	code, _, _ = TestAPICall(http.MethodPost, "/validate", b, ValidateJWTRoute)
	assert.Equal(t, http.StatusUnauthorized, code)

	// a missing jwt
	b.Reset()
	enc.Encode(map[string]string{
		"jwt": "",
	})

	code, _, _ = TestAPICall(http.MethodPost, "/validate", b, ValidateJWTRoute)
	assert.Equal(t, http.StatusUnauthorized, code)

	// finally, try an obviously fake jwt
	b.Reset()
	enc.Encode(map[string]string{
		"jwt": "afakejwt",
	})

	code, _, _ = TestAPICall(http.MethodPost, "/validate", b, ValidateJWTRoute)
	assert.Equal(t, http.StatusUnauthorized, code)

}
