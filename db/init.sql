-- Create the user_sessions table if it doesn't exist (DuckDB & PostgreSQL compatible)
CREATE TABLE IF NOT EXISTS jobbtid (
    -- Primary Key
    id BIGINT PRIMARY KEY,

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
