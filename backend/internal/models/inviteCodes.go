package models

import "time"

type InviteCode struct {
	ID           uint64 `gorm:"primary_key;auto_increment"`
	Token        string `json:"token"`
	Type         string `json:"type"`   //UserInvite/Construct
	Count        int    `json:"count"`  //已邀请数
	Amount       int    `json:"amount"` //邀请总额,-1记为无穷
	InvitedBy    string `json:"invited_by"`
	InviteAsRole string `json:"invite_as_role"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
