package dto

// CreateAccountRequest 创建账户请求
type CreateAccountRequest struct {
	Type    string  `json:"type" binding:"required"`
	Balance float64 `json:"balance"`
}

// AdjustBalanceRequest 调整余额请求
type AdjustBalanceRequest struct {
	Amount float64 `json:"amount" binding:"required"` // 正数增加，负数减少
}
