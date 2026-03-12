package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"context"
	"errors"
	"log"
	"time"
)

type AccountService interface {
	CreateAccount(*models.Account) (*models.Account, error)                                                             //创建账户
	DeleteAccount(*models.Account) error                                                                                //更新账户
	AdjustBalance(userID, accountID uint64, amount float64) (float64, error)                                            //同时调整Balance,NetBalance字段
	GetByAccountID(userID, accountID uint64) (*models.Account, error)                                                   //通过 ID获取单个账户信息
	ListByUserID(userID uint64) ([]models.Account, error)                                                               //列出一个用户拥有的所有账户
	ListLinkedTask(userID, accountID uint64, startTime, endTime time.Time) ([]models.Task, []models.TaskPayment, error) //列出关联的所有任务
	ListLinkedTaskPayment(userID, accountID uint64, startTime, endTime time.Time) ([]models.TaskPayment, error)         //列出关联的所有账单
}

func NewAccountService(
	accountRepo repository.AccountRepository,
	taskRepo repository.TaskRepository,
	taskBudgetServ repository.TaskBudgetRepository,
	transactor repository.Transactor,
) AccountService {
	return &accountService{accountRepo: accountRepo, taskRepo: taskRepo, taskBudgetServ: taskBudgetServ, transactor: transactor}
}

type accountService struct {
	accountRepo    repository.AccountRepository
	taskRepo       repository.TaskRepository
	taskBudgetServ repository.TaskBudgetRepository
	transactor     repository.Transactor
}

func (s *accountService) CreateAccount(account *models.Account) (*models.Account, error) {
	acc, err := s.accountRepo.Create(account)
	if err != nil {
		log.Printf("CreateAccount error: %v", err)
		return nil, ErrInternal
	}
	return acc, nil
}

func (s *accountService) DeleteAccount(account *models.Account) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		// 1. 查找所有引用该账户的项目预算
		budgets, err := txRepo.ProjectBudget.GetByAccountID(account.ID)
		if err != nil {
			return ErrInternal
		}

		// 2. 将这些预算的账户ID置零
		for _, b := range budgets {
			if err := txRepo.ProjectBudget.UpdateAccountID(b.ID, 0); err != nil {
				return ErrInternal
			}
		}

		// 3. 删除账户
		if err := txRepo.Account.Delete(account); err != nil {
			return ErrInternal
		}
		return nil
	})
}

func (s *accountService) AdjustBalance(userID, accountID uint64, amount float64) (float64, error) {
	balance, err := s.accountRepo.AdjustBalance(userID, accountID, amount)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return 0, ErrAccountNotFound
		default:
			log.Printf("AdjustBalance error: %v", err)
			return 0, ErrInternal
		}
	}
	return balance, nil
}

func (s *accountService) GetByAccountID(userID, accountID uint64) (*models.Account, error) {
	account, err := s.accountRepo.GetByID(userID, accountID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			return nil, ErrAccountNotFound
		default:
			log.Printf("ListByAccountID error: %v", err)
			return nil, ErrInternal
		}
	}
	return account, nil
}

func (s *accountService) ListByUserID(userID uint64) ([]models.Account, error) {
	accounts, err := s.accountRepo.ListByUserID(userID)
	if err != nil {
		return nil, ErrInternal
	}
	return accounts, nil
}

func (s *accountService) ListLinkedTask(userID, accountID uint64, startTime, endTime time.Time) ([]models.Task, []models.TaskPayment, error) {
	tasks, err := s.taskRepo.ListByAccountID(userID, accountID, startTime, endTime)
	if err != nil {
		log.Printf("ListLinkedTask error: %v", err)
		return nil, nil, ErrInternal
	}
	payments, err := s.taskBudgetServ.ListByAccountID(accountID, startTime, endTime)
	if err != nil {
		log.Printf("ListLinkedTask error: %v", err)
		return nil, nil, ErrInternal
	}
	return tasks, payments, nil
}

func (s *accountService) ListLinkedTaskPayment(userID, accountID uint64, startTime, endTime time.Time) ([]models.TaskPayment, error) {
	payments, err := s.taskBudgetServ.ListByAccountID(accountID, startTime, endTime)
	if err != nil {
		log.Printf("ListLinkedTask error: %v", err)
		return nil, ErrInternal
	}
	return payments, nil
}
