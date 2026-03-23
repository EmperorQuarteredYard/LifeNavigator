package repository

import (
	"LifeNavigator/internal/interfaces/Repository"
	"LifeNavigator/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// NewTaskPaymentRepository 创建一个 TaskBudgetRepository 实例
func NewTaskPaymentRepository(db *gorm.DB) Repository.TaskBudgetRepository {
	return &taskBudgetRepository{baseRepository: &baseRepository{db: db}}
}

type taskBudgetRepository struct {
	*baseRepository
}

func (r *taskBudgetRepository) ListByAccountID(accountID uint64, startTime, endTime time.Time) (payments []models.TaskPayment, err error) {
	var tasks []models.Task
	var tempRes []models.TaskPayment
	err = r.db.Model(models.Task{}).Where("completed_at < ? and completed_at >= ?", endTime, startTime).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		err = r.db.Model(models.TaskPayment{}).Where("task_id = ? and account_id = ?", task.ID).Find(&tempRes).Error
		if err != nil {
			return nil, err
		}
		payments = append(payments, models.TaskPayment{})
	}
	return payments, nil
}

func (r *taskBudgetRepository) Create(budget *models.TaskPayment) error {

	return r.create(budget)
}

func (r *taskBudgetRepository) GetByID(id uint64) (*models.TaskPayment, error) {
	var budget models.TaskPayment
	result := r.db.First(&budget, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, Repository.ErrNotFound
		}
		return nil, result.Error
	}
	return &budget, nil
}

func (r *taskBudgetRepository) GetByTaskID(taskID uint64) ([]models.TaskPayment, error) {
	var budgets []models.TaskPayment
	result := r.db.Where("task_id = ?", taskID).Find(&budgets)
	return budgets, result.Error
}

func (r *taskBudgetRepository) Update(budget *models.TaskPayment) error {
	result := r.db.Save(budget)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return Repository.ErrNotFound
	}
	return nil
}

func (r *taskBudgetRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.TaskPayment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return Repository.ErrNotFound
	}
	return nil
}

func (r *taskBudgetRepository) DeleteByTaskID(taskID uint64) error {
	result := r.db.Where("task_id = ?", taskID).Delete(&models.TaskPayment{})
	return result.Error
}
