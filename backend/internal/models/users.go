package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `gorm:"primary_key;auto_increment"`
	CreatedAt time.Time      `gorm:"auto_now_add" json:"created_at"`
	UpdatedAt time.Time      `gorm:"auto_now" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Username string `json:"username"`
	Password string `json:"-"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Version  uint64 `json:"version" gorm:"default:0"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Profile  string `json:"profile" gorm:"type:varchar(2000);default:''"` // 用户画像，最大2000字符
}
