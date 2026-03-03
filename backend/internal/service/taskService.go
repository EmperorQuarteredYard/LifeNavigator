package service

import (
	"LifeNavigator/backend/internal/models"
	"LifeNavigator/backend/internal/repository"
	"context"
	"errors"
	"log"
	"time"
)

type TaskService interface {
	Create(task *models.Task, currentUserID uint64) error
	GetByID(id uint64, currentUserID uint64) (*models.Task, error)
	ListByProjectID(projectID uint64, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error)
	ListByUserID(currentUserID uint64, offset, limit int) ([]models.Task, error)
	Update(task *models.Task, currentUserID uint64) error
	Delete(id uint64, currentUserID uint64) error
	GetByStatus(projectID uint64, status uint8, currentUserID uint64) ([]models.Task, error)
	GetByDeadlineBefore(projectID uint64, deadline time.Time, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error)
	GetByDeadlineAfter(projectID uint64, deadline time.Time, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error)
	GetByTimePeriod(projectID uint64, start, end time.Time, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error)
	AddBudget(taskID uint64, budget *models.TaskBudget, currentUserID uint64) error
	UpdateBudget(budget *models.TaskBudget, currentUserID uint64) error
	DeleteBudget(budgetID uint64, currentUserID uint64) error
	GetBudgetByTaskID(taskID uint64, currentUserID uint64) ([]models.TaskBudget, error)
}

type taskService struct {
	transactor     repository.Transactor
	taskRepo       repository.TaskRepository
	taskBudgetRepo repository.TaskBudgetRepository
	projectRepo    repository.ProjectRepository
}

func NewTaskService(
	transactor repository.Transactor,
	taskRepo repository.TaskRepository,
	taskBudgetRepo repository.TaskBudgetRepository,
	projectRepo repository.ProjectRepository,
) TaskService {
	return &taskService{
		transactor:     transactor,
		taskRepo:       taskRepo,
		taskBudgetRepo: taskBudgetRepo,
		projectRepo:    projectRepo,
	}
}

// checkTaskOwnership 检查任务是否存在且属于当前用户
func (s *taskService) checkTaskOwnership(taskID uint64, currentUserID uint64) (*models.Task, error) {
	task, err := s.taskRepo.GetByUserID(currentUserID, taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrForbidden
		}
		log.Printf("failed to get task %d: %v", taskID, err)
		return nil, ErrInternal
	}
	return task, nil
}

// checkProjectOwnership 检查项目是否存在且属于当前用户（项目ID可能为0）
func (s *taskService) checkProjectOwnership(projectID uint64, currentUserID uint64) error {
	_, err := s.projectRepo.GetByUserID(currentUserID, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrForbidden
		}
		log.Printf("failed to get project %d: %v", projectID, err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) Create(task *models.Task, currentUserID uint64) error {
	task.UserID = currentUserID
	if task.ProjectID != 0 {
		if err := s.checkProjectOwnership(task.ProjectID, currentUserID); err != nil {
			return err
		}
	}
	if err := s.taskRepo.Create(task); err != nil {
		log.Printf("failed to create task: %v", err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) GetByID(id uint64, currentUserID uint64) (*models.Task, error) {
	return s.checkTaskOwnership(id, currentUserID)
}

func (s *taskService) ListByProjectID(projectID uint64, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error) {
	if err := s.checkProjectOwnership(projectID, currentUserID); err != nil {
		return nil, 0, err
	}
	tasks, total, err := s.taskRepo.ListByProjectID(projectID, page, pageSize)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidInput) {
			return nil, 0, ErrInvalidInput
		}
		log.Printf("failed to list tasks for project %d: %v", projectID, err)
		return nil, 0, ErrInternal
	}
	return tasks, total, nil
}

func (s *taskService) ListByUserID(currentUserID uint64, offset, limit int) ([]models.Task, error) {
	tasks, err := s.taskRepo.ListByUserID(currentUserID, offset, limit)
	if err != nil {
		log.Printf("failed to list tasks for user %d: %v", currentUserID, err)
		return nil, ErrInternal
	}
	return tasks, nil
}

func (s *taskService) Update(task *models.Task, currentUserID uint64) error {
	// 检查任务所有权
	_, err := s.checkTaskOwnership(task.ID, currentUserID)
	if err != nil {
		return err
	}
	if task.ProjectID != 0 {
		if err := s.checkProjectOwnership(task.ProjectID, currentUserID); err != nil {
			return err
		}
	}
	if err := s.taskRepo.Update(task); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Printf("failed to update task %d: %v", task.ID, err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) Delete(id uint64, currentUserID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		// 检查任务所有权
		task, err := txRepo.Task.GetByID(id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrTaskNotFound
			}
			log.Printf("failed to get task %d in transaction: %v", id, err)
			return ErrInternal
		}
		if task.UserID != currentUserID {
			return ErrForbidden
		}

		// 删除关联预算
		if err := txRepo.TaskBudget.DeleteByTaskID(id); err != nil {
			log.Printf("failed to delete task budgets: %v", err)
			return ErrInternal
		}

		// 删除任务
		if err := txRepo.Task.Delete(id); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrTaskNotFound
			}
			log.Printf("failed to delete task %d: %v", id, err)
			return ErrInternal
		}
		return nil
	})
}

func (s *taskService) GetByStatus(projectID uint64, status uint8, currentUserID uint64) ([]models.Task, error) {
	if err := s.checkProjectOwnership(projectID, currentUserID); err != nil {
		return nil, err
	}
	tasks, err := s.taskRepo.GetByStatus(projectID, status)
	if err != nil {
		log.Printf("failed to get tasks by status: %v", err)
		return nil, ErrInternal
	}
	return tasks, nil
}

func (s *taskService) GetByDeadlineBefore(projectID uint64, deadline time.Time, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error) {
	if err := s.checkProjectOwnership(projectID, currentUserID); err != nil {
		return nil, 0, err
	}
	tasks, total, err := s.taskRepo.GetByDeadlineBefore(projectID, deadline, page, pageSize)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidInput) {
			return nil, 0, ErrInvalidInput
		}
		log.Printf("failed to get tasks by deadline before: %v", err)
		return nil, 0, ErrInternal
	}
	return tasks, total, nil
}

func (s *taskService) GetByDeadlineAfter(projectID uint64, deadline time.Time, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error) {
	if err := s.checkProjectOwnership(projectID, currentUserID); err != nil {
		return nil, 0, err
	}
	tasks, total, err := s.taskRepo.GetByDeadlineAfter(projectID, deadline, page, pageSize)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidInput) {
			return nil, 0, ErrInvalidInput
		}
		log.Printf("failed to get tasks by deadline after: %v", err)
		return nil, 0, ErrInternal
	}
	return tasks, total, nil
}

func (s *taskService) GetByTimePeriod(projectID uint64, start, end time.Time, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error) {
	if err := s.checkProjectOwnership(projectID, currentUserID); err != nil {
		return nil, 0, err
	}
	tasks, total, err := s.taskRepo.GetByTimePeriod(projectID, start, end, page, pageSize)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidInput) {
			return nil, 0, ErrInvalidInput
		}
		log.Printf("failed to get tasks by time period: %v", err)
		return nil, 0, ErrInternal
	}
	return tasks, total, nil
}

func (s *taskService) AddBudget(taskID uint64, budget *models.TaskBudget, currentUserID uint64) error {
	// 检查任务所有权
	_, err := s.checkTaskOwnership(taskID, currentUserID)
	if err != nil {
		return err
	}
	budget.TaskID = taskID
	if err := s.taskBudgetRepo.Create(budget); err != nil {
		log.Printf("failed to create task budget: %v", err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) UpdateBudget(budget *models.TaskBudget, currentUserID uint64) error {
	// 获取预算对应的任务ID
	existing, err := s.taskBudgetRepo.GetByID(budget.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to get task budget %d: %v", budget.ID, err)
		return ErrInternal
	}
	// 检查任务所有权
	_, err = s.checkTaskOwnership(existing.TaskID, currentUserID)
	if err != nil {
		return err
	}
	if err := s.taskBudgetRepo.Update(budget); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to update task budget %d: %v", budget.ID, err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) DeleteBudget(budgetID uint64, currentUserID uint64) error {
	// 获取预算对应的任务ID
	existing, err := s.taskBudgetRepo.GetByID(budgetID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to get task budget %d: %v", budgetID, err)
		return ErrInternal
	}
	// 检查任务所有权
	_, err = s.checkTaskOwnership(existing.TaskID, currentUserID)
	if err != nil {
		return err
	}
	if err := s.taskBudgetRepo.Delete(budgetID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to delete task budget %d: %v", budgetID, err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) GetBudgetByTaskID(taskID uint64, currentUserID uint64) ([]models.TaskBudget, error) {
	// 检查任务所有权
	_, err := s.checkTaskOwnership(taskID, currentUserID)
	if err != nil {
		return nil, err
	}
	budgets, err := s.taskBudgetRepo.GetByTaskID(taskID)
	if err != nil {
		log.Printf("failed to get budgets for task %d: %v", taskID, err)
		return nil, ErrInternal
	}
	return budgets, nil
}
