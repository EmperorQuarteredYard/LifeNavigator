package repository

import (
	"LifeNavigator/internal/interfaces/Repository"
	"LifeNavigator/internal/models"
	"errors"

	"gorm.io/gorm"
)

// NewProjectBudgetRepository 创建一个 ProjectBudgetRepository 实例
func NewProjectBudgetRepository(db *gorm.DB) Repository.ProjectBudgetRepository {
	return &projectBudgetRepository{baseRepository: &baseRepository{db: db}}
}

type projectBudgetRepository struct {
	*baseRepository
}

func (r *projectBudgetRepository) SetUsedZero(budgetID, projectID uint64) error {
	err := r.db.Model(&models.ProjectBudget{}).Where("project_id = ? and id = ?", projectID, budgetID).Update("used", 0).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Repository.ErrNotFound
		}
		return err
	}
	return nil
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
	return r.create(budget)
}

func (r *projectBudgetRepository) GetByID(id uint64) (*models.ProjectBudget, error) {
	var budget models.ProjectBudget
	result := r.db.First(&budget, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, Repository.ErrNotFound
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
		return Repository.ErrNotFound
	}
	return nil
}

func (r *projectBudgetRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.ProjectBudget{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return Repository.ErrNotFound
	}
	return nil
}

func (r *projectBudgetRepository) DeleteByProjectID(projectID uint64) error {
	result := r.db.Where("project_id = ?", projectID).Delete(&models.ProjectBudget{})
	return result.Error
}
