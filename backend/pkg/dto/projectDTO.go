package dto

import "time"

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description"`
	RefreshInterval uint32 `json:"refresh_interval"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description"`
	RefreshInterval uint32 `json:"refresh_interval"`
}

// ProjectBudgetRequest 项目预算请求（创建/更新共用）
type ProjectBudgetRequest struct {
	AccountID uint64  `json:"account_id"`                                   // 关联账户，0 表示未关联
	Type      string  `json:"type" binding:"oneof=time money token energy"` // 预算类型
	Budget    float64 `json:"budget" binding:"min=0"`                       // 预算总额
	Used      float64 `json:"used" binding:"min=0"`                         // 已用金额（创建时通常传0）
}

// ProjectResponse 项目详情响应
type ProjectResponse struct {
	ID              uint64                   `json:"id"`
	UserID          uint64                   `json:"user_id"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	RefreshInterval uint32                   `json:"refresh_interval"`
	LastRefresh     time.Time                `json:"last_refresh"`
	MaxTaskID       uint64                   `json:"max_task_id"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	Budgets         []*ProjectBudgetResponse `json:"budgets,omitempty"`
}

// ProjectBudgetResponse 项目预算响应
type ProjectBudgetResponse struct {
	ID        uint64    `json:"id"`
	ProjectID uint64    `json:"project_id"`
	AccountID uint64    `json:"account_id"` // 关联账户，0 表示未关联
	Type      string    `json:"type"`       // 预算类型：time/money/token/energy
	Budget    float64   `json:"budget"`
	Used      float64   `json:"used"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProjectBudgetSummaryResponse 项目预算汇总响应
type ProjectBudgetSummaryResponse struct {
	Budgets     []*ProjectBudgetResponse `json:"budgets"`
	TotalBudget float64                  `json:"total_budget"`
	TotalUsed   float64                  `json:"total_used"`
}

// ProjectListResponse 项目列表响应（分页）
type ProjectListResponse struct {
	Total int64              `json:"total"`
	Items []*ProjectResponse `json:"items"`
}
