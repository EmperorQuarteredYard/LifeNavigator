package Service

import (
	"time"

	"LifeNavigator/internal/models"
	"LifeNavigator/pkg/dto"
)

// TaskService 定义任务管理接口
type TaskService interface {
	// Create 创建任务
	// 参数:
	//   userID: 当前用户ID，用于权限校验（需要项目编辑权限）
	//   task: 任务模型（应包含项目ID、名称等）
	// 返回值:
	//   *dto.TaskResponse: 创建成功的任务详情
	//   error: 可能返回 ErrForbidden、ErrInternal
	Create(userID uint64, task *models.Task) (*dto.TaskResponse, error)

	// GetByID 获取任务详情
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   id: 任务ID
	// 返回值:
	//   *dto.TaskResponse: 任务详情（包含支付记录）
	//   error: 可能返回 ErrForbidden、ErrTaskNotFound、ErrInternal
	GetByID(userID, id uint64) (*dto.TaskResponse, error)

	// ListByProjectID 获取项目下的任务列表（分页）
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   projectID: 项目ID
	//   page: 页码
	//   pageSize: 每页数量
	// 返回值:
	//   *dto.TaskListResponse: 任务列表
	//   error: 可能返回 ErrForbidden、ErrInvalidInput、ErrInternal
	ListByProjectID(userID, projectID uint64, page, pageSize int) (*dto.TaskListResponse, error)

	// Update 更新任务信息
	// 参数:
	//   userID: 当前用户ID，用于权限校验（需要项目编辑权限）
	//   task: 更新后的任务模型（必须包含ID）
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrTaskNotFound、ErrInternal
	Update(userID uint64, task *models.Task) error

	// UpdateStatus 更新任务状态
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   id: 任务ID
	//   status: 新状态（0-未开始，1-进行中，2-已完成等）
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrTaskNotFound、ErrInternal
	UpdateStatus(userID, id uint64, status uint8) error

	// Delete 删除任务
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   id: 任务ID
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrTaskNotFound、ErrInternal
	Delete(userID, id uint64) error

	// GetByStatus 获取项目下指定状态的任务列表
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   projectID: 项目ID
	//   status: 任务状态
	// 返回值:
	//   []*dto.TaskResponse: 任务列表
	//   error: 可能返回 ErrForbidden、ErrInternal
	GetByStatus(userID, projectID uint64, status uint8) ([]*dto.TaskResponse, error)

	// GetByDeadlineBefore 获取截止时间在指定时间之前的任务
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   projectID: 项目ID
	//   deadline: 截止时间
	//   page: 页码
	//   pageSize: 每页数量
	// 返回值:
	//   *dto.TaskListResponse: 任务列表
	//   error: 可能返回 ErrForbidden、ErrInvalidInput、ErrInternal
	GetByDeadlineBefore(userID, projectID uint64, deadline time.Time, page, pageSize int) (*dto.TaskListResponse, error)

	// GetByDeadlineAfter 获取截止时间在指定时间之后的任务
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   projectID: 项目ID
	//   deadline: 截止时间
	//   page: 页码
	//   pageSize: 每页数量
	// 返回值:
	//   *dto.TaskListResponse: 任务列表
	//   error: 可能返回 ErrForbidden、ErrInvalidInput、ErrInternal
	GetByDeadlineAfter(userID, projectID uint64, deadline time.Time, page, pageSize int) (*dto.TaskListResponse, error)

	// GetByTimePeriod 获取创建时间在指定时间段内的任务
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   projectID: 项目ID
	//   start: 开始时间
	//   end: 结束时间
	//   page: 页码
	//   pageSize: 每页数量
	// 返回值:
	//   *dto.TaskListResponse: 任务列表
	//   error: 可能返回 ErrForbidden、ErrInvalidInput、ErrInternal
	GetByTimePeriod(userID, projectID uint64, start, end time.Time, page, pageSize int) (*dto.TaskListResponse, error)

	// SetPrerequisiteTask 设置任务依赖关系（前置任务）
	// 参数:
	//   userID: 当前用户ID，用于权限校验（需要项目编辑权限）
	//   prerequisiteID: 前置任务ID
	//   taskID: 后置任务ID
	// 返回值:
	//   *dto.DependencyResponse: 创建的依赖关系
	//   error: 可能返回 ErrForbidden、ErrTaskNotFound、ErrInternal
	SetPrerequisiteTask(userID, prerequisiteID, taskID uint64) (*dto.DependencyResponse, error)

	// UnsetPrerequisiteTask 移除任务依赖关系
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   prerequisiteID: 前置任务ID
	//   taskID: 后置任务ID
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrTaskDependencyNotFound、ErrInternal
	UnsetPrerequisiteTask(userID, prerequisiteID, taskID uint64) error

	// GetPrerequisites 获取任务的所有前置任务
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   taskID: 任务ID
	// 返回值:
	//   []*dto.DependencyResponse: 前置任务列表
	//   error: 可能返回 ErrForbidden、ErrTaskDependencyNotFound、ErrInternal
	GetPrerequisites(userID, taskID uint64) ([]*dto.DependencyResponse, error)

	// GetPostrequisite 获取任务的所有后置任务
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   taskID: 任务ID
	// 返回值:
	//   []*dto.DependencyResponse: 后置任务列表
	//   error: 可能返回 ErrForbidden、ErrTaskDependencyNotFound、ErrInternal
	GetPostrequisite(userID, taskID uint64) ([]*dto.DependencyResponse, error)
}
