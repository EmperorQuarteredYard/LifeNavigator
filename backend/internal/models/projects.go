package models

import (
	"LifeNavigator/backend/pkg/refreshTime"
	"time"

	"gorm.io/gorm"
)

type Arrangement struct {
	ID        uint64         `Gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time      `Gorm:"auto_now_add" json:"created_at"`
	UpdatedAt time.Time      `Gorm:"auto_now" json:"updated_at"`
	DeletedAt gorm.DeletedAt `Gorm:"index" json:"deleted_at"`

	SonCount      uint64        `Gorm:"not null;default:0;column:son_count" json:"son_count"` //当子节点为零时，认为这是一个叶节点，此时下面自定义的TimeBudget等才认为是有效的，且可以作为天然的乐观锁
	OwnerID       uint64        `Gorm:"not null;default:0;column:owner_id" json:"owner_id"`
	IsRoot        bool          `Gorm:"not null;default:true" json:"is_root"`
	FatherID      uint64        `Gorm:"not null;column:father_id" json:"father_id"` //默认为自身
	Name          string        `Gorm:"not null;column:name" json:"name"`
	Description   string        `Gorm:"not null;longtext;column:description" json:"description"` //json格式，分为基础description和Steps,其中steps可以无线嵌套，但是Step只允许保留description字段
	EstimatedTime time.Duration `Gorm:"not null;default:0" json:"estimated_time"`
	ActualTime    time.Duration `Gorm:"not null;default:0" json:"actual_time"`
	Budget        float32       `Gorm:"not null;default:0" json:"budget"`
	ActualExpense float32       `Gorm:"not null;default:0" json:"actual_expense"`
	Currency      string        //币种，默认人民币，暂留空位不作处理
	Category      uint64        `Gorm:"not null;default:10000;column:category" json:"category"` //可能无效哈哈
	ForWhom       uint64        `Gorm:"default:0;column:for_whom" json:"for_whom"`              //这一项会用到用户的ID 和人物的私有ID
	Projects      []Project     `Gorm:"many2many:project_arrangements;" json:"projects"`
}

type Project struct {
	ID        uint64         `Gorm:"primary_key;auto_increment"`
	CreatedAt time.Time      `Gorm:"auto_now_add" json:"created_at"`
	UpdatedAt time.Time      `Gorm:"auto_now" json:"updated_at"`
	DeletedAt gorm.DeletedAt `Gorm:"index" json:"deleted_at"`

	Name            string        `Gorm:"size:255;not null" json:"name"`
	FundBudget      float32       `Gorm:"not null;default:0" json:"fund_budget"`
	FundUsed        float32       `Gorm:"not null;default:0" json:"fund_used"`
	FundLastRefresh time.Time     `Gorm:"not null;default:0" json:"fund_last_refresh"`
	RefreshStrategy string        `Gorm:"size:255;not null" json:"refresh_strategy"`
	RefreshGap      uint8         `Gorm:"not null;default:1" json:"refresh_gap"`
	Description     string        `Gorm:"not null"`
	Arrangements    []Arrangement `Gorm:"many2many:project_arrangements;" json:"arrangements"`
}

func RefreshProject(pro *Project) {
	if refreshTime.ShouldRefresh(pro.RefreshStrategy, pro.RefreshGap, pro.FundLastRefresh) {
		pro.FundUsed = 0
	}
}
