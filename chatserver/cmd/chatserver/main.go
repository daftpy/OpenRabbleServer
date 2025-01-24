package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"auth_url":  "placeholder",
			"token_url": "placeholder",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Start the server on port 8080
	http.ListenAndServe("0.0.0.0:8080", nil)
}
