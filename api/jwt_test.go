package api

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	originalUser := os.Getenv("SA_USERNAME")
	originalPassword := os.Getenv("SA_PASSWORD")
	originalSalt := os.Getenv("SA_TOKEN_SALT")
	os.Setenv("SA_USERNAME", "sample_user")
	os.Setenv("SA_PASSWORD", "sample_password")
	os.Setenv("SA_TOKEN_SALT", "salt!")
	defer os.Setenv("SA_USERNAME", originalUser)
	defer os.Setenv("SA_PASSWORD", originalPassword)
	defer os.Setenv("SA_TOKEN_SALT", originalSalt)
	Setup()
	input := JWTUser{
		Username: "testuser",
	}

	jwt, err := createJwt(&input)
	assert.Nil(t, err)
	assert.NotEqual(t, "", jwt)

	parsed, err := parseJwt(jwt)
	assert.Nil(t, err)
	assert.Equal(t, input.Username, parsed.Username)
	assert.NotEqual(t, "", parsed.Expires)
	expires, err := time.Parse("2006-01-02T15:04:05Z", parsed.Expires)
	assert.Nil(t, err)
	assert.True(t, expires.After(time.Now()))

}
