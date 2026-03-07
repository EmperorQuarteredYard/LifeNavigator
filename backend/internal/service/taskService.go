package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"context"
	"errors"
	"log"
	"time"
)

type TaskService interface {
	Create(task *models.Task, currentUserID uint64) error                                                     //创建任务
	GetByID(id uint64, currentUserID uint64) (*models.Task, error)                                            //通过任务 ID获取任务详情
	ListByProjectID(projectID uint64, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error) //列举一个项目下的所有任务
	ListByUserID(currentUserID uint64, offset, limit int) ([]models.Task, int64, error)                       //列举一个用户的所有任务
	Update(task *models.Task, currentUserID uint64) error                                                     //更新任务
	Delete(id uint64, currentUserID uint64) error                                                             //删除任务，并删除关联关系
	GetByStatus(projectID uint64, status uint8, currentUserID uint64) ([]models.Task, error)
	GetByDeadlineBefore(projectID uint64, deadline time.Time, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error)
	GetByDeadlineAfter(projectID uint64, deadline time.Time, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error)
	GetByTimePeriod(projectID uint64, start, end time.Time, page, pageSize int, currentUserID uint64) ([]models.Task, int64, error)
	AddPayment(taskID uint64, budget *models.TaskPayment, currentUserID uint64) error
	UpdatePayment(budget *models.TaskPayment, currentUserID uint64) error
	DeletePayment(budgetID uint64, currentUserID uint64) error
	GetPaymentByTaskID(taskID uint64, currentUserID uint64) ([]models.TaskPayment, error)
	SetPrerequisiteTask(prerequisiteID, taskID uint64) (dependency *models.TaskDependency, err error)
	UnsetPrerequisiteTask(prerequisiteID, taskID, userID uint64) (err error) //不会检查ID是否有效、用户是否正确
	GetPrerequisites(taskID, userID uint64) (prerequisites []models.TaskDependency, err error)
	GetPostrequisite(taskID, userID uint64) (prerequisites []models.TaskDependency, err error)
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

type taskService struct {
	transactor     repository.Transactor
	taskRepo       repository.TaskRepository
	taskBudgetRepo repository.TaskBudgetRepository
	projectRepo    repository.ProjectRepository
}

func (s *taskService) SetPrerequisiteTask(prerequisiteID, taskID uint64) (dependency *models.TaskDependency, err error) {
	dependency, err = s.taskRepo.SetPrerequisiteTask(prerequisiteID, taskID)
	if err != nil {
		log.Println("Fail to set prerequisite task", err)
		return nil, ErrInternal
	}
	return dependency, nil
}

func (s *taskService) UnsetPrerequisiteTask(prerequisiteID, taskID, userID uint64) (err error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrTaskNotFound
		}
		log.Println("Fail to get task:", err)
		return ErrInternal
	}
	if task.UserID != userID {
		return ErrForbidden
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

func (s *taskService) GetPrerequisites(taskID, userID uint64) (prerequisites []models.TaskDependency, err error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		log.Println("Fail to get task:", err)
		return nil, ErrInternal
	}
	if task.UserID != userID {
		return nil, ErrForbidden
	}
	prerequisites, err = s.taskRepo.GetPrerequisites(taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskDependencyNotFound
		}
		log.Println("Fail to get prerequisites by taskID:", err)
		return nil, ErrInternal
	}
	return prerequisites, nil
}

func (s *taskService) GetPostrequisite(taskID, userID uint64) (prerequisites []models.TaskDependency, err error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskNotFound
		}
		log.Println("Fail to get task:", err)
		return nil, ErrInternal
	}
	if task.UserID != userID {
		return nil, ErrForbidden
	}
	prerequisites, err = s.taskRepo.GetPostrequisites(taskID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrTaskDependencyNotFound
		}
		log.Println("Fail to get prerequisites by TaskID:", err)
		return nil, ErrInternal
	}
	return prerequisites, nil
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

func (s *taskService) ListByUserID(currentUserID uint64, offset, limit int) ([]models.Task, int64, error) {
	tasks, total, err := s.taskRepo.ListByUserID(currentUserID, offset, limit)
	if err != nil {
		log.Printf("failed to list tasks for user %d: %v", currentUserID, err)
		return nil, 0, ErrInternal
	}
	return tasks, total, nil
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
		if err := txRepo.TaskPayment.DeleteByTaskID(id); err != nil {
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

func (s *taskService) AddPayment(taskID uint64, payment *models.TaskPayment, currentUserID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		task, err := txRepo.Task.GetByUserID(currentUserID, taskID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrForbidden
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
			return ErrInvalidInput // 预算不属于该任务所在的项目
		}
		if err := txRepo.ProjectBudget.AddUsed(payment.BudgetID, payment.Amount); err != nil {
			log.Printf("failed to add used to project budgets: %v", err)
			return ErrInternal
		}
		if projBudget.AccountID != 0 {
			_, err := txRepo.Account.AdjustBalance(currentUserID, projBudget.AccountID, -payment.Amount)
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

func (s *taskService) UpdatePayment(payment *models.TaskPayment, currentUserID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		oldPayment, err := txRepo.TaskPayment.GetByID(payment.ID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}

		//  验证任务所有权
		_, err = txRepo.Task.GetByUserID(currentUserID, oldPayment.TaskID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrForbidden
			}
			return ErrInternal
		}

		//不允许修改 BudgetID
		if payment.BudgetID != oldPayment.BudgetID {
			return ErrInvalidInput
		}

		projBudget, err := txRepo.ProjectBudget.GetByID(payment.BudgetID)
		if err != nil {
			return ErrInternal
		}

		//计算金额变化，并同步预算与账户
		delta := payment.Amount - oldPayment.Amount
		if delta != 0 {
			if delta > 0 {
				// 增加支出
				if err := txRepo.ProjectBudget.AddUsed(payment.BudgetID, delta); err != nil {
					return ErrInternal
				}
				if projBudget.AccountID != 0 {
					_, err := txRepo.Account.AdjustBalance(currentUserID, projBudget.AccountID, -delta)
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
				// 减少支出
				dec := -delta
				if err := txRepo.ProjectBudget.SubtractUsed(payment.BudgetID, dec); err != nil {
					return ErrInternal
				}
				if projBudget.AccountID != 0 {
					_, err := txRepo.Account.AdjustBalance(currentUserID, projBudget.AccountID, dec)
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

func (s *taskService) DeletePayment(paymentID uint64, currentUserID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		payment, err := txRepo.TaskPayment.GetByID(paymentID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}

		_, err = txRepo.Task.GetByUserID(currentUserID, payment.TaskID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrForbidden
			}
			return ErrInternal
		}

		projBudget, err := txRepo.ProjectBudget.GetByID(payment.BudgetID)
		if err != nil {
			return ErrInternal
		}

		if err := txRepo.ProjectBudget.SubtractUsed(payment.BudgetID, payment.Amount); err != nil {
			return ErrInternal
		}

		if projBudget.AccountID != 0 {
			_, err := txRepo.Account.AdjustBalance(currentUserID, projBudget.AccountID, payment.Amount)
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

func (s *taskService) GetPaymentByTaskID(taskID uint64, currentUserID uint64) ([]models.TaskPayment, error) {
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
