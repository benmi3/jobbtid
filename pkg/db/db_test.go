/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package db

import (
	"encoding/json"
	"testing"
)

func TestDefaultDbSetup(t *testing.T) {
	t.Run("Debug Setup", func(t *testing.T) {
		_, err := setupDbCon()
		if err != nil {
			t.Error("Wanted default is true, got false")
		}
	})
}

func TestDefaultDbRun1(t *testing.T) {
	t.Run("One run db", func(t *testing.T) {
		uid := "testUser"
		jobbdag := "2025/5/10"
		starttime := "2025/5/10 10:10"
		stoptime := "2025/5/10 10:11"
		id, err := Create(uid, jobbdag, starttime, stoptime)
		if err != nil {
			t.Error("error happended during creation")
		}
		u_id, err := Update(id, uid, jobbdag, starttime, stoptime)
		if err != nil {
			t.Errorf("error happended during updating errror %q", err)
		}
		if u_id != id {
			t.Error("update id is not the same as create id")
		}
		testItem, err := GetByDate(uid, jobbdag)
		if err == nil {
			var myTestStruct Jobbtid
			marErr := json.Unmarshal(testItem, &myTestStruct)
			if marErr != nil {
				t.Error("failed at parsig GetByDate data")
			}
			if myTestStruct.Id != id {
				t.Errorf("myTestStruct.Id(%d) did not match id(%d) and update_id(%d)", myTestStruct.Id, id, u_id)
			}

		} else {
			t.Error("error occured when executing GetByDate")
		}
	})
}
