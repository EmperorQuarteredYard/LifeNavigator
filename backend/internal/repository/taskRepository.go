package repository

import (
	"LifeNavigator/internal/interfaces/Repository"
	"LifeNavigator/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

func NewTaskRepository(db *gorm.DB) Repository.TaskRepository {
	return &taskRepository{baseRepository: &baseRepository{db: db}}
}

type taskRepository struct {
	*baseRepository
}

func (r *taskRepository) CheckOwnership(userID, taskID uint64) (bool, error) {
	var task models.Task
	if err := r.db.Select("project_id").First(&task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	var count int64
	err := r.db.Table("project_users").
		Where("user_id = ? AND project_id = ?", userID, task.ProjectID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *taskRepository) ListByAccountID(accountID uint64, startTime, endTime time.Time) ([]models.Task, error) {
	var tasks []models.Task
	var taskIDs []uint64

	err := r.db.Table("task_payments").
		Select("DISTINCT task_id").
		Where("account_id = ?", accountID).
		Pluck("task_id", &taskIDs).Error
	if err != nil {
		return nil, err
	}

	if len(taskIDs) == 0 {
		return tasks, nil
	}

	err = r.db.Where("id IN (?) AND completed_at >= ? AND completed_at <= ?", taskIDs, startTime, endTime).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) GetPostrequisites(prerequisiteID uint64) (prerequisites []models.TaskDependency, err error) {
	err = r.db.Where("prerequisite_id = ?", prerequisiteID).Find(&prerequisites).Error
	if err != nil {
		return nil, err
	}
	return prerequisites, err
}

func (r *taskRepository) GetPrerequisites(taskID uint64) (prerequisites []models.TaskDependency, err error) {
	err = r.db.Where("task_id = ?", taskID).Find(&prerequisites).Error
	if err != nil {
		return nil, err
	}
	return prerequisites, nil
}

func (r *taskRepository) UnsetPrerequisiteTask(prerequisite, task uint64) (err error) {
	err = r.db.Where("prerequisite_id = ? AND task_id = ?", prerequisite, task).Delete(&models.TaskDependency{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Repository.ErrNotFound
		}
		return
	}
	return
}

func (r *taskRepository) SetPrerequisiteTask(prerequisiteID, taskID uint64) (dependency *models.TaskDependency, err error) {
	task := &models.Task{}
	prerequisite := &models.Task{}

	if err = r.db.First(task, taskID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, Repository.ErrNotFound
		}
		return nil, err
	}
	if err = r.db.First(prerequisite, prerequisiteID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, Repository.ErrNotFound
		}
		return nil, err
	}
	if prerequisite.ProjectID != task.ProjectID {
		return nil, Repository.ErrPermissionDenied
	}

	dependency = &models.TaskDependency{
		TaskID:         task.ID,
		ProjectID:      task.ProjectID,
		PrerequisiteID: prerequisiteID,
	}
	if err := r.db.Where("task_id = ? AND prerequisite_id = ?", taskID, prerequisiteID).First(dependency).Error; err == nil {
		return dependency, Repository.ErrRecordExist
	}
	err = r.db.Create(dependency).Error
	if err != nil {
		return nil, err
	}
	return dependency, nil
}

func (r *taskRepository) GetByDeadlineAfter(projectID uint64, deadline time.Time, page, pageSize int) ([]models.Task, int64, error) {
	if pageSize < 0 {
		return nil, 0, Repository.ErrInvalidInput
	}

	var tasks []models.Task
	var total int64

	query := r.db.Model(&models.Task{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	query = query.Where("deadline > ?", deadline)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("deadline ASC")
	if pageSize > 0 {
		offset := page * pageSize
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) Create(task *models.Task) error {
	return r.create(task)
}

func (r *taskRepository) GetByID(id uint64) (*models.Task, error) {
	var task models.Task
	result := r.db.First(&task, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, Repository.ErrNotFound
		}
		return nil, result.Error
	}
	return &task, nil
}

func (r *taskRepository) ListByProjectID(projectID uint64, page, pageSize int) ([]models.Task, int64, error) {
	if pageSize < 0 {
		return nil, 0, Repository.ErrInvalidInput
	}

	var tasks []models.Task
	var total int64

	query := r.db.Model(&models.Task{}).Where("project_id = ?", projectID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC")
	if pageSize > 0 {
		offset := page * pageSize
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(pageSize)
	}

	result := query.Find(&tasks)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return tasks, total, nil
}

func (r *taskRepository) ListByUserID(userID uint64, offset, limit int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	subQuery := r.db.Table("project_users").
		Select("project_id").
		Where("user_id = ?", userID)

	if err := r.db.Model(&models.Task{}).Where("project_id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("project_id IN (?)", subQuery).Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (r *taskRepository) ListCompletedByUserIDAndTimeRange(userID uint64, startTime, endTime time.Time) ([]models.Task, error) {
	var tasks []models.Task

	subQuery := r.db.Table("project_users").
		Select("project_id").
		Where("user_id = ?", userID)

	err := r.db.Where("project_id IN (?) AND completed_at IS NOT NULL AND completed_at >= ? AND completed_at <= ?", subQuery, startTime, endTime).
		Order("completed_at DESC").
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) Update(task *models.Task) error {
	result := r.db.Save(task)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return Repository.ErrNotFound
	}
	return nil
}

func (r *taskRepository) UpdateStatus(id uint64, status uint8) error {
	result := r.db.Model(&models.Task{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return Repository.ErrNotFound
	}
	return nil
}

func (r *taskRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return Repository.ErrNotFound
	}
	return nil
}

func (r *taskRepository) GetByStatus(projectID uint64, status uint8) ([]models.Task, error) {
	var tasks []models.Task
	result := r.db.Where("project_id = ? AND status = ?", projectID, status).Find(&tasks)
	return tasks, result.Error
}

func (r *taskRepository) GetByDeadlineBefore(projectID uint64, deadline time.Time, page, pageSize int) ([]models.Task, int64, error) {
	if pageSize < 0 {
		return nil, 0, Repository.ErrInvalidInput
	}

	var tasks []models.Task
	var total int64
	query := r.db.Model(&models.Task{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	query = query.Where("deadline < ?", deadline)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	query = query.Order("deadline ASC")
	if pageSize > 0 {
		offset := page * pageSize
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(pageSize)
	}
	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) GetByTimePeriod(projectID uint64, start, end time.Time, page, pageSize int) ([]models.Task, int64, error) {
	if pageSize < 0 {
		return nil, 0, Repository.ErrInvalidInput
	}

	var tasks []models.Task
	var total int64

	query := r.db.Model(&models.Task{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	query = query.Where("deadline >= ? AND deadline <= ?", start, end)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("deadline ASC")
	if pageSize > 0 {
		offset := page * pageSize
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(pageSize)
	}

	result := query.Find(&tasks)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return tasks, total, nil
}
