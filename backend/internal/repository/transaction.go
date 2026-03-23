package repository

import (
	"LifeNavigator/internal/interfaces/Repository"

	"gorm.io/gorm"
)

func NewTransactor(db *gorm.DB) Repository.Transactor {
	return &transactor{db: db}
}

type transactor struct {
	db *gorm.DB
}

func (t *transactor) WithinTransaction(fn func(txRepo Repository.TxRepositories) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		txRepo := Repository.TxRepositories{
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
