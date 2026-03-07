package repository

import (
	"LifeNavigator/internal/models"
	"errors"

	"gorm.io/gorm"
)

// ProjectBudgetRepository 定义了项目预算的数据访问接口
type ProjectBudgetRepository interface {
	Create(budget *models.ProjectBudget) error                       // 创建一个新的项目预算记录
	GetByID(id uint64) (*models.ProjectBudget, error)                // 根据主键 ID 查询项目预算
	GetByProjectID(projectID uint64) ([]models.ProjectBudget, error) // 根据项目 ID 查询该项目的所有预算项
	Update(budget *models.ProjectBudget) error                       // 更新现有的项目预算记录
	Delete(id uint64) error                                          // 根据 ID 删除项目预算记录
	DeleteByProjectID(projectID uint64) error                        // 删除指定项目的所有预算记录
	AddUsed(budgetID uint64, amount float64) error                   // 原子增加预算已用金额（不检查超额）
	SubtractUsed(budgetID uint64, amount float64) error              // 原子减少预算已用金额（不检查负数）
	GetByAccountID(accountID uint64) ([]models.ProjectBudget, error) // 查询指定账户关联的所有项目预算
	UpdateAccountID(budgetID uint64, accountID uint64) error         // 更新指定预算的账户 ID
}

// NewProjectBudgetRepository 创建一个 ProjectBudgetRepository 实例
func NewProjectBudgetRepository(db *gorm.DB) ProjectBudgetRepository {
	return &projectBudgetRepository{db: db}
}

type projectBudgetRepository struct {
	db *gorm.DB
}

func (r *projectBudgetRepository) AddUsed(budgetID uint64, amount float64) error {
	return r.db.Model(&models.ProjectBudget{}).
		Where("id = ?", budgetID).
		Update("used", gorm.Expr("used + ?", amount)).Error
}

func (r *projectBudgetRepository) SubtractUsed(budgetID uint64, amount float64) error {
	return r.db.Model(&models.ProjectBudget{}).
		Where("id = ?", budgetID).
		Update("used", gorm.Expr("used - ?", amount)).Error
}

func (r *projectBudgetRepository) GetByAccountID(accountID uint64) ([]models.ProjectBudget, error) {
	var budgets []models.ProjectBudget
	err := r.db.Where("account_id = ?", accountID).Find(&budgets).Error
	return budgets, err
}

func (r *projectBudgetRepository) UpdateAccountID(budgetID uint64, accountID uint64) error {
	return r.db.Model(&models.ProjectBudget{}).
		Where("id = ?", budgetID).
		Update("account_id", accountID).Error
}
func (r *projectBudgetRepository) Create(budget *models.ProjectBudget) error {
	result := r.db.Create(budget)
	return result.Error
}

func (r *projectBudgetRepository) GetByID(id uint64) (*models.ProjectBudget, error) {
	var budget models.ProjectBudget
	result := r.db.First(&budget, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &budget, nil
}

func (r *projectBudgetRepository) GetByProjectID(projectID uint64) ([]models.ProjectBudget, error) {
	var budgets []models.ProjectBudget
	result := r.db.Where("project_id = ?", projectID).Find(&budgets)
	return budgets, result.Error
}

func (r *projectBudgetRepository) Update(budget *models.ProjectBudget) error {
	result := r.db.Save(budget)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *projectBudgetRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.ProjectBudget{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *projectBudgetRepository) DeleteByProjectID(projectID uint64) error {
	result := r.db.Where("project_id = ?", projectID).Delete(&models.ProjectBudget{})
	return result.Error
}
