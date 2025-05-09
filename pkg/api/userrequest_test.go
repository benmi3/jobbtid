/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package api

import (
	"net/http"
	"testing"
)

func TestUserRquestHandleOK(t *testing.T) {
	t.Run("No Query returns Bad Reqest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/stop?username=TEST", nil)

		var userReq UserRequest
		resCode := userReq.handleRequest(request)

		got := userReq.Username
		want := `TEST`

		got_date := userReq.Time
		dont_want_date := ``

		got_time := userReq.Time
		dont_want_time := ``

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if got_date == dont_want_date {
			t.Errorf("Date was %q, and it should not be that", got_date)
		}
		if got_time == dont_want_time {
			t.Errorf("Time was %q, and it should not be that", got_time)
		}

		if resCode != http.StatusOK {
			t.Errorf(" recode was %q, and it should be %q", resCode, http.StatusOK)
		}
	})
}

func TestUserRquestHandleNGDueToUsername(t *testing.T) {
	t.Run("No Query returns Bad Reqest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/stop?date=x&time=y", nil)

		var userReq UserRequest
		resCode := userReq.handleRequest(request)

		got := userReq.Username
		want := ``

		got_date := userReq.Time
		dont_want_date := ``

		got_time := userReq.Time
		dont_want_time := ``

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if got_date == dont_want_date {
			t.Errorf("Date was %q, and it should not be that", got_date)
		}
		if got_time == dont_want_time {
			t.Errorf("Time was %q, and it should not be that", got_time)
		}

		if resCode != http.StatusBadRequest {
			t.Errorf(" recode was %q, and it should be %q", resCode, http.StatusBadRequest)
		}
	})
}
