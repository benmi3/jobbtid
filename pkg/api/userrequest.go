package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type UserRequest struct {
	Username string `json:"username" `
	Date     string `json:"date" `
	Time     string `json:"time" `
}

func (ur *UserRequest) setDate(d string) {
	if d == "" {
		dateNow := time.Now()
		ur.Date = dateNow.Format(time.DateOnly)
		return
	}
	ur.Date = d
}

func (ur *UserRequest) checkDate() {
	if ur.Date == "" {
		ur.setDate(ur.Date)
	}
}

func (ur *UserRequest) setTime(t string) {
	if t == "" {
		dateNow := time.Now()
		ur.Time = dateNow.Format(time.DateTime)
		return
	}
	ur.Time = t
}

func (ur *UserRequest) checkTime() {
	if ur.Time == "" {
		ur.setTime(ur.Time)
	}
}

func (ur *UserRequest) checkDateTime() {
	ur.checkDate()
	ur.checkTime()
}

func (ur *UserRequest) okUsername() bool {
	return ur.Username != ""
}

func (ur *UserRequest) handleRequest(r *http.Request) int {
	switch r.Method {
	case "GET":
		username := r.URL.Query().Get("username")
		date := r.URL.Query().Get("date")
		time := r.URL.Query().Get("time")
		ur.Username = username
		// this will change
		ur.setDate(date)
		ur.setTime(time)
	case "POST":
		err := json.NewDecoder(r.Body).Decode(&ur)
		if err != nil {
			return http.StatusBadRequest
		}
		ur.checkDateTime()

	default:
		return http.StatusMethodNotAllowed
	}
	if !ur.okUsername() {
		return http.StatusBadRequest
	}
	return http.StatusOK
}
