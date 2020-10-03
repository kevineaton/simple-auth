package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	api "github.com/kevineaton/simple-auth/api"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	r := api.Setup()

	fmt.Printf("\n====================== SIMPLE AUTH ====================")
	fmt.Printf("\n====================== Starting on port %s ============\n", api.Config.RootAPIPort)
	err := http.ListenAndServe(api.Config.RootAPIPort, r)
	if err != nil {
		panic(err.Error())
	}
}
