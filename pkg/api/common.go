package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type CreatedId struct {
	ItemId int64 `json:"id"`
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

	RespondWithCodeBody(w, code, jsonData)
}

func RespondWithCodeBody(w http.ResponseWriter, code int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err := w.Write(body)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
