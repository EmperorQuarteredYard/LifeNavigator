package Service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/pkg/dto"
)

// BudgetService 定义预算相关的业务接口
// 包括项目预算管理和任务支付管理
type BudgetService interface {
	// AddBudget 为项目添加预算项
	// 参数:
	//   userID: 当前用户ID，用于权限校验（需要项目编辑权限）
	//   projectID: 项目ID
	//   budget: 预算项模型，应包含单位、预算金额、关联账户ID等
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrProjectNotFound、ErrAccountNotFound、ErrInternal
	AddBudget(userID, projectID uint64, budget *models.ProjectBudget) error

	// UpdateBudget 更新项目预算项
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   budget: 更新后的预算项模型（必须包含ID）
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrBudgetNotFound、ErrAccountNotFound、ErrInternal
	UpdateBudget(userID uint64, budget *models.ProjectBudget) error

	// DeleteBudget 删除项目预算项
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   budgetID: 预算项ID
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrBudgetNotFound、ErrInternal
	DeleteBudget(userID, budgetID uint64) error

	// RefreshBudget 刷新预算（将预算项中已使用的金额回滚到账户净余额）
	// 通常在项目结束或手动触发时调用。
	// 参数:
	//   projectID: 项目ID
	// 返回值:
	//   error: 可能返回 ErrProjectNotFound、ErrAccountNotFound、ErrInternal
	RefreshBudget(projectID uint64) error

	// AddPayment 为任务添加支付记录
	// 参数:
	//   userID: 当前用户ID，用于权限校验（需要项目编辑权限）
	//   taskID: 任务ID
	//   payment: 支付项模型，包含预算ID和支付金额
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrTaskNotFound、ErrBudgetNotFound、ErrAccountNotFound、ErrInternal
	AddPayment(userID, taskID uint64, payment *models.TaskPayment) error

	// UpdatePayment 更新任务支付记录
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   payment: 更新后的支付项模型（必须包含ID）
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrBudgetNotFound、ErrAccountNotFound、ErrInternal
	UpdatePayment(userID uint64, payment *models.TaskPayment) error

	// DeletePayment 删除任务支付记录
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   paymentID: 支付项ID
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrBudgetNotFound、ErrAccountNotFound、ErrInternal
	DeletePayment(userID, paymentID uint64) error

	// GetPaymentByTaskID 获取指定任务的所有支付记录
	// 参数:
	//   userID: 当前用户ID，用于权限校验（需要项目只读权限）
	//   taskID: 任务ID
	// 返回值:
	//   []*dto.TaskPaymentResponse: 支付记录列表
	//   error: 可能返回 ErrForbidden、ErrTaskNotFound、ErrInternal
	GetPaymentByTaskID(userID, taskID uint64) ([]*dto.TaskPaymentResponse, error)

	// StartAutoRefresh 启动自动刷新预算的调度器
	// 根据项目的刷新间隔和上次刷新时间，自动定时执行 RefreshBudget
	// 返回值:
	//   error: 启动失败时返回 ErrInternal
	StartAutoRefresh() error

	// EndAutoRefresh 停止自动刷新预算的调度器
	EndAutoRefresh()
}
