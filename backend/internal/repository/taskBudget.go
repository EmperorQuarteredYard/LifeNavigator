package repository

import (
	"LifeNavigator/internal/models"
	"errors"

	"gorm.io/gorm"
)

// TaskBudgetRepository 定义了任务预算的数据访问接口
type TaskBudgetRepository interface {
	// Create 创建一个新的任务预算记录
	// 成功返回 nil，失败返回错误（如 ErrDuplicate 或其他数据库错误）
	Create(budget *models.TaskBudget) error

	// GetByID 根据主键 ID 查询任务预算
	// 返回预算对象指针，若记录不存在返回 ErrNotFound
	GetByID(id uint64) (*models.TaskBudget, error)

	// GetByTaskID 根据任务 ID 查询该任务的所有预算项
	// 返回预算切片，不会返回 nil（若无记录返回空切片）
	GetByTaskID(taskID uint64) ([]models.TaskBudget, error)

	// Update 更新现有的任务预算记录
	// 根据预算对象的 ID 进行更新，若记录不存在返回 ErrNotFound
	Update(budget *models.TaskBudget) error

	// Delete 根据 ID 删除任务预算记录
	// 若记录不存在返回 ErrNotFound
	Delete(id uint64) error

	// DeleteByTaskID 删除指定任务的所有预算记录（通常用于级联删除）
	// 即使任务没有预算也不返回错误
	DeleteByTaskID(taskID uint64) error
}

type taskBudgetRepository struct {
	db *gorm.DB
}

// NewTaskBudgetRepository 创建一个 TaskBudgetRepository 实例
func NewTaskBudgetRepository(db *gorm.DB) TaskBudgetRepository {
	return &taskBudgetRepository{db: db}
}

func (r *taskBudgetRepository) Create(budget *models.TaskBudget) error {
	result := r.db.Create(budget)
	if result.Error != nil {
		// 可根据错误类型封装自定义错误，例如唯一索引冲突
		return result.Error
	}
	return nil
}

func (r *taskBudgetRepository) GetByID(id uint64) (*models.TaskBudget, error) {
	var budget models.TaskBudget
	result := r.db.First(&budget, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &budget, nil
}

func (r *taskBudgetRepository) GetByTaskID(taskID uint64) ([]models.TaskBudget, error) {
	var budgets []models.TaskBudget
	result := r.db.Where("task_id = ?", taskID).Find(&budgets)
	return budgets, result.Error
}

func (r *taskBudgetRepository) Update(budget *models.TaskBudget) error {
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
	result := r.db.Delete(&models.TaskBudget{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *taskBudgetRepository) DeleteByTaskID(taskID uint64) error {
	result := r.db.Where("task_id = ?", taskID).Delete(&models.TaskBudget{})
	return result.Error
}
