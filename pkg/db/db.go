/*
Copyright © 2025 Benjamin Jørgensen <me@benmi.me>
*/
package db

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/marcboeker/go-duckdb"
)

var (
	db             *sql.DB
	table_creation = `
-- Create the user_sessions table if it doesn't exist (DuckDB & PostgreSQL compatible)
CREATE TABLE IF NOT EXISTS jobbtid (
    -- Primary Key
    id BIGSERIAL PRIMARY KEY,

    -- User Information
    username VARCHAR(255) NOT NULL,

		-- Work Time
    jobbdag TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    starttime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    stoptime TIMESTAMP NULL,

    -- Record creation/update Timestamps
    create_dt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_dt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Audit User IDs
    create_uid VARCHAR(255) NULL,
    update_uid VARCHAR(255) NULL

		-- Delete state
		delete_flag TIMESTAMP NULL
);

-- Optional: Add indexes
CREATE INDEX IF NOT EXISTS idx_user_sessions_username ON user_sessions(username);
CREATE INDEX IF NOT EXISTS idx_user_sessions_starttime ON user_sessions(starttime);
CREATE INDEX IF NOT EXISTS idx_user_sessions_stoptime ON user_sessions(stoptime);
`
)

type Jobbtid struct {
	Id int

	Uid string

	JobbDag time.Time

	Starttime time.Time
	Stoptime  time.Time

	Create_dt time.Time
	Update_dt time.Time

	Create_uid string
	Update_uid string
}

func init() {
	db, err := sql.Open("duckdb", "?access_mode=READ_WRITE")
	if err != nil {
		// DB is essential
		// panic and tell the user why
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		// Pingcheck
		// panic and tell the user why
		panic(err)
	}
	setupTables()
	defer db.Close()
}

func setupTables() {
	setting := db.QueryRowContext(context.Background(), "SELECT current_setting('access_mode')")
	var accessMode string
	err := setting.Scan(&accessMode)
	if err != nil {
		log.Println("Could not get accessmode")
	} else {
		log.Printf("DB opened with access mode %s", accessMode)
	}

	_, err = db.ExecContext(context.Background(), table_creation)
	if err != nil {
		panic(err)
	}
}

func setupDbCon() (*sql.DB, error) {
	db, err := sql.Open("duckdb", "?access_mode=READ_WRITE")
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
	jobbtid string,
	starttime string,
	stoptime string,
) (int64, error) {
	db, err := setupDbCon()
	if err != nil {
		return -1, err
	}

	query := `INSERT INTO jobbtid VALUES(?, ?, ?, ?,CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?)`

	res, err := db.ExecContext(context.Background(), query, uid, jobbtid, starttime, stoptime, uid, uid)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func Update(
	uid string,
	jobbtid string,
	starttime string,
	stoptime string,
) (int64, error) {
	db, err := setupDbCon()
	if err != nil {
		return -1, err
	}

	query := `
        UPDATE jobbtid
        SET
            starttime = ?,
            stoptime = ?,
            update_dt = CURRENT_TIMESTAMP,
            update_uid = ?
        WHERE
            id = ? AND jobbtid = ?`

	res, err := db.ExecContext(context.Background(), query, starttime, stoptime, uid, jobbtid)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func Get(
	uid string,
) (*bytes.Buffer, error) {
	db, err := sql.Open("duckdb", "?access_mode=READ_WRITE")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(
		context.Background(), `
		SELECT (*)
		FROM jobbtid
		WHERE (uid = ? )`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buff bytes.Buffer

	for rows.Next() {
		j := new(Jobbtid)
		err := rows.Scan(j.Id, j.Uid, j.JobbDag, j.Starttime, j.Stoptime, j.Create_dt, j.Update_dt, j.Create_uid, j.Update_uid)
		if err != nil {
			log.Fatal(err)
		}
		jsonData, err := json.Marshal(j)
		if err != nil {
			return nil, fmt.Errorf("error marshaling jobbtid due to %s¥n", err)
		}
		err = appendJson(jsonData, &buff)
		if err != nil {
			return nil, fmt.Errorf("error appending %s", err)
		}
	}
	return &buff, nil
}

func List() (*bytes.Buffer, error) {
	db, err := setupDbCon()
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(
		context.Background(), `
		SELECT (*)
		FROM jobbtid`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buff bytes.Buffer

	for rows.Next() {
		j := new(Jobbtid)
		err := rows.Scan(j.Id, j.Uid, j.JobbDag, j.Starttime, j.Stoptime, j.Create_dt, j.Update_dt, j.Create_uid, j.Update_uid)
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
