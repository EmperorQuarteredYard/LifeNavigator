package repository

import (
	"context"

	"gorm.io/gorm"
)

type TxRepositories struct {
	Project       ProjectRepository
	ProjectBudget ProjectBudgetRepository
	Task          TaskRepository
	TaskPayment   TaskBudgetRepository
	Account       AccountRepository
	User          UserRepository
}

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(txRepo TxRepositories) error) error
}

type transactor struct {
	db *gorm.DB
}

func NewTransactor(db *gorm.DB) Transactor {
	return &transactor{db: db}
}

func (t *transactor) WithinTransaction(ctx context.Context, fn func(txRepo TxRepositories) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		txRepo := TxRepositories{
			Project:       NewProjectRepository(tx),
			ProjectBudget: NewProjectBudgetRepository(tx),
			Task:          NewTaskRepository(tx),
			TaskPayment:   NewTaskPaymentRepository(tx),
			Account:       NewAccountRepository(tx),
			User:          NewUserRepository(tx),
		}
		return fn(txRepo)
	})
}
