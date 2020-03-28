package api

import "net/http"

// GetStatusRoute sees if the service is up
func GetStatusRoute(w http.ResponseWriter, r *http.Request) {
	Send(w, http.StatusOK, map[string]bool{
		"up": true,
	})
}

// GetHealthRoute sees if the service is healthy. Since there is no DB, cache, or anything, as soon as it can respond, it's healthy.
func GetHealthRoute(w http.ResponseWriter, r *http.Request) {
	Send(w, http.StatusOK, map[string]bool{
		"healthy": true,
	})
}
