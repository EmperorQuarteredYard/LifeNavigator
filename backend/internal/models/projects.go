package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID        uint64         `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserID    uint64         `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name            string    `json:"name"`
	Description     string    `json:"description"`
	RefreshInterval uint32    `json:"refresh_interval"`
	LastRefresh     time.Time `json:"last_refresh"`
	MaxTaskID       uint64    `json:"MaxTaskID" gorm:"default:0"`

	// 关联的项目预算
	Budgets []ProjectBudget `gorm:"foreignKey:ProjectID" json:"budgets"`
}
