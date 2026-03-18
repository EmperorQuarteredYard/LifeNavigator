package dto

import "time"

type CreateAccountRequest struct {
	Type    string  `json:"type" binding:"required"`
	Name    string  `json:"name" binding:"required"`
	Balance float64 `json:"balance"`
}

type AdjustBalanceRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

type Account struct {
	ID         uint64  `json:"id"`
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Unit       string  `json:"unit"`
	Balance    float64 `json:"balance"`
	NetBalance float64 `json:"net_balance"`
}

type AccountList struct {
	Items []*Account `json:"items"`
}

type TaskPayment struct {
	ID        uint64    `json:"id"`
	TaskID    uint64    `json:"task_id"`
	BudgetID  uint64    `json:"budget_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Task struct {
	ID          uint64         `json:"id"`
	ProjectID   uint64         `json:"project_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        uint8          `json:"type"`
	Status      uint8          `json:"status"`
	Category    string         `json:"category"`
	Deadline    *time.Time     `json:"deadline,omitempty"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Payments    []*TaskPayment `json:"payments,omitempty"`
}

type TaskList struct {
	Page  int64   `json:"page"`
	Total int64   `json:"total"`
	Size  int64   `json:"size"`
	List  []*Task `json:"list"`
}
