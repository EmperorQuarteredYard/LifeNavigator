package repository

import (
	"LifeNavigator/internal/interfaces/repositoryInte"
	"context"

	"gorm.io/gorm"
)

type TxRepositories struct {
	Project       repositoryInte.ProjectRepository
	ProjectBudget repositoryInte.ProjectBudgetRepository
	Task          repositoryInte.TaskRepository
	TaskPayment   repositoryInte.TaskBudgetRepository
	Account       repositoryInte.AccountRepository
	User          repositoryInte.UserRepository
}

func NewTransactor(db *gorm.DB) repositoryInte.Transactor {
	return &transactor{db: db}
}

type transactor struct {
	db *gorm.DB
}

func (t *transactor) WithinTransaction(fn func(txRepo TxRepositories) error) error {
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
