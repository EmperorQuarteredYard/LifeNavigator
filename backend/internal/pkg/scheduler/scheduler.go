package scheduler

import "time"

type AccountSchedule struct {
	AccountId   uint64    `json:"account_id"`
	RefreshTime time.Time `json:"refresh_time"`
}

type ProjectSchedule struct {
	ProjectId   uint64    `gorm:"primary_key"`
	RefreshTime time.Time `gorm:"primary_key"`
}
