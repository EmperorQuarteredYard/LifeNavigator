package dto

import "time"

// CreateInviteCodeRequest 创建邀请码请求
type CreateInviteCodeRequest struct {
	Amount int    `json:"amount" binding:"required,min=1"`
	Role   string `json:"role" binding:"required"`
}

// InviteCodeResponse 邀请码详细信息
type InviteCodeResponse struct {
	Token     string     `json:"token"`
	CreatedBy uint64     `json:"created_by"`
	MaxUsage  int        `json:"max_usage"`
	UsedCount int        `json:"used_count"`
	Status    int        `json:"status"`               // 0:无效 1:有效
	ExpiresAt *time.Time `json:"expires_at,omitempty"` // 过期时间，可能为空
	CreatedAt time.Time  `json:"created_at"`
}
