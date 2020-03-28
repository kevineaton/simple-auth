package api

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusRoutes(t *testing.T) {
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
	code, _, _ := TestAPICall(http.MethodGet, "/status", nil, GetStatusRoute)
	assert.Equal(t, http.StatusOK, code)
	code, _, _ = TestAPICall(http.MethodGet, "/health", nil, GetHealthRoute)
	assert.Equal(t, http.StatusOK, code)
}
