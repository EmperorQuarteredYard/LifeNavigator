package dto

import "time"

// CreateTaskPaymentRequest 创建任务付款请求
type CreateTaskPaymentRequest struct {
	BudgetID uint64  `json:"budget_id" binding:"required"` // 所属项目预算ID
	Amount   float64 `json:"amount" binding:"required"`    // 付款金额
}

// UpdateTaskPaymentRequest 更新任务付款请求
type UpdateTaskPaymentRequest struct {
	Amount float64 `json:"amount" binding:"required"` // 新的付款金额
}

// FinishTaskRequest 完成任务请求
type FinishTaskRequest struct {
	Time string `json:"time" binding:"required"` // RFC3339 完成时间
}

// PrerequisitesRequest 设置/取消前置任务请求
type PrerequisitesRequest struct {
	PrerequisiteID uint64 `json:"prerequisite_id" binding:"required"`
}

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	ProjectID      uint64  `json:"project_id" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	AutoCalculated bool    `json:"auto_calculated"`
	Type           uint8   `json:"type"`
	Status         uint8   `json:"status"`
	Category       string  `json:"category"`
	Deadline       *string `json:"deadline"` // RFC3339 格式字符串
}

// UpdateTaskRequest 更新任务请求（所有字段均为零值时表示不更新，Deadline 为空字符串表示清空）
type UpdateTaskRequest struct {
	Name           string `json:"name"`
	ProjectID      uint64 `json:"project_id"`
	Description    string `json:"description"`
	AutoCalculated bool   `json:"auto_calculated"`
	Type           uint8  `json:"type"`
	Status         uint8  `json:"status"`
	Category       string `json:"category"`
	Deadline       string `json:"deadline"` // RFC3339 格式，空字符串表示清空
}

// TaskResponse 任务详情响应
type TaskResponse struct {
	ID             uint64                 `json:"id"`
	UserID         uint64                 `json:"user_id"`
	ProjectID      uint64                 `json:"project_id"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	AutoCalculated bool                   `json:"auto_calculated"`
	Type           uint8                  `json:"type"`
	Status         uint8                  `json:"status"`
	Category       string                 `json:"category"`
	Deadline       *time.Time             `json:"deadline,omitempty"`
	CompletedAt    *time.Time             `json:"completed_at,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	Payments       []*TaskPaymentResponse `json:"payments,omitempty"` // 改用 Payments
}

// TaskPaymentResponse 任务付款响应
type TaskPaymentResponse struct {
	ID        uint64    `json:"id"`
	TaskID    uint64    `json:"task_id"`
	BudgetID  uint64    `json:"budget_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 可选：如果需要展示预算的附加信息，可额外添加字段（如预算类型、账户ID等）
	// 可通过服务层查询后填充
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Page  int64          `json:"page"`
	Total int64          `json:"total"`
	Size  int64          `json:"size"`
	List  []TaskResponse `json:"list"`
}

// SetPrerequisiteRequest 设置前置任务请求（备用）
type SetPrerequisiteRequest struct {
	PrerequisiteID uint64 `json:"prerequisite_id" binding:"required"`
	TaskID         uint64 `json:"task_id" binding:"required"`
}

// GetPrerequisiteRequest 获取前置任务请求（备用）
type GetPrerequisiteRequest struct {
	TaskID uint64 `json:"task_id" binding:"required"`
}

// GetPostrequisiteRequest 获取后置任务请求（备用）
type GetPostrequisiteRequest struct {
	PrerequisiteID uint64 `json:"prerequisite_id" binding:"required"`
}

// DependencyResponse 依赖关系响应
type DependencyResponse struct {
	PrerequisiteID uint64 `json:"prerequisite_id"`
	TaskID         uint64 `json:"task_id"`
}
