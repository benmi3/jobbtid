/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"log"
	"net/http"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token != "Bearer valid-token" {
			log.Printf("Authentication failed: Invalid or missing token ('%s')", token)
			RespondWithCodeMessage(w, http.StatusUnauthorized, "Unauthorized: Invalid or missing authentication token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
