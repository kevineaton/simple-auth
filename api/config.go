package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// Config is a global configuration object that consolidates the various settings
var Config = Configuration{}

// Configuration is a global configuration object that consolidates the various settings
type Configuration struct {
	RootAPIPort         string
	Username            string
	Password            string
	TokenExpiresMinutes time.Duration
	RateLimit           float64
	TokenSalt           string
}

// Setup initializes the configuration for the server. It is an in-memory Global struct
func Setup() *chi.Mux {
	// first we parse the environment to populate the global config store
	port := envHelper("SA_API_PORT", "8090")
	Config.RootAPIPort = fmt.Sprintf(":%s", port)

	Config.Username = envHelper("SA_USERNAME", "")
	Config.Password = envHelper("SA_PASSWORD", "")

	if Config.Username == "" || Config.Password == "" {
		panic("SA_USERNAME and/or SA_PASSWORD are empty!")
	}

	Config.TokenSalt = envHelper("SA_TOKEN_SALT", "")
	if Config.TokenSalt == "" {
		panic("SA_TOKEN_SALT is empty and that is a security issue!")
	}

	rateLimit := envHelper("SA_RATE", "1000")
	rateLimiteParsed, err := strconv.ParseFloat(rateLimit, 64)
	if err != nil {
		panic("invalid SA_RATE rate limit set")
	}
	Config.RateLimit = rateLimiteParsed

	tokenExpiresMinutes := envHelper("SA_TOKEN_EXPIRES_MINUTES", "1440") // defaults to a day
	tokenExpiresMinutesParsed, err := strconv.ParseInt(tokenExpiresMinutes, 10, 64)
	if err != nil {
		panic("invalid SA_TOKEN_EXPIRES_MINUTES rate limit set")
	}
	tokenDuration, err := time.ParseDuration(fmt.Sprintf("%dm", tokenExpiresMinutesParsed))
	if err != nil {
		panic("invalid SA_TOKEN_EXPIRES_MINUTES rate limit set")
	}
	Config.TokenExpiresMinutes = tokenDuration

	// now setup the routes
	r := chi.NewRouter()
	middlewares := getMiddlewares()
	for _, m := range middlewares {
		r.Use(m)
	}

	r.Get("/", GetStatusRoute)
	r.Get("/health", GetHealthRoute)
	r.Get("/status", GetStatusRoute)

	r.Post("/verify", VerifyLoginRoute)
	r.Post("/validate", ValidateJWTRoute)

	return r
}

func envHelper(key, defaultValue string) string {
	found := os.Getenv(key)
	if found == "" {
		return defaultValue
	}
	return found
}

func getMiddlewares() []func(http.Handler) http.Handler {
	handlers := []func(http.Handler) http.Handler{}

	h := middleware.RequestID
	handlers = append(handlers, h)

	h = middleware.RealIP
	handlers = append(handlers, h)

	h = middleware.Logger
	handlers = append(handlers, h)

	h = middleware.Recoverer
	handlers = append(handlers, h)

	handlers = append(handlers, render.SetContentType(render.ContentTypeJSON))

	handlers = append(handlers, middleware.Timeout(120*time.Second))
	handlers = append(handlers, middleware.StripSlashes)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Types", "X-CSRF-TOKEN", "RANGE", "ACCEPT-RANGE"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	h = cors.Handler
	handlers = append(handlers, h)

	return handlers
}
