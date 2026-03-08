package repository

import (
	"LifeNavigator/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id uint64) (*models.Task, error)
	GetByUserID(userID, taskID uint64) (*models.Task, error)
	// ListByProjectID 分页查询指定项目的任务，返回任务列表和总数。
	//   - projectID: 项目ID
	//   - page:      页码，从0开始
	//   - pageSize:  每页大小。若 pageSize < 0 返回 ErrInvalidInput；
	//                若 pageSize == 0 则返回所有记录（不分页），总数仍返回。
	ListByProjectID(projectID uint64, page, pageSize int) ([]models.Task, int64, error)
	ListByUserID(userID uint64, offset, limit int) ([]models.Task, int64, error)
	// ListCompletedByUserIDAndTimeRange 查询用户在指定时间区间内已完成的任务
	//   - userID:    用户ID
	//   - startTime: 开始时间
	//   - endTime:   结束时间
	//   返回CompletedAt在[startTime, endTime]范围内且不为nil的任务
	ListCompletedByUserIDAndTimeRange(userID uint64, startTime, endTime time.Time) ([]models.Task, error)
	Update(task *models.Task) error
	Delete(id uint64) error
	GetByStatus(projectID uint64, status uint8) ([]models.Task, error)
	// GetByDeadlineBefore 查询指定项目下截止时间早于给定时间戳的任务，支持分页并返回总数。
	//   - projectID: 项目ID，若为0则忽略项目过滤（查询所有项目）
	//   - pageSize:  每页记录数（<=0 时返回所有记录，不进行分页，此时 page 参数无效）
	//   - page:      页码，从0开始
	// 返回值：
	//   - tasks: 任务切片
	//   - total: 符合条件的总记录数
	//   - err:   错误信息
	GetByDeadlineBefore(projectID uint64, deadline time.Time, pageSize, page int) ([]models.Task, int64, error)
	GetByDeadlineAfter(projectID uint64, deadline time.Time, page, pageSize int) ([]models.Task, int64, error)

	// GetByTimePeriod 查询指定项目下截止时间在 [start, end] 范围内的任务，支持分页并返回总数。
	//   - projectID: 项目ID，若为0则忽略项目过滤
	//   - page:      页码（从0开始）
	//   - pageSize:  每页大小，负数返回 ErrInvalidInput，0表示返回所有记录
	GetByTimePeriod(projectID uint64, start, end time.Time, page, pageSize int) ([]models.Task, int64, error)
	//SetPrerequisiteTask 不允许跨任务、跨用户创建，会检查ID是否有效
	//返回错误：
	// - ErrNotFound:ID无效
	// - ErrPermissionDenied 跨任务或跨用户
	SetPrerequisiteTask(prerequisiteID, taskID uint64) (dependency *models.TaskDependency, err error)
	//UnsetPrerequisiteTask 不会检查ID是否有效、用户是否正确
	UnsetPrerequisiteTask(prerequisiteID, taskID uint64) (err error)
	GetPrerequisites(taskID uint64) (prerequisites []models.TaskDependency, err error)
	GetPostrequisites(prerequisiteID uint64) (prerequisites []models.TaskDependency, err error)
	ListByAccountID(userID, accountID uint64, startTime, endTime time.Time) ([]models.Task, error)
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

type taskRepository struct {
	db *gorm.DB
}

func (r *taskRepository) ListByAccountID(userID, accountID uint64, startTime, endTime time.Time) ([]models.Task, error) {
	var tasks, result []models.Task
	var payments []models.TaskPayment
	err := r.db.Where("record_time < ? AND record_time > ? and user_id = ?", endTime, startTime, userID).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		if r.db.Where("task_id = ? and account_id", task.ID, accountID).Find(&payments).RowsAffected != 0 {
			result = append(result, task)
		}
	}
	return result, nil
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
	err = r.db.Where("prerequisite_id = ? and task_id ?", prerequisite, task).Delete(&models.TaskDependency{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return
	}
	return
}

func (r *taskRepository) SetPrerequisiteTask(prerequisiteID, taskID uint64) (dependency *models.TaskDependency, err error) {
	task := &models.Task{}
	prerequisite := &models.Task{}

	if err = r.db.Where("id = ?", taskID).First(task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if err = r.db.Where("id = ?", prerequisiteID).First(prerequisite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if prerequisite.UserID != task.UserID || prerequisite.ProjectID != task.ProjectID {
		return nil, ErrPermissionDenied
	}

	dependency = &models.TaskDependency{
		TaskID:         task.ID,
		ProjectID:      task.ProjectID,
		PrerequisiteID: prerequisiteID,
		UserID:         task.UserID,
	}
	if err := r.db.Where("task_id = ? and prerequisite_id = ?", taskID, prerequisiteID).First(dependency).Error; err == nil {
		return dependency, ErrRecordExist
	}
	err = r.db.Create(dependency).Error
	if err != nil {
		return nil, err
	}
	return dependency, nil
}

func (r *taskRepository) GetByUserID(userID, taskID uint64) (*models.Task, error) {
	task := &models.Task{}
	err := r.db.Where("user_id = ? AND task_id = ?", userID, taskID).First(task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return task, nil
}

func (r *taskRepository) GetByDeadlineAfter(projectID uint64, deadline time.Time, page, pageSize int) ([]models.Task, int64, error) {
	if pageSize < 0 {
		return nil, 0, ErrInvalidInput
	}

	var tasks []models.Task
	var total int64

	// 构建基础查询
	query := r.db.Model(&models.Task{})
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	query = query.Where("deadline > ?", deadline) // 严格大于

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序和分页（按截止时间升序，越早的越靠前；也可按需改为降序）
	query = query.Order("deadline ASC")
	if pageSize > 0 {
		offset := page * pageSize
		if offset < 0 {
			offset = 0
		}
		query = query.Offset(offset).Limit(pageSize)
	}

	// 执行查询
	if err := query.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *taskRepository) Create(task *models.Task) error {
	result := r.db.Create(task)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *taskRepository) GetByID(id uint64) (*models.Task, error) {
	var task models.Task
	result := r.db.First(&task, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &task, nil
}

func (r *taskRepository) ListByProjectID(projectID uint64, page, pageSize int) ([]models.Task, int64, error) {
	if pageSize < 0 {
		return nil, 0, ErrInvalidInput
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
	if err := r.db.Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

func (r *taskRepository) ListCompletedByUserIDAndTimeRange(userID uint64, startTime, endTime time.Time) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("user_id = ? AND completed_at IS NOT NULL AND completed_at >= ? AND completed_at <= ?", userID, startTime, endTime).
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
		return ErrNotFound
	}
	return nil
}

func (r *taskRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
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
		return nil, 0, ErrInvalidInput
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
		return nil, 0, ErrInvalidInput
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

	query = query.Order("deadline ASC") // 按截止时间升序
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
