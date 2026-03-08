package dto

type AIReduceProjectRequest struct {
	ProjectDescription string `json:"project_description" binding:"required"`
}

type AIReduceProjectStreamEvent struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type ProjectCreatedEvent struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TaskCreatedEvent struct {
	ID             uint64  `json:"id"`
	ProjectID      uint64  `json:"project_id"`
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Category       string  `json:"category"`
	Deadline       *string `json:"deadline,omitempty"`
	PrerequisiteID *uint64 `json:"prerequisite_id,omitempty"`
}

type BudgetCreatedEvent struct {
	ID        uint64  `json:"id"`
	ProjectID uint64  `json:"project_id"`
	Type      string  `json:"type"`
	Budget    float64 `json:"budget"`
}

type StreamCompleteEvent struct {
	Message string `json:"message"`
}

type StreamErrorEvent struct {
	Error string `json:"error"`
}

type AISummaryRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type AISummaryStreamEvent struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type SummaryContentEvent struct {
	Content string `json:"content"`
}

type SummaryCompleteEvent struct {
	Message       string `json:"message"`
	ProfileUpdate string `json:"profile_update"`
}
