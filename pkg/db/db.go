/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package db

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	//_ "github.com/marcboeker/go-duckdb"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db      *sql.DB
	initSql = `
-- Create the user_sessions table if it doesn't exist (DuckDB & PostgreSQL compatible)
CREATE TABLE IF NOT EXISTS jobbtid (
    -- Primary Key
    id INTEGER PRIMARY KEY,

    -- User Information
    uid VARCHAR(255) NOT NULL,

		-- Work Time
    jobbdag TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    starttime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    stoptime TIMESTAMP NULL,

    -- Record creation/update Timestamps
    create_dt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_dt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Audit User IDs
    create_uid VARCHAR(255) NULL,
    update_uid VARCHAR(255) NULL,

		-- Delete state
		delete_flag TIMESTAMP NULL
);

-- Optional: Add indexes
CREATE INDEX IF NOT EXISTS idx_jobbtid_uid ON jobbtid(uid);
CREATE INDEX IF NOT EXISTS idx_jobbtid_starttime ON jobbtid(starttime);
CREATE INDEX IF NOT EXISTS idx_jobbtid_stoptime ON jobbtid(stoptime);
	`
)

type Jobbtid struct {
	Id int64 `json:"id"`

	Uid string `json:"userId"`

	JobbDag time.Time `json:"jobbDag"`

	Starttime time.Time `json:"startTime"`
	Stoptime  time.Time `json:"stopTime"`

	Create_dt time.Time `json:"createDt"`
	Update_dt time.Time `json:"updateDt"`

	Create_uid string `json:"createUserId"`
	Update_uid string `json:"updateUserId"`
}

func init() {
	db, err := setupDbCon()
	if err != nil {
		// DB is essential
		// panic and tell the user why
		panic(err)
	}
	defer db.Close()
	// setting := db.QueryRowContext(context.Background(), "SELECT current_setting('access_mode')")
	// var accessMode string
	// err = setting.Scan(&accessMode)
	// if err != nil {
	// 	log.Println("Could not get accessmode")
	// } else {
	// 	log.Printf("DB opened with access mode %s", accessMode)
	// }

	// initfile, err := os.ReadFile("db/init.sql")
	// if err != nil {
	// 	panic(err)
	// }
	_, err = db.ExecContext(context.Background(), string(initSql))
	if err != nil {
		panic(err)
	}
}

func setupDbCon() (*sql.DB, error) {
	// db, err := sql.Open("duckdb", "?access_mode=READ_WRITE")
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		return db, err
	}
	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

func Create(
	uid string,
	jobbdag string,
	starttime string,
	stoptime string,
) (int64, error) {
	db, err := setupDbCon()
	if err != nil {
		return -1, err
	}

	// query := `INSERT INTO jobbtid VALUES(?, ?, ?, ?,CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?, NULL)`
	query := `
INSERT INTO jobbtid (
    uid,
    jobbdag,
    starttime,
    stoptime,
    create_uid,
    update_uid
) VALUES (?, ?, ?, ?, ?, ?)`

	res, err := db.ExecContext(context.Background(), query, uid, jobbdag, starttime, stoptime, uid, uid)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func Update(
	id int64,
	uid string,
	jobbdag string,
	starttime string,
	stoptime string,
) (int64, error) {
	if starttime == "" && stoptime == "" {
		return -1, errors.New("you need either or both of starttime and stoptime")
	}
	db, err := setupDbCon()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	query := `
UPDATE jobbtid
SET `
	params := []any{}
	setClauses := []string{}

	if starttime != "" {
		setClauses = append(setClauses, "starttime = ?")
		params = append(params, starttime)
	}

	if stoptime != "" {
		setClauses = append(setClauses, "stoptime = ?")
		params = append(params, stoptime)
	}

	setClauses = append(setClauses, "update_dt = CURRENT_TIMESTAMP", "update_uid = ?")
	params = append(params, uid)

	query += strings.Join(setClauses, ", ") + `
WHERE
		uid = ?
	AND jobbdag = ?
	AND delete_flag IS NULL`

	params = append(params, uid, jobbdag)

	_, err = db.ExecContext(context.Background(), query, params...)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func GetById(
	uid string,
) ([]byte, error) {
	db, err := setupDbCon()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRowContext(
		context.Background(), `
SELECT
		id,
		uid,
		jobbdag,
		starttime,
		stoptime,
		create_dt,
		update_dt,
		create_uid,
		update_uid
FROM jobbtid
		WHERE 
			uid = ?
			AND delete_flag IS NULL`,
		uid,
	)

	j := new(Jobbtid)
	err = row.Scan(&j.Id, &j.Uid, &j.JobbDag, &j.Starttime, &j.Stoptime, &j.Create_dt, &j.Update_dt, &j.Create_uid, &j.Update_uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	jsonData, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("error marshaling jobbtid due to %s¥n", err)
	}

	return jsonData, nil
}

func GetByDate(
	userid string,
	jobbdag string,
) ([]byte, error) {
	db, err := setupDbCon()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRowContext(
		context.Background(), `
SELECT
		id,
		uid,
		jobbdag,
		starttime,
		stoptime,
		create_dt,
		update_dt,
		create_uid,
		update_uid
FROM jobbtid
WHERE 
	jobbdag = ? 
	AND uid = ?
	AND delete_flag IS NULL`,
		jobbdag, userid,
	)

	j := new(Jobbtid)
	err = row.Scan(&j.Id, &j.Uid, &j.JobbDag, &j.Starttime, &j.Stoptime, &j.Create_dt, &j.Update_dt, &j.Create_uid, &j.Update_uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	jsonData, err := json.Marshal(j)
	if err != nil {
		return nil, fmt.Errorf("error marshaling jobbtid due to %s¥n", err)
	}

	return jsonData, nil
}

func List() (*bytes.Buffer, error) {
	db, err := setupDbCon()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(
		context.Background(), `
SELECT
		id,
		uid,
		jobbdag,
		starttime,
		stoptime,
		create_dt,
		update_dt,
		create_uid,
		update_uid
FROM jobbtid
WHERE 
	AND delete_flag IS NULL`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buff bytes.Buffer

	for rows.Next() {
		j := new(Jobbtid)
		err := rows.Scan(&j.Id, &j.Uid, &j.JobbDag, &j.Starttime, &j.Stoptime, &j.Create_dt, &j.Update_dt, &j.Create_uid, &j.Update_uid)
		if err != nil {
			log.Fatal(err)
		}
		jsonData, err := json.Marshal(j)
		if err != nil {
			return nil, fmt.Errorf("error marshaling item due to %s¥n", err)
		}
		err = appendJson(jsonData, &buff)
		if err != nil {
			return nil, fmt.Errorf("error appending due to %s¥n", err)
		}
	}
	return &buff, nil
}

func appendJson(bytesData []byte, buff *bytes.Buffer) error {
	_, err := buff.Write(bytesData)
	if err != nil {
		return err
	}

	// Append a newline character to delimit JSON objects
	err = buff.WriteByte('\n')
	if err != nil {
		return err
	}
	return nil
}
