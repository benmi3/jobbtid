package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserRequest struct {
	Username string `json:"username" `
	Date     string `json:"date" `
	Time     string `json:"time" `
}

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func RespondWithCodeMessage(w http.ResponseWriter, code int, message string) {
	errResp := HttpResponse{
		StatusCode: code,
		Message:    message,
	}

	jsonData, err := json.Marshal(errResp)
	if err != nil {
		log.Printf("Error marshaling error response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing error response: %v", err)
	}
}
