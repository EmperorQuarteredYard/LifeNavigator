package models

import (
	"time"

	"gorm.io/gorm"
)

type Kanban struct {
	ID        uint64         `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	UserID      uint64 `json:"user_id" gorm:"index"`
	Name        string `json:"name" gorm:"size:100"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default" gorm:"default:false"`
	SortOrder   int    `json:"sort_order" gorm:"default:0"`
}

type KanbanProject struct {
	KanbanID  uint64    `gorm:"primaryKey;autoIncrement:false" json:"kanban_id"`
	ProjectID uint64    `gorm:"primaryKey;autoIncrement:false" json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}
