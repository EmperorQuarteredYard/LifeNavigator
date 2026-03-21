package models

import (
	"LifeNavigator/pkg/permission"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID        uint64         `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name            string                   `json:"name"`
	Owner           uint64                   `json:"owner"`
	Description     string                   `json:"description"`
	RefreshInterval uint32                   `json:"refresh_interval"`
	LastRefresh     time.Time                `json:"last_refresh"` //TODO 这东西得迁移进预算/账户
	Permission      permission.PermissionSet `json:"permission" gorm:"default:31"`

	// 关联的项目预算
	Budgets []ProjectBudget `gorm:"foreignKey:ProjectID" json:"budgets"`
	Users   []User          `gorm:"many2many:account_users;"`
}

func (m *Project) SetID(id uint64) {
	m.ID = id
}
func (m *Project) GetID() uint64 {
	return m.ID
}
