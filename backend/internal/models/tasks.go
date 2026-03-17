package models

import "time"

type Task struct {
	ID          uint64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ProjectID   uint64 `json:"project_id" gorm:"index:idx_user_project"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        uint8  `json:"type"`
	Status      uint8  `json:"status"`

	Category    string     `gorm:"column:category;type:varchar(50)" json:"category"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
type TaskDependency struct {
	UserID         uint64    `json:"user_id" gorm:"not null"`
	ProjectID      uint64    `json:"project_id" gorm:"not null"`
	TaskID         uint64    `gorm:"not null;index:dependency" json:"task_id"`         // 当前任务 ID，外键
	PrerequisiteID uint64    `gorm:"not null;index:dependency" json:"prerequisite_id"` // 前置的任务 ID，外键
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`                 // 创建时间，自动管理
}
