package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/discovery", func(w http.ResponseWriter, r *http.Request) {
		hostname := os.Getenv("PUBLIC_HOSTNAME")
		log.Printf(hostname)
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
