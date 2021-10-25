package model

import "time"

type Event struct {
	ID           int       `db:"id"`
	Title        string    `db:"title"`
	CreationTime time.Time `db:"creation_time"`
	PlanTime     time.Time `db:"plan_time"`
}
