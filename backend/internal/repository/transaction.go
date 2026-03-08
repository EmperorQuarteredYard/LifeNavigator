package repository

import (
	"context"

	"gorm.io/gorm"
)

// TxRepositories 聚合了所有事务性的 repository 接口
type TxRepositories struct {
	Project       ProjectRepository
	ProjectBudget ProjectBudgetRepository
	Task          TaskRepository
	TaskPayment   TaskBudgetRepository
	Account       AccountRepository
	User          UserRepository
	// 可继续添加其他
}

// Transactor 定义了事务执行接口
type Transactor interface {
	// WithinTransaction 在事务中执行给定的函数，函数内使用事务性的 repository 实例
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
