/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	var user UserRequest
	switch r.Method {
	case "GET":
		username := r.URL.Query().Get("username")
		date := r.URL.Query().Get("date")
		time := r.URL.Query().Get("time")
		user.Username = username
		// this will change
		user.Date = date
		user.Time = time
	case "POST":
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			// return HTTP 400 bad request
			RespondWithCodeMessage(w, http.StatusBadRequest, "Bad Request")
		}
	default:
		// For now, only post and get
		RespondWithCodeMessage(w, http.StatusBadRequest, "Bad Request")
	}

	fmt.Printf("Username is %s\n", user.Username)
}

// http.StatusBadRequest (400): Malformed request, missing parameters.
// http.StatusUnauthorized (401): Authentication required.
// http.StatusForbidden (403): Authenticated but not authorized.
// http.StatusNotFound (404): Resource not found.
// http.StatusMethodNotAllowed (405): HTTP method not supported for the route.
// http.StatusInternalServerError (500): Generic server-side error.
