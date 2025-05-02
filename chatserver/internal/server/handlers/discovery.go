package handlers

import (
	"chatserver/internal/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func HandleDiscovery(identity db.ServerIdentity) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		KChostname := os.Getenv("KC_HOSTNAME")
		chatClientName := os.Getenv("CHAT_CLIENT_NAME")
		realmName := os.Getenv("REALM_NAME")
		hostname := os.Getenv("PUBLIC_HOSTNAME")

		log.Printf("KC_HOSTNAME: %s", KChostname)
		log.Printf("REALM_NAME: %s", realmName)

		URL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect", KChostname, realmName)
		response := map[string]string{
			"auth_url":    URL + "/auth",                    // URL to authenticate against
			"chat_client": chatClientName,                   // The keycloak client to connect through
			"chat_url":    "wss://chat." + hostname + "/ws", // Websocket URL for chat connections
			"token_url":   URL + "/token",
			"health_url":  fmt.Sprintf("%s/health", KChostname), // Health check URL for the authentication service
			"server_name": identity.Name,
			"server_id":   identity.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
