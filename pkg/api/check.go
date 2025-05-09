/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"jobbtid/pkg/db"
	"net/http"
)

func check(w http.ResponseWriter, r *http.Request) {
	var userReq UserRequest
	resCode := userReq.handleRequest(r)
	if resCode > http.StatusIMUsed {
		RespondWithCodeMessage(w, resCode, http.StatusText(resCode))
		return
	}
	resBody, err := db.GetByDate(userReq.Username, userReq.Date)
	if err != nil {
		// recId is an auto increment BIGING, so if its 0 or less, something is wrong
		RespondWithCodeMessage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if resBody == nil {
		RespondWithCodeMessage(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	RespondWithCodeBody(w, http.StatusOK, resBody)
}
