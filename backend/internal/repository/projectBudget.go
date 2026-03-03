package repository

import (
	"LifeNavigator/backend/internal/models"
	"errors"

	"gorm.io/gorm"
)

// ProjectBudgetRepository 定义了项目预算的数据访问接口
type ProjectBudgetRepository interface {
	// Create 创建一个新的项目预算记录
	Create(budget *models.ProjectBudget) error

	// GetByID 根据主键 ID 查询项目预算
	GetByID(id uint64) (*models.ProjectBudget, error)

	// GetByProjectID 根据项目 ID 查询该项目的所有预算项
	GetByProjectID(projectID uint64) ([]models.ProjectBudget, error)

	// Update 更新现有的项目预算记录
	Update(budget *models.ProjectBudget) error

	// Delete 根据 ID 删除项目预算记录
	Delete(id uint64) error

	// DeleteByProjectID 删除指定项目的所有预算记录
	DeleteByProjectID(projectID uint64) error
}

type projectBudgetRepository struct {
	db *gorm.DB
}

// NewProjectBudgetRepository 创建一个 ProjectBudgetRepository 实例
func NewProjectBudgetRepository(db *gorm.DB) ProjectBudgetRepository {
	return &projectBudgetRepository{db: db}
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
