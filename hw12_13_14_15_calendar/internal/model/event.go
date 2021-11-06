package model

import (
	"database/sql"
	"time"
)

type Event struct {
	ID                     int            `db:"id"`
	Title                  string         `db:"title"`
	CreationTime           time.Time      `db:"creation_time"`
	StartTime              time.Time      `db:"start_time"`
	EndTime                time.Time      `db:"end_time"`
	Description            sql.NullString `db:"description"`
	TimeBeforeNotification sql.NullTime   `db:"time_before_notification"`
}
