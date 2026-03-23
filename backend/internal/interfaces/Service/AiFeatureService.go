package Service

import (
	"context"
	"time"

	"LifeNavigator/internal/models"
)

// StreamEvent 表示 AI 流式响应事件
type StreamEvent struct {
	Type    string      // 事件类型，如 "project_created", "task_created", "budget_created", "stream_complete", "stream_error", "summary_content"
	Content interface{} // 事件内容，类型根据事件类型变化
}

// AIFeatureService 定义 AI 辅助功能接口
// 提供基于 AI 的项目规划、总结及用户画像更新等功能
type AIFeatureService interface {
	// ReduceProject 根据用户描述生成项目计划（流式返回创建事件）
	// 该方法会调用 AI 模型解析用户描述，生成项目、任务和预算，并通过事件通道实时返回创建结果。
	// 参数:
	//   ctx: 上下文，用于控制超时或取消
	//   userID: 用户ID，用于关联创建的资源
	//   projectDescription: 用户输入的项目描述文本
	//   accountSummary: 用户账户摘要信息，用于辅助预算规划（可为空字符串）
	//   eventChan: 事件通道，用于接收流式事件（调用方负责关闭通道，本方法会在返回前关闭通道）
	// 返回值:
	//   error: 可能返回 ErrInternal（AI调用失败、解析失败等）、ErrForbidden 等
	ReduceProject(ctx context.Context, userID uint64, projectDescription string, accountSummary string, eventChan chan<- StreamEvent) error

	// Summary 生成用户在指定时间段内的成就总结，并更新用户画像
	// 该方法会获取用户在该时间段内完成的任务和当前用户画像，调用 AI 生成总结，并通过事件通道流式返回总结内容。
	// 参数:
	//   ctx: 上下文
	//   userID: 用户ID
	//   startTime: 时间段开始时间
	//   endTime: 时间段结束时间
	//   eventChan: 事件通道，用于接收流式总结片段（调用方负责关闭通道，本方法会在返回前关闭通道）
	// 返回值:
	//   string: 完整的总结文本
	//   error: 可能返回 ErrInternal、ErrUserNotFound、ErrInvalidInput 等
	Summary(ctx context.Context, userID uint64, startTime, endTime time.Time, eventChan chan<- StreamEvent) (string, error)

	// UpdateUserProfile 更新用户画像
	// 参数:
	//   userID: 用户ID
	//   profile: 新的用户画像内容（通常由 AI 生成）
	// 返回值:
	//   error: 可能返回 ErrUserNotFound、ErrInternal
	UpdateUserProfile(userID uint64, profile string) error

	// GetUserCompletedTasks 获取用户在指定时间段内完成的任务列表和当前用户画像
	// 参数:
	//   userID: 用户ID
	//   startTime: 时间段开始
	//   endTime: 时间段结束
	// 返回值:
	//   []models.Task: 完成的任务列表
	//   string: 当前用户画像
	//   error: 可能返回 ErrUserNotFound、ErrInternal
	GetUserCompletedTasks(userID uint64, startTime, endTime time.Time) ([]models.Task, string, error)
}

const (
	EventTypeProjectCreated = "project_created"
	EventTypeTaskCreated    = "task_created"
	EventTypeBudgetCreated  = "budget_created"
	EventTypeComplete       = "stream_complete"
	EventTypeError          = "stream_error"
	EventTypeSummaryContent = "summary_content"
	StreamEndMarker         = "[STREAM_END]"
)
