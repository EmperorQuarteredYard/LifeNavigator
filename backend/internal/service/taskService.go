package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/permission"
	"context"
	"errors"
	"log"
	"time"
)

type TaskService interface {
	Create(userID uint64, task *models.Task) (*dto.TaskResponse, error)
	GetByID(userID, id uint64) (*dto.TaskResponse, error)
	ListByProjectID(userID, projectID uint64, page, pageSize int) (*dto.TaskListResponse, error)
	Update(userID uint64, task *models.Task) error
	UpdateStatus(userID, id uint64, status uint8) error
	Delete(userID, id uint64) error
	GetByStatus(userID, projectID uint64, status uint8) ([]*dto.TaskResponse, error)
	GetByDeadlineBefore(userID, projectID uint64, deadline time.Time, page, pageSize int) (*dto.TaskListResponse, error)
	GetByDeadlineAfter(userID, projectID uint64, deadline time.Time, page, pageSize int) (*dto.TaskListResponse, error)
	GetByTimePeriod(userID, projectID uint64, start, end time.Time, page, pageSize int) (*dto.TaskListResponse, error)
	SetPrerequisiteTask(userID, prerequisiteID, taskID uint64) (*dto.DependencyResponse, error)
	UnsetPrerequisiteTask(userID, prerequisiteID, taskID uint64) error
	GetPrerequisites(userID, taskID uint64) ([]*dto.DependencyResponse, error)
	GetPostrequisite(userID, taskID uint64) ([]*dto.DependencyResponse, error)
}

func NewTaskService(
	transactor repository.Transactor,
	taskRepo repository.TaskRepository,
	taskBudgetRepo repository.TaskBudgetRepository,
	projectRepo repository.ProjectRepository,
	accountRepo repository.AccountRepository,
) TaskService {
	return &taskService{
		transactor:     transactor,
		taskRepo:       taskRepo,
		taskBudgetRepo: taskBudgetRepo,
		projectBase:    &projectBase{projectRepo: projectRepo},
		accountRepo:    accountRepo,
	}
}

type taskService struct {
	transactor     repository.Transactor
	taskRepo       repository.TaskRepository
	taskBudgetRepo repository.TaskBudgetRepository
	accountRepo    repository.AccountRepository
	*projectBase
}

func (s *taskService) toTaskResponse(task *models.Task) *dto.TaskResponse {
	payments, _ := s.taskBudgetRepo.GetByTaskID(task.ID) //TODO 这里应当迁移到budgetService
	paymentDtos := make([]*dto.TaskPaymentResponse, len(payments))
	for i, p := range payments {
		paymentDtos[i] = &dto.TaskPaymentResponse{
			ID:       p.ID,
			TaskID:   p.TaskID,
			BudgetID: p.BudgetID,
			Amount:   p.Amount,
		}
	}
	return &dto.TaskResponse{
		ID:          task.ID,
		ProjectID:   task.ProjectID,
		Name:        task.Name,
		Description: task.Description,
		Type:        task.Type,
		Status:      task.Status,
		Category:    task.Category,
		Deadline:    task.Deadline,
		CompletedAt: task.CompletedAt,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Payments:    paymentDtos,
	}
}

func (s *taskService) Create(userID uint64, task *models.Task) (*dto.TaskResponse, error) {
	if task.ProjectID != 0 {
		if err := s.checkProjectAccessibility(userID, task.ProjectID, permission.OpCreate); err != nil {
			return nil, err
		}
	}
	if err := s.taskRepo.Create(task); err != nil {
		log.Printf("failed to create task: %v", err)
		return nil, ErrInternal
	}
	return s.toTaskResponse(task), nil
}

func (s *taskService) GetByID(userID, id uint64) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		log.Printf("failed to get task %d: %v", id, err)
		return nil, ErrInternal
	}
	if err := s.checkProjectAccessibility(userID, task.ProjectID, permission.OpRead); err != nil {
		return nil, err
	}
	return s.toTaskResponse(task), nil
}

func (s *taskService) ListByProjectID(userID, projectID uint64, page, pageSize int) (*dto.TaskListResponse, error) {
	if err := s.checkProjectAccessibility(userID, projectID, permission.OpRead); err != nil {
		return nil, err
	}
	tasks, total, err := s.taskRepo.ListByProjectID(projectID, page, pageSize)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidInput) {
			return nil, ErrInvalidInput
		}
		log.Printf("failed to list tasks for project %d: %v", projectID, err)
		return nil, ErrInternal
	}

	list := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		list[i] = *(s.toTaskResponse(&task))
	}
	return &dto.TaskListResponse{
		Page:  int64(page),
		Total: total,
		Size:  int64(pageSize),
		List:  list,
	}, nil
}

func (s *taskService) Update(userID uint64, task *models.Task) (err error) {
	var oldTask *models.Task
	oldTask, err = s.taskRepo.GetByID(task.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Printf("failed to get task %d: %v", task.ID, err)
		return ErrInternal
	}
	if oldTask.ProjectID != task.ProjectID && task.ProjectID != 0 { //行为：迁移
		err = s.checkProjectAccessibility(userID, oldTask.ProjectID, permission.OpDelete)
		if err != nil {
			return err
		}
	} else { //行为：更新
		err = s.checkProjectAccessibility(userID, oldTask.ProjectID, permission.OpUpdate)
		if err != nil {
			return err
		}
	}
	if err = s.taskRepo.Update(task); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Printf("failed to update task %d: %v", task.ID, err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) Delete(userID, id uint64) error {

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		var task *models.Task
		task, err := txRepo.Task.GetByID(id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrTaskNotFound
			}
			log.Printf("failed to get task %d: %v", id, err)
			return ErrInternal
		}
		if err = s.checkProjectAccessibility(userID, task.ProjectID, permission.OpDelete); err != nil {
			return err
		}
		if err = txRepo.Task.Delete(id); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrTaskNotFound
			}
			log.Printf("failed to delete task %d: %v", id, err)
			return ErrInternal
		}
		if err = txRepo.TaskPayment.DeleteByTaskID(id); err != nil {
			log.Printf("failed to delete task budgets: %v", err)
			return ErrInternal
		}
		return nil
	})
}

func (s *taskService) UpdateStatus(userID, id uint64, status uint8) (err error) {
	var oldTask *models.Task
	oldTask, err = s.taskRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Printf("failed to get task %d: %v", id, err)
		return ErrInternal
	}
	if err = s.checkProjectAccessibility(userID, oldTask.ProjectID, permission.OpUpdate); err != nil {
		return err
	}

	if err = s.taskRepo.UpdateStatus(id, status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Printf("failed to update task status %d: %v", id, err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) GetByStatus(userID, projectID uint64, status uint8) ([]*dto.TaskResponse, error) {
	if err := s.checkProjectAccessibility(userID, projectID, permission.OpRead); err != nil {
		return nil, err
	}
	tasks, err := s.taskRepo.GetByStatus(projectID, status)
	if err != nil {
		log.Printf("failed to get tasks by status: %v", err)
		return nil, ErrInternal
	}

	list := make([]*dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		list[i] = s.toTaskResponse(&task)
	}
	return list, nil
}

func (s *taskService) GetByDeadlineBefore(userID, projectID uint64, deadline time.Time, page, pageSize int) (*dto.TaskListResponse, error) {
	if err := s.checkProjectAccessibility(userID, projectID, permission.OpRead); err != nil {
		return nil, err
	}
	tasks, total, err := s.taskRepo.GetByDeadlineBefore(projectID, deadline, page, pageSize)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidInput) {
			return nil, ErrInvalidInput
		}
		log.Printf("failed to get tasks by deadline before: %v", err)
		return nil, ErrInternal
	}

	list := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		list[i] = *(s.toTaskResponse(&task))
	}
	return &dto.TaskListResponse{
		Page:  int64(page),
		Total: total,
		Size:  int64(pageSize),
		List:  list,
	}, nil
}

func (s *taskService) GetByDeadlineAfter(userID, projectID uint64, deadline time.Time, page, pageSize int) (*dto.TaskListResponse, error) {
	if err := s.checkProjectAccessibility(userID, projectID, permission.OpRead); err != nil {
		return nil, err
	}
	tasks, total, err := s.taskRepo.GetByDeadlineAfter(projectID, deadline, page, pageSize)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidInput) {
			return nil, ErrInvalidInput
		}
		log.Printf("failed to get tasks by deadline after: %v", err)
		return nil, ErrInternal
	}

	list := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		list[i] = *(s.toTaskResponse(&task))
	}
	return &dto.TaskListResponse{
		Page:  int64(page),
		Total: total,
		Size:  int64(pageSize),
		List:  list,
	}, nil
}

func (s *taskService) GetByTimePeriod(userID, projectID uint64, start, end time.Time, page, pageSize int) (*dto.TaskListResponse, error) {
	if err := s.checkProjectAccessibility(userID, projectID, permission.OpRead); err != nil {
		return nil, err
	}
	tasks, total, err := s.taskRepo.GetByTimePeriod(projectID, start, end, page, pageSize)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidInput) {
			return nil, ErrInvalidInput
		}
		log.Printf("failed to get tasks by time period: %v", err)
		return nil, ErrInternal
	}

	list := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		list[i] = *(s.toTaskResponse(&task))
	}
	return &dto.TaskListResponse{
		Page:  int64(page),
		Total: total,
		Size:  int64(pageSize),
		List:  list,
	}, nil
}

func (s *taskService) SetPrerequisiteTask(userID, prerequisiteID, taskID uint64) (*dto.DependencyResponse, error) {
	var task, prerequisite *models.Task
	var err error
	if task, err = s.taskRepo.GetByID(taskID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		log.Printf("failed to get task %d: %v", taskID, err)
		return nil, ErrInternal
	}
	if err = s.checkProjectAccessibility(userID, task.ProjectID, permission.OpUpdate); err != nil {
		return nil, err
	}
	if prerequisite, err = s.taskRepo.GetByID(taskID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		log.Printf("failed to get task %d: %v", taskID, err)
		return nil, ErrInternal
	}
	if err = s.checkProjectAccessibility(userID, prerequisite.ProjectID, permission.OpBranch); err != nil {
		return nil, err
	}

	dependency, err := s.taskRepo.SetPrerequisiteTask(prerequisiteID, taskID)
	if err != nil {
		log.Println("Fail to set prerequisite task", err)
		return nil, ErrInternal
	}
	return &dto.DependencyResponse{
		PrerequisiteID: dependency.PrerequisiteID,
		TaskID:         dependency.TaskID,
	}, nil
}

func (s *taskService) UnsetPrerequisiteTask(userID, prerequisiteID, taskID uint64) error {
	var task *models.Task
	var err error
	if task, err = s.taskRepo.GetByID(taskID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Printf("failed to get task %d: %v", taskID, err)
		return ErrInternal
	}
	if err = s.checkProjectAccessibility(userID, task.ProjectID, permission.OpUpdate); err != nil {
		return err
	}

	err = s.taskRepo.UnsetPrerequisiteTask(prerequisiteID, taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskDependencyNotFound
		}
		log.Println("Fail to unset prerequisite task:", err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) GetPrerequisites(userID, taskID uint64) ([]*dto.DependencyResponse, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskDependencyNotFound
		}
		log.Println("Fail to get prerequisites by taskID:", err)
		return nil, ErrInternal
	}

	if err = s.checkProjectAccessibility(userID, task.ProjectID, permission.OpRead); err != nil {
		return nil, err
	}
	prerequisites, err := s.taskRepo.GetPrerequisites(taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskDependencyNotFound
		}
		log.Println("Fail to get prerequisites by taskID:", err)
		return nil, ErrInternal
	}

	result := make([]*dto.DependencyResponse, len(prerequisites))
	for i, p := range prerequisites {
		result[i] = &dto.DependencyResponse{
			PrerequisiteID: p.PrerequisiteID,
			TaskID:         p.TaskID,
		}
	}
	return result, nil
}

func (s *taskService) GetPostrequisite(userID, taskID uint64) ([]*dto.DependencyResponse, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskDependencyNotFound
		}
		log.Println("Fail to get prerequisites by taskID:", err)
		return nil, ErrInternal
	}

	if err = s.checkProjectAccessibility(userID, task.ProjectID, permission.OpRead); err != nil {
		return nil, err
	}
	prerequisites, err := s.taskRepo.GetPostrequisites(taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskDependencyNotFound
		}
		log.Println("Fail to get prerequisites by TaskID:", err)
		return nil, ErrInternal
	}

	result := make([]*dto.DependencyResponse, len(prerequisites))
	for i, p := range prerequisites {
		result[i] = &dto.DependencyResponse{
			PrerequisiteID: p.PrerequisiteID,
			TaskID:         p.TaskID,
		}
	}
	return result, nil
}
