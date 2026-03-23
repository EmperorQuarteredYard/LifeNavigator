package service

import (
	"LifeNavigator/internal/interfaces/Repository"
	"LifeNavigator/internal/interfaces/Service"
	"LifeNavigator/internal/models"
	"LifeNavigator/pkg/dto"
	"errors"
	"log"
	"time"
)

func NewAccountService(
	accountRepo Repository.AccountRepository,
	taskRepo Repository.TaskRepository,
	taskBudgetRepo Repository.TaskBudgetRepository,
	userRepo Repository.UserRepository,
	transactor Repository.Transactor,
) Service.AccountService {
	return &accountService{
		accountRepo:    accountRepo,
		taskRepo:       taskRepo,
		taskBudgetRepo: taskBudgetRepo,
		userRepo:       userRepo,
		transactor:     transactor,
	}
}

type accountService struct {
	accountRepo    Repository.AccountRepository
	taskRepo       Repository.TaskRepository
	taskBudgetRepo Repository.TaskBudgetRepository
	userRepo       Repository.UserRepository
	transactor     Repository.Transactor
}

func (s *accountService) GetUserName(userID uint64) (string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return "", Service.ErrUserNotFound
		}
		log.Printf("GetUserName error: %v", err)
		return "", Service.ErrInternal
	}
	return user.Username, nil
}

func (s *accountService) CreateAccount(userID uint64, account *models.Account) (*dto.Account, error) {
	created, err := s.accountRepo.Create(account, []uint64{userID})
	if err != nil {
		log.Printf("CreateAccount error: %v", err)
		return nil, Service.ErrInternal
	}
	return &dto.Account{
		ID:         created.ID,
		Name:       created.Name,
		Type:       created.Type,
		Unit:       created.Unit,
		Balance:    created.Balance,
		NetBalance: created.NetBalance,
	}, nil
}

func (s *accountService) DeleteAccount(userID, accountID uint64) error {
	owned, err := s.accountRepo.CheckOwnership(userID, accountID)
	if err != nil {
		log.Printf("CheckOwnership error: %v", err)
		return Service.ErrInternal
	}
	if !owned {
		return Service.ErrForbidden
	}

	return s.transactor.WithinTransaction(func(txRepo Repository.TxRepositories) error {
		budgets, err := txRepo.ProjectBudget.GetByAccountID(accountID)
		if err != nil {
			return Service.ErrInternal
		}

		for _, b := range budgets {
			if err := txRepo.ProjectBudget.UpdateAccountID(b.ID, 0); err != nil {
				return Service.ErrInternal
			}
		}

		account := &models.Account{ID: accountID}
		if err := txRepo.Account.Delete(account); err != nil {
			return Service.ErrInternal
		}
		return nil
	})
}

func (s *accountService) AdjustBalance(userID, accountID uint64, amount float64) (float64, error) {
	owned, err := s.accountRepo.CheckOwnership(userID, accountID)
	if err != nil {
		log.Printf("CheckOwnership error: %v", err)
		return 0, Service.ErrInternal
	}
	if !owned {
		return 0, Service.ErrForbidden
	}

	balance, err := s.accountRepo.AdjustBalance(accountID, amount)
	if err != nil {
		switch {
		case errors.Is(err, Repository.ErrNotFound):
			return 0, Service.ErrAccountNotFound
		default:
			log.Printf("AdjustBalance error: %v", err)
			return 0, Service.ErrInternal
		}
	}
	return balance, nil
}

func (s *accountService) GetByAccountID(userID, accountID uint64) (*dto.Account, error) {
	owned, err := s.accountRepo.CheckOwnership(userID, accountID)
	if err != nil {
		log.Printf("CheckOwnership error: %v", err)
		return nil, Service.ErrInternal
	}
	if !owned {
		return nil, Service.ErrForbidden
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		switch {
		case errors.Is(err, Repository.ErrNotFound):
			return nil, Service.ErrAccountNotFound
		default:
			log.Printf("GetByAccountID error: %v", err)
			return nil, Service.ErrInternal
		}
	}
	return &dto.Account{
		ID:         account.ID,
		Name:       account.Name,
		Type:       account.Type,
		Unit:       account.Unit,
		Balance:    account.Balance,
		NetBalance: account.NetBalance,
	}, nil
}

func (s *accountService) ListByUserID(userID uint64) (*dto.AccountList, error) {
	accounts, err := s.accountRepo.ListByUserID(userID)
	if err != nil {
		log.Printf("ListByUserID error: %v", err)
		return nil, Service.ErrInternal
	}

	items := make([]*dto.Account, len(accounts))
	for i, acc := range accounts {
		items[i] = &dto.Account{
			ID:         acc.ID,
			Name:       acc.Name,
			Type:       acc.Type,
			Unit:       acc.Unit,
			Balance:    acc.Balance,
			NetBalance: acc.NetBalance,
		}
	}
	return &dto.AccountList{Items: items}, nil
}

func (s *accountService) ListLinkedTask(userID, accountID uint64, startTime, endTime time.Time) (*dto.TaskList, error) {
	owned, err := s.accountRepo.CheckOwnership(userID, accountID)
	if err != nil {
		log.Printf("CheckOwnership error: %v", err)
		return nil, Service.ErrInternal
	}
	if !owned {
		return nil, Service.ErrForbidden
	}

	tasks, err := s.taskRepo.ListByAccountID(accountID, startTime, endTime)
	if err != nil {
		log.Printf("ListLinkedTask error: %v", err)
		return nil, Service.ErrInternal
	}

	list := make([]*dto.Task, len(tasks))
	for i, task := range tasks {
		payments, _ := s.taskBudgetRepo.GetByTaskID(task.ID)
		paymentDtos := make([]*dto.TaskPayment, len(payments))
		for j, p := range payments {
			paymentDtos[j] = &dto.TaskPayment{
				ID:       p.ID,
				TaskID:   p.TaskID,
				BudgetID: p.BudgetID,
				Amount:   p.Amount,
			}
		}
		list[i] = &dto.Task{
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

	return &dto.TaskList{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
