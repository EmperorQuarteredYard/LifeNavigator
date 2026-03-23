package Repository

import (
	"LifeNavigator/internal/models"
)

// ProjectBudgetRepository 定义了项目预算的数据访问接口。
type ProjectBudgetRepository interface {
	// Create 创建一个新的项目预算记录。
	Create(budget *models.ProjectBudget) error

	// GetByID 根据主键 ID 查询项目预算。
	// 错误：如果记录不存在返回 ErrNotFound。
	GetByID(id uint64) (*models.ProjectBudget, error)

	// GetByProjectID 根据项目 ID 查询该项目下的所有预算项。
	GetByProjectID(projectID uint64) ([]models.ProjectBudget, error)

	// Update 更新现有的项目预算记录。
	// 错误：如果记录不存在返回 ErrNotFound。
	Update(budget *models.ProjectBudget) error

	// Delete 根据 ID 删除项目预算记录。
	// 错误：如果记录不存在返回 ErrNotFound。
	Delete(id uint64) error

	// DeleteByProjectID 删除指定项目的所有预算记录。
	DeleteByProjectID(projectID uint64) error

	// AddUsed 原子增加预算已用金额（不检查是否超出预算）。
	AddUsed(budgetID uint64, amount float64) error

	// SubtractUsed 原子减少预算已用金额（不检查负数）。
	SubtractUsed(budgetID uint64, amount float64) error

	// GetByAccountID 查询指定账户关联的所有项目预算。
	GetByAccountID(accountID uint64) ([]models.ProjectBudget, error)

	// UpdateAccountID 更新指定预算的账户 ID。
	UpdateAccountID(budgetID uint64, accountID uint64) error

	// SetUsedZero 将指定预算的 used 字段置为 0。
	// 注意：此方法没有使用事务，调用方需确保一致性。
	SetUsedZero(budgetID, projectID uint64) error
}
