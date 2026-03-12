package dto

// CreateAccountRequest 创建账户请求
type CreateAccountRequest struct {
	Type    string  `json:"type" binding:"required"`
	Name    string  `json:"name" binding:"required"`
	Balance float64 `json:"balance"`
}

// AdjustBalanceRequest 调整余额请求
type AdjustBalanceRequest struct {
	Amount float64 `json:"amount" binding:"required"` // 正数增加，负数减少
}

// AccountModel 一般用的Account模型
type AccountModel struct {
	ID         uint64  `json:"id" binding:"required"`
	UserID     uint64  `json:"user_id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Balance    float64 `json:"balance"`
	NetBalance float64 `json:"net_balance"`
}
