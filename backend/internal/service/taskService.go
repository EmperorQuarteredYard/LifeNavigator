package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/pkg/dto"
	"context"
	"errors"
	"log"
	"time"
)

type TaskService interface {
	Create(userID uint64, task *models.Task) (*dto.TaskResponse, error)
	GetByID(userID, id uint64) (*dto.TaskResponse, error)
	ListByProjectID(userID, projectID uint64, page, pageSize int) (*dto.TaskListResponse, error)
	ListByUserID(userID uint64, offset, limit int) (*dto.TaskListResponse, error)
	Update(userID uint64, task *models.Task) error
	UpdateStatus(userID, id uint64, status uint8) error
	Delete(userID, id uint64) error
	GetByStatus(userID, projectID uint64, status uint8) ([]*dto.TaskResponse, error)
	GetByDeadlineBefore(userID, projectID uint64, deadline time.Time, page, pageSize int) (*dto.TaskListResponse, error)
	GetByDeadlineAfter(userID, projectID uint64, deadline time.Time, page, pageSize int) (*dto.TaskListResponse, error)
	GetByTimePeriod(userID, projectID uint64, start, end time.Time, page, pageSize int) (*dto.TaskListResponse, error)
	AddPayment(userID, taskID uint64, payment *models.TaskPayment) error
	UpdatePayment(userID uint64, payment *models.TaskPayment) error
	DeletePayment(userID, paymentID uint64) error
	GetPaymentByTaskID(userID, taskID uint64) ([]*dto.TaskPaymentResponse, error)
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
		projectRepo:    projectRepo,
		accountRepo:    accountRepo,
	}
}

type taskService struct {
	transactor     repository.Transactor
	taskRepo       repository.TaskRepository
	taskBudgetRepo repository.TaskBudgetRepository
	projectRepo    repository.ProjectRepository
	accountRepo    repository.AccountRepository
}

func (s *taskService) checkTaskOwnership(userID, taskID uint64) error {
	owned, err := s.taskRepo.CheckOwnership(userID, taskID)
	if err != nil {
		log.Printf("failed to check task ownership %d: %v", taskID, err)
		return ErrInternal
	}
	if !owned {
		return ErrForbidden
	}
	return nil
}

func (s *taskService) checkProjectOwnership(userID, projectID uint64) error {
	owned, err := s.projectRepo.CheckOwnership(userID, projectID)
	if err != nil {
		log.Printf("failed to check project ownership %d: %v", projectID, err)
		return ErrInternal
	}
	if !owned {
		return ErrForbidden
	}
	return nil
}

func (s *taskService) toTaskResponse(task *models.Task) *dto.TaskResponse {
	payments, _ := s.taskBudgetRepo.GetByTaskID(task.ID)
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
		if err := s.checkProjectOwnership(userID, task.ProjectID); err != nil {
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
	if err := s.checkTaskOwnership(userID, id); err != nil {
		return nil, err
	}
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		log.Printf("failed to get task %d: %v", id, err)
		return nil, ErrInternal
	}
	return s.toTaskResponse(task), nil
}

func (s *taskService) ListByProjectID(userID, projectID uint64, page, pageSize int) (*dto.TaskListResponse, error) {
	if err := s.checkProjectOwnership(userID, projectID); err != nil {
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

func (s *taskService) ListByUserID(userID uint64, offset, limit int) (*dto.TaskListResponse, error) {
	tasks, total, err := s.taskRepo.ListByUserID(userID, offset, limit)
	if err != nil {
		log.Printf("failed to list tasks for user %d: %v", userID, err)
		return nil, ErrInternal
	}

	list := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		list[i] = *(s.toTaskResponse(&task))
	}
	return &dto.TaskListResponse{
		Total: total,
		List:  list,
	}, nil
}

func (s *taskService) Update(userID uint64, task *models.Task) error {
	if err := s.checkTaskOwnership(userID, task.ID); err != nil {
		return err
	}
	if task.ProjectID != 0 {
		if err := s.checkProjectOwnership(userID, task.ProjectID); err != nil {
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

func (s *taskService) Delete(userID, id uint64) error {
	if err := s.checkTaskOwnership(userID, id); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		if err := txRepo.TaskPayment.DeleteByTaskID(id); err != nil {
			log.Printf("failed to delete task budgets: %v", err)
			return ErrInternal
		}
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

func (s *taskService) UpdateStatus(userID, id uint64, status uint8) error {
	if err := s.checkTaskOwnership(userID, id); err != nil {
		return err
	}

	if err := s.taskRepo.UpdateStatus(id, status); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Printf("failed to update task status %d: %v", id, err)
		return ErrInternal
	}
	return nil
}

func (s *taskService) GetByStatus(userID, projectID uint64, status uint8) ([]*dto.TaskResponse, error) {
	if err := s.checkProjectOwnership(userID, projectID); err != nil {
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
	if err := s.checkProjectOwnership(userID, projectID); err != nil {
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
	if err := s.checkProjectOwnership(userID, projectID); err != nil {
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
	if err := s.checkProjectOwnership(userID, projectID); err != nil {
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

func (s *taskService) AddPayment(userID, taskID uint64, payment *models.TaskPayment) error {
	if err := s.checkTaskOwnership(userID, taskID); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		task, err := txRepo.Task.GetByID(taskID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrTaskNotFound
			}
			return ErrInternal
		}

		projBudget, err := txRepo.ProjectBudget.GetByID(payment.BudgetID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}
		if projBudget.ProjectID != task.ProjectID {
			return ErrInvalidInput
		}
		if err := txRepo.ProjectBudget.AddUsed(payment.BudgetID, payment.Amount); err != nil {
			log.Printf("failed to add used to project budgets: %v", err)
			return ErrInternal
		}
		if projBudget.AccountID != 0 {
			_, err := txRepo.Account.AdjustBalance(projBudget.AccountID, -payment.Amount)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return ErrAccountNotFound
				}
				if errors.Is(err, repository.ErrConcurrentUpdate) {
					return ErrConcurrentUpdate
				}
				log.Printf("failed to update account balance: %v", err)
				return ErrInternal
			}
		}

		payment.TaskID = taskID
		if err := txRepo.TaskPayment.Create(payment); err != nil {
			log.Printf("failed to create task payment: %v", err)
			return ErrInternal
		}
		return nil
	})
}

func (s *taskService) UpdatePayment(userID uint64, payment *models.TaskPayment) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		oldPayment, err := txRepo.TaskPayment.GetByID(payment.ID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}

		if err := s.checkTaskOwnership(userID, oldPayment.TaskID); err != nil {
			return err
		}

		if payment.BudgetID != oldPayment.BudgetID {
			return ErrInvalidInput
		}

		projBudget, err := txRepo.ProjectBudget.GetByID(payment.BudgetID)
		if err != nil {
			return ErrInternal
		}

		delta := payment.Amount - oldPayment.Amount
		if delta != 0 {
			if delta > 0 {
				if err := txRepo.ProjectBudget.AddUsed(payment.BudgetID, delta); err != nil {
					return ErrInternal
				}
				if projBudget.AccountID != 0 {
					_, err := txRepo.Account.AdjustBalance(projBudget.AccountID, -delta)
					if err != nil {
						if errors.Is(err, repository.ErrNotFound) {
							return ErrAccountNotFound
						}
						if errors.Is(err, repository.ErrConcurrentUpdate) {
							return ErrConcurrentUpdate
						}
						log.Printf("failed to update account balance: %v", err)
						return ErrInternal
					}
				}
			} else {
				dec := -delta
				if err := txRepo.ProjectBudget.SubtractUsed(payment.BudgetID, dec); err != nil {
					return ErrInternal
				}
				if projBudget.AccountID != 0 {
					_, err := txRepo.Account.AdjustBalance(projBudget.AccountID, dec)
					if err != nil {
						if errors.Is(err, repository.ErrNotFound) {
							return ErrAccountNotFound
						}
						if errors.Is(err, repository.ErrConcurrentUpdate) {
							return ErrConcurrentUpdate
						}
						log.Printf("failed to update account balance: %v", err)
						return ErrInternal
					}
				}
			}
		}

		if err := txRepo.TaskPayment.Update(payment); err != nil {
			return ErrInternal
		}
		return nil
	})
}

func (s *taskService) DeletePayment(userID, paymentID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		payment, err := txRepo.TaskPayment.GetByID(paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}

		if err := s.checkTaskOwnership(userID, payment.TaskID); err != nil {
			return err
		}

		projBudget, err := txRepo.ProjectBudget.GetByID(payment.BudgetID)
		if err != nil {
			return ErrInternal
		}

		if err := txRepo.ProjectBudget.SubtractUsed(payment.BudgetID, payment.Amount); err != nil {
			return ErrInternal
		}

		if projBudget.AccountID != 0 {
			_, err := txRepo.Account.AdjustBalance(projBudget.AccountID, payment.Amount)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return ErrAccountNotFound
				}
				if errors.Is(err, repository.ErrConcurrentUpdate) {
					return ErrConcurrentUpdate
				}
				log.Printf("failed to update account balance: %v", err)
				return ErrInternal
			}
		}

		if err := txRepo.TaskPayment.Delete(paymentID); err != nil {
			return ErrInternal
		}
		return nil
	})
}

func (s *taskService) GetPaymentByTaskID(userID, taskID uint64) ([]*dto.TaskPaymentResponse, error) {
	if err := s.checkTaskOwnership(userID, taskID); err != nil {
		return nil, err
	}
	budgets, err := s.taskBudgetRepo.GetByTaskID(taskID)
	if err != nil {
		log.Printf("failed to get budgets for task %d: %v", taskID, err)
		return nil, ErrInternal
	}

	result := make([]*dto.TaskPaymentResponse, len(budgets))
	for i, b := range budgets {
		result[i] = &dto.TaskPaymentResponse{
			ID:       b.ID,
			TaskID:   b.TaskID,
			BudgetID: b.BudgetID,
			Amount:   b.Amount,
		}
	}
	return result, nil
}

func (s *taskService) SetPrerequisiteTask(userID, prerequisiteID, taskID uint64) (*dto.DependencyResponse, error) {
	if err := s.checkTaskOwnership(userID, taskID); err != nil {
		return nil, err
	}
	if err := s.checkTaskOwnership(userID, prerequisiteID); err != nil {
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
	if err := s.checkTaskOwnership(userID, taskID); err != nil {
		return err
	}

	err := s.taskRepo.UnsetPrerequisiteTask(prerequisiteID, taskID)
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
	if err := s.checkTaskOwnership(userID, taskID); err != nil {
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
	if err := s.checkTaskOwnership(userID, taskID); err != nil {
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
