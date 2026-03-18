package dto

import "time"

type CreateKanbanRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	ProjectIDs  []uint64 `json:"project_ids"`
}

type UpdateKanbanRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ProjectIDs  []uint64 `json:"project_ids"`
	SortOrder   int      `json:"sort_order"`
}

type KanbanResponse struct {
	ID          uint64            `json:"id"`
	UserID      uint64            `json:"user_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	IsDefault   bool              `json:"is_default"`
	SortOrder   int               `json:"sort_order"`
	Projects    []KanbanProject   `json:"projects,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type KanbanProject struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type KanbanListResponse struct {
	Total int64             `json:"total"`
	List  []KanbanResponse  `json:"list"`
}

type KanbanTaskResponse struct {
	ID             uint64                  `json:"id"`
	ProjectID      uint64                  `json:"project_id"`
	ProjectName    string                  `json:"project_name"`
	Name           string                  `json:"name"`
	Description    string                  `json:"description"`
	Type           uint8                   `json:"type"`
	Status         uint8                   `json:"status"`
	StatusText     string                  `json:"status_text"`
	Category       string                  `json:"category"`
	Deadline       *time.Time              `json:"deadline,omitempty"`
	CompletedAt    *time.Time              `json:"completed_at,omitempty"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
	Prerequisites  []*TaskDependencyInfo   `json:"prerequisites,omitempty"`
	Payments       []*TaskPaymentResponse  `json:"payments,omitempty"`
}

type TaskDependencyInfo struct {
	TaskID         uint64 `json:"task_id"`
	TaskName       string `json:"task_name"`
	TaskStatus     uint8  `json:"task_status"`
	TaskStatusText string `json:"task_status_text"`
}

type KanbanTaskListResponse struct {
	Total int64                `json:"total"`
	List  []*KanbanTaskResponse `json:"list"`
}

type UpdateTaskStatusRequest struct {
	Status uint8 `json:"status" binding:"required"`
}
