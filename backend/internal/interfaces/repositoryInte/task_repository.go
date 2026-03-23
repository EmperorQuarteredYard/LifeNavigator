package repositoryInte

import (
	"LifeNavigator/internal/models"
	"time"
)

// TaskRepository 定义了任务相关的数据访问接口。
type TaskRepository interface {
	// Create 创建新任务。
	Create(task *models.Task) error

	// GetByID 根据任务 ID 查询任务信息。
	GetByID(id uint64) (*models.Task, error)

	// ListByProjectID 分页查询指定项目下的所有任务。
	// 返回结果按创建时间倒序排列。
	ListByProjectID(projectID uint64, page, pageSize int) ([]models.Task, int64, error)

	// ListByUserID 分页查询用户有权访问的所有任务。
	// 返回结果和总数。
	ListByUserID(userID uint64, offset, limit int) ([]models.Task, int64, error)

	// ListCompletedByUserIDAndTimeRange 查询用户在指定时间范围内完成的任务。
	// 参数 startTime, endTime: 时间范围（包含边界），基于 completed_at 字段。
	ListCompletedByUserIDAndTimeRange(userID uint64, startTime, endTime time.Time) ([]models.Task, error)

	// Update 更新任务信息（全量更新）。
	// 错误：如果任务不存在返回 ErrNotFound。
	Update(task *models.Task) error

	// UpdateStatus 更新任务状态。
	UpdateStatus(id uint64, status uint8) error

	// Delete 删除任务。
	// 错误：如果任务不存在返回 ErrNotFound。
	Delete(id uint64) error

	// GetByStatus 查询指定项目下特定状态的任务列表。
	GetByStatus(projectID uint64, status uint8) ([]models.Task, error)

	// GetByDeadlineBefore 查询指定项目下截止时间早于给定时间的任务（按截止时间升序）。
	// 支持分页，若 projectID == 0 则忽略项目条件。
	GetByDeadlineBefore(projectID uint64, deadline time.Time, page, pageSize int) ([]models.Task, int64, error)

	// GetByDeadlineAfter 查询指定项目下截止时间晚于给定时间的任务（按截止时间升序）。
	GetByDeadlineAfter(projectID uint64, deadline time.Time, page, pageSize int) ([]models.Task, int64, error)

	// GetByTimePeriod 查询指定项目下截止时间在区间内的任务（按截止时间升序）。
	GetByTimePeriod(projectID uint64, start, end time.Time, page, pageSize int) ([]models.Task, int64, error)

	// SetPrerequisiteTask 为任务设置前置任务。
	// 返回值 dependency: 创建的任务依赖关系对象。
	// 错误：若任一任务不存在返回 ErrNotFound；若两个任务不属于同一项目返回 ErrPermissionDenied；
	//       若依赖关系已存在返回 ErrRecordExist。
	SetPrerequisiteTask(prerequisiteID, taskID uint64) (dependency *models.TaskDependency, err error)

	// UnsetPrerequisiteTask 移除任务的前置任务关系。
	// 错误：若依赖关系不存在返回 ErrNotFound。
	UnsetPrerequisiteTask(prerequisiteID, taskID uint64) (err error)

	// GetPrerequisites 获取任务的所有前置任务（依赖关系）。
	GetPrerequisites(taskID uint64) (prerequisites []models.TaskDependency, err error)

	// GetPostrequisites 获取依赖该任务的所有后续任务。
	GetPostrequisites(prerequisiteID uint64) (prerequisites []models.TaskDependency, err error)

	// ListByAccountID 查询指定账户在时间范围内的所有任务。
	// 参数 startTime, endTime: 基于任务的完成时间（completed_at）。
	// 实现行为：通过 task_payments 表找到账户关联的任务，再按完成时间过滤。
	ListByAccountID(accountID uint64, startTime, endTime time.Time) ([]models.Task, error)

	// CheckOwnership 检查用户是否对任务有访问权限（通过任务所在项目与用户关联）。
	CheckOwnership(userID, taskID uint64) (bool, error)
}
