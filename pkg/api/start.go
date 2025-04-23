/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jobbtid/pkg/config"
	"net/http"
	"time"
)

func start(w http.ResponseWriter, r *http.Request) {
	var userReq UserRequest
	resCode := userReq.handleRequest(r)
	if resCode > http.StatusIMUsed {
		RespondWithCodeMessage(w, resCode, http.StatusText(resCode))
	}
	code, body, err := sendPostRequest(userReq)
	if err != nil || code >= http.StatusIMUsed {
		// If code is not 2XX then it will come as an error
		// the http.StatusIMUsed is just incase I change it later
		//  hehe
		RespondWithCodeMessage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
	// Temp just to satisfy lsp
	fmt.Fprint(w, body)
}

// http.StatusBadRequest (400): Malformed request, missing parameters.
// http.StatusUnauthorized (401): Authentication required.
// http.StatusForbidden (403): Authenticated but not authorized.
// http.StatusNotFound (404): Resource not found.
// http.StatusMethodNotAllowed (405): HTTP method not supported for the route.
// http.StatusInternalServerError (500): Generic server-side error.

var HOST_URL string = "https://project-eidetic/rpc"

func sendPostRequest(payload interface{}) (int, []byte, error) {
	baseURL := config.MainConfig.EideticConfig.Host

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, fmt.Errorf("error marshaling payload to JSON: %w", err)
	}
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// For now I will create a new client each time... might change later
	client := &http.Client{
		Timeout: 30 * time.Second, // Set a reasonable timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("error reading response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, responseBody, fmt.Errorf("received non-success status code %d: %s", resp.StatusCode, string(responseBody))
	}
	return resp.StatusCode, responseBody, nil
}
