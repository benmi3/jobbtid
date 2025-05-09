/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"encoding/json"
	"jobbtid/pkg/db"
	"net/http"
)

func start(w http.ResponseWriter, r *http.Request) {
	var userReq UserRequest
	var err error
	var recId int64
	resCode := userReq.handleRequest(r)
	if resCode > http.StatusIMUsed {
		RespondWithCodeMessage(w, resCode, http.StatusText(resCode))
		return
	}

	recItem, err := db.GetByDate(userReq.Username, userReq.Date)
	if err == nil {
		if recItem == nil {
			recId, err = db.Create(userReq.Username, userReq.Date, userReq.Time, "")
		}
		var myJobbtid db.Jobbtid
		marErr := json.Unmarshal(recItem, &myJobbtid)
		if marErr == nil {
			recId, err = db.Update(myJobbtid.Id, userReq.Username, userReq.Date, userReq.Time, "")
		}
	}
	if err != nil || recId <= 0 {
		// recId is an auto increment BIGING, so if its 0 or less, something is wrong
		RespondWithCodeMessage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	resJson := CreatedId{
		ItemId: recId,
	}
	resBody, err := json.Marshal(resJson)
	if err != nil {
		RespondWithCodeMessage(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	RespondWithCodeBody(w, http.StatusCreated, resBody)
}
