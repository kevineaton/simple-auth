package main

import (
	"math/rand"
	"net/http"
	"time"

	api "github.com/kevineaton/simple-auth/api"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	r := api.Setup()

	err := http.ListenAndServe(api.Config.RootAPIPort, r)
	if err != nil {
		panic(err.Error())
	}
}
