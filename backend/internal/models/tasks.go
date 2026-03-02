package models

import "time"

type Task struct {
	ID             uint64 `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserID         uint64 `json:"user_id" gorm:"index:idx_user_project"`
	ProjectID      uint64 `json:"project_id" gorm:"index:idx_user_project"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	AutoCalculated bool   `json:"auto_calculated" gorm:"default:false"`
	Type           uint8  `json:"type"`
	Status         uint8  `json:"status"`

	Category    string     `gorm:"column:category;type:varchar(50)" json:"category"`
	ForWhom     string     `gorm:"column:for_whom;type:varchar(100)" json:"for_whom"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// 关联的任务预算
	Budgets []TaskBudget `gorm:"foreignKey:TaskID" json:"budgets"`
}
type TaskDependency struct {
	ID        uint64    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserID    uint64    `json:"user_id" gorm:"index:idx_user_project"`
	ProjectID uint64    `json:"project_id" gorm:"index:idx_user_project"`
	TaskID    uint64    `gorm:"not null;index" json:"task_id"`    // 当前任务 ID，外键
	DependsOn uint64    `gorm:"not null;index" json:"depends_on"` // 依赖的任务 ID，外键
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间，自动管理
}
