package repository

import (
	"LifeNavigator/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

// TaskBudgetRepository 定义了任务预算的数据访问接口
type TaskBudgetRepository interface {
	Create(budget *models.TaskPayment) error                 // 创建一个新的任务预算记录,成功返回 nil，失败返回错误（如 ErrDuplicate 或其他数据库错误）
	GetByID(id uint64) (*models.TaskPayment, error)          // 根据主键 ID 查询任务预算 返回预算对象指针，若记录不存在返回 ErrNotFound
	GetByTaskID(taskID uint64) ([]models.TaskPayment, error) // GetByTaskID 根据任务 ID 查询该任务的所有预算项,返回预算切片，不会返回 nil（若无记录返回空切片）
	Update(budget *models.TaskPayment) error                 // Update 更新现有的任务预算记录,根据预算对象的 ID 进行更新，若记录不存在返回 ErrNotFound
	Delete(id uint64) error                                  // Delete 根据 ID 删除任务预算记录,若记录不存在返回 ErrNotFound
	DeleteByTaskID(taskID uint64) error                      //  删除指定任务ime.Time) ([]models.TaskPayment, error)
	ListByAccountID(accountID uint64, startTime, endTime time.Time) ([]models.TaskPayment, error)
}

// NewTaskPaymentRepository 创建一个 TaskBudgetRepository 实例
func NewTaskPaymentRepository(db *gorm.DB) TaskBudgetRepository {
	return &taskBudgetRepository{db: db}
}

type taskBudgetRepository struct {
	db *gorm.DB
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
	result := r.db.Create(budget)
	if result.Error != nil {
		// 可根据错误类型封装自定义错误，例如唯一索引冲突
		return result.Error
	}
	return nil
}

func (r *taskBudgetRepository) GetByID(id uint64) (*models.TaskPayment, error) {
	var budget models.TaskPayment
	result := r.db.First(&budget, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
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
		return ErrNotFound
	}
	return nil
}

func (r *taskBudgetRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.TaskPayment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *taskBudgetRepository) DeleteByTaskID(taskID uint64) error {
	result := r.db.Where("task_id = ?", taskID).Delete(&models.TaskPayment{})
	return result.Error
}
