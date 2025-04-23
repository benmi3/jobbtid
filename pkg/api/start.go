/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"fmt"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	var user UserRequest
	resCode := user.handleRequest(r)
	if resCode > http.StatusIMUsed {
		RespondWithCodeMessage(w, resCode, http.StatusText(resCode))
	}
	fmt.Printf("Username is %s\n", user.Username)
}

// http.StatusBadRequest (400): Malformed request, missing parameters.
// http.StatusUnauthorized (401): Authentication required.
// http.StatusForbidden (403): Authenticated but not authorized.
// http.StatusNotFound (404): Resource not found.
// http.StatusMethodNotAllowed (405): HTTP method not supported for the route.
// http.StatusInternalServerError (500): Generic server-side error.
