package dto

import "time"

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	ProjectID      uint64  `json:"project_id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	AutoCalculated bool    `json:"auto_calculated"`
	Type           uint8   `json:"type"`
	Status         uint8   `json:"status"`
	Category       string  `json:"category"`
	ForWhom        string  `json:"for_whom"`
	Deadline       *string `json:"deadline"` // RFC3339 格式字符串，需在控制器中解析
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	ProjectID      uint64  `json:"project_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	AutoCalculated bool    `json:"auto_calculated"`
	Type           uint8   `json:"type"`
	Status         uint8   `json:"status"`
	Category       string  `json:"category"`
	ForWhom        string  `json:"for_whom"`
	Deadline       *string `json:"deadline"`
}

// TaskBudgetRequest 任务预算请求
type TaskBudgetRequest struct {
	Type   string  `json:"type" binding:"required"`
	Budget float64 `json:"budget" binding:"min=0"`
	Used   float64 `json:"used" binding:"min=0"`
}

// TaskResponse 任务详情响应
type TaskResponse struct {
	ID             uint64                `json:"id"`
	UserID         uint64                `json:"user_id"`
	ProjectID      uint64                `json:"project_id"`
	Name           string                `json:"name"`
	Description    string                `json:"description"`
	AutoCalculated bool                  `json:"auto_calculated"`
	Type           uint8                 `json:"type"`
	Status         uint8                 `json:"status"`
	Category       string                `json:"category"`
	ForWhom        string                `json:"for_whom"`
	Deadline       *time.Time            `json:"deadline,omitempty"`
	CompletedAt    *time.Time            `json:"completed_at,omitempty"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	Budgets        []*TaskBudgetResponse `json:"budgets,omitempty"`
}

// TaskBudgetResponse 任务预算响应
type TaskBudgetResponse struct {
	ID        uint64    `json:"id"`
	TaskID    uint64    `json:"task_id"`
	Type      string    `json:"type"`
	Budget    float64   `json:"budget"`
	Used      float64   `json:"used"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskBudgetSummaryResponse 任务预算汇总响应
type TaskBudgetSummaryResponse struct {
	Budgets []*TaskBudgetResponse `json:"budgets"`
}

type SetPrerequisiteRequest struct {
	PrerequisiteID uint64 `json:"prerequisite_id" binding:"required"`
	TaskID         uint64 `json:"task_id" binding:"required"`
}
type GetPrerequisiteRequest struct {
	TaskID uint64 `json:"task_id" binding:"required"`
}
type GetPostrequisiteRequest struct {
	PrerequisiteID uint64 `json:"prerequisite_id" binding:"required"`
}
type DependencyResponse struct {
	PrerequisiteID uint64 `json:"prerequisite_id" binding:"required"`
	TaskID         uint64 `json:"task_id" binding:"required"`
}
