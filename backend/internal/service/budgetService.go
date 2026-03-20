package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/refresh"
	"LifeNavigator/pkg/scheduler"
	"context"
	"errors"
	"log"
	"time"
)

type BudgetService interface {
	// ProjectBudget 相关
	AddBudget(userID, projectID uint64, budget *models.ProjectBudget) error
	UpdateBudget(userID uint64, budget *models.ProjectBudget) error
	DeleteBudget(userID, budgetID uint64) error
	RefreshBudget(projectID uint64) error

	// TaskPayment 相关
	AddPayment(userID, taskID uint64, payment *models.TaskPayment) error
	UpdatePayment(userID uint64, payment *models.TaskPayment) error
	DeletePayment(userID, paymentID uint64) error
	GetPaymentByTaskID(userID, taskID uint64) ([]*dto.TaskPaymentResponse, error)

	// 调度控制（原 projectService 中与预算刷新相关的调度）
	StartAutoRefresh() error
	EndAutoRefresh()
}

func NewBudgetService(
	taskRepo repository.TaskRepository,
	taskBudgetRepo repository.TaskBudgetRepository,
	projectBudgetRepo repository.ProjectBudgetRepository,
	accountRepo repository.AccountRepository,
	projectRepo repository.ProjectRepository,
	transactor repository.Transactor,
) BudgetService {
	// 调度器暂时用空实现，StartAutoRefresh 中会重新创建
	scheduleServ := scheduler.NewScheduleService(func(uint64) error { return nil }, func(int, int) (int64, []*scheduler.Schedule) { return 0, nil })
	return &budgetService{
		transactor:        transactor,
		taskRepo:          taskRepo,
		taskBudgetRepo:    taskBudgetRepo,
		projectBudgetRepo: projectBudgetRepo,
		accountRepo:       accountRepo,
		projectRepo:       projectRepo,
		scheduleService:   scheduleServ,
	}
}

type budgetService struct {
	transactor        repository.Transactor
	taskRepo          repository.TaskRepository
	taskBudgetRepo    repository.TaskBudgetRepository
	projectBudgetRepo repository.ProjectBudgetRepository
	accountRepo       repository.AccountRepository
	projectRepo       repository.ProjectRepository
	scheduleService   *scheduler.ScheduleService
}

// ---------- 权限检查（从原 service 复制，不改变逻辑）----------
func (s *budgetService) checkTaskOwnership(userID, taskID uint64) error {
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

func (s *budgetService) checkProjectOwnership(userID, projectID uint64) error {
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

// ---------- ProjectBudget 相关实现 ----------
func (s *budgetService) AddBudget(userID, projectID uint64, budget *models.ProjectBudget) error {
	if err := s.checkProjectOwnership(userID, projectID); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		budget.ProjectID = projectID
		if err := txRepo.ProjectBudget.Create(budget); err != nil {
			log.Printf("failed to create project budget: %v", err)
			return ErrInternal
		}
		if budget.AccountID != 0 {
			if _, err := txRepo.Account.AdjustNetBalance(budget.AccountID, -1*budget.Budget); err != nil {
				log.Printf("failed to adjust net balance: %v", err)
				return ErrInternal
			}
		}
		return nil
	})
}

func (s *budgetService) UpdateBudget(userID uint64, budget *models.ProjectBudget) error {
	existing, err := s.projectBudgetRepo.GetByID(budget.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to get budget %d: %v", budget.ID, err)
		return ErrInternal
	}

	if err := s.checkProjectOwnership(userID, existing.ProjectID); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		oldBudget, err := txRepo.ProjectBudget.GetByID(budget.ID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}

		if budget.AccountID != oldBudget.AccountID {
			if oldBudget.AccountID != 0 {
				if _, err := txRepo.Account.AdjustNetBalance(oldBudget.AccountID, oldBudget.Budget); err != nil {
					return ErrInternal
				}
			}
			if budget.AccountID != 0 {
				if _, err := txRepo.Account.AdjustNetBalance(budget.AccountID, -1*budget.Budget); err != nil {
					return ErrInternal
				}
			}
		} else if budget.Budget != oldBudget.Budget && budget.AccountID != 0 {
			if _, err := txRepo.Account.AdjustNetBalance(oldBudget.AccountID, oldBudget.Budget-budget.Budget); err != nil {
				return ErrInternal
			}
		}

		if err := txRepo.ProjectBudget.Update(budget); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}
		return nil
	})
}

func (s *budgetService) DeleteBudget(userID, budgetID uint64) error {
	existing, err := s.projectBudgetRepo.GetByID(budgetID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to get budget %d: %v", budgetID, err)
		return ErrInternal
	}

	if err := s.checkProjectOwnership(userID, existing.ProjectID); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		budget, err := txRepo.ProjectBudget.GetByID(budgetID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}

		if budget.AccountID != 0 {
			if _, err := txRepo.Account.AdjustNetBalance(budget.AccountID, budget.Budget); err != nil {
				return ErrInternal
			}
		}

		if err := txRepo.ProjectBudget.Delete(budgetID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}
		return nil
	})
}

func (s *budgetService) RefreshBudget(projectID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		project, err := s.projectRepo.GetByID(projectID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				return ErrProjectNotFound
			default:
				log.Printf("fail to get project %d :%v", projectID, err)
				return ErrInternal
			}
		}
		project.LastRefresh = time.Now()
		if err = txRepo.Project.Update(project); err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				return ErrProjectNotFound
			default:
				log.Printf("fail to update project %d :%v", projectID, err)
				return ErrInternal
			}
		}
		budgets, err := txRepo.ProjectBudget.GetByProjectID(projectID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				return ErrBudgetNotFound
			}
		}
		for _, budget := range budgets {
			if _, err = txRepo.Account.AdjustNetBalance(budget.AccountID, budget.Used); err != nil {
				switch {
				case errors.Is(err, repository.ErrNotFound):
					return ErrAccountNotFound
				default:
					log.Printf("fail to add budget %d :%v", budget.ID, err)
					return ErrInternal
				}
			}
		}
		if _, strategy := refresh.ReduceRefreshInterval(project.RefreshInterval); strategy != refresh.TypeRefreshNever {
			s.scheduleService.AddRefreshSchedule(projectID, *refresh.GetNextRefreshTime(project.RefreshInterval, project.LastRefresh))
		}

		return nil
	})
}

// ---------- TaskPayment 相关实现 ----------
func (s *budgetService) AddPayment(userID, taskID uint64, payment *models.TaskPayment) error {
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

func (s *budgetService) UpdatePayment(userID uint64, payment *models.TaskPayment) error {
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

func (s *budgetService) DeletePayment(userID, paymentID uint64) error {
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

func (s *budgetService) GetPaymentByTaskID(userID, taskID uint64) ([]*dto.TaskPaymentResponse, error) {
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

// ---------- 调度控制（原 projectService 中与预算刷新相关的部分）----------
func (s *budgetService) StartAutoRefresh() error {
	var schExeFunc = func(ID uint64) error {
		return s.RefreshBudget(ID)
	}
	var schGetFunc = func(page, pageSize int) (int64, []*scheduler.Schedule) {
		sch, totalZ, errZ := s.projectRepo.GetRefreshInformation(page, pageSize)
		if errZ != nil {
			log.Printf("failed to start auto-refresh:%v", errZ)
		}
		return totalZ, sch
	}
	s.scheduleService = scheduler.NewScheduleService(schExeFunc, schGetFunc)
	s.scheduleService.Start()
	return nil
}

func (s *budgetService) EndAutoRefresh() {
	s.scheduleService.Stop()
}
