package Repository

import (
	"LifeNavigator/internal/models"
	"time"
)

// TaskBudgetRepository 定义了任务预算（任务支付）的数据访问接口。
type TaskBudgetRepository interface {
	// Create 创建一个新的任务预算记录。
	Create(budget *models.TaskPayment) error

	// GetByID 根据主键 ID 查询任务预算。
	// 错误：如果记录不存在返回 ErrNotFound。
	GetByID(id uint64) (*models.TaskPayment, error)

	// GetByTaskID 根据任务 ID 查询该任务的所有预算项。
	GetByTaskID(taskID uint64) ([]models.TaskPayment, error)

	// Update 更新现有的任务预算记录。
	// 错误：如果记录不存在返回 ErrNotFound。
	Update(budget *models.TaskPayment) error

	// Delete 根据 ID 删除任务预算记录。
	// 错误：如果记录不存在返回 ErrNotFound。
	Delete(id uint64) error

	// DeleteByTaskID 删除指定任务的所有预算记录。
	DeleteByTaskID(taskID uint64) error

	// ListByAccountID 查询指定账户在时间范围内的所有任务支付记录。
	// 参数 startTime, endTime: 时间范围（包含边界），用于过滤任务的完成时间（completed_at）。
	// 实现行为：先通过 task_payments 表获取该账户关联的任务 ID，
	// 再查询这些任务中 completed_at 在时间范围内的任务支付记录。
	ListByAccountID(accountID uint64, startTime, endTime time.Time) ([]models.TaskPayment, error)
}
