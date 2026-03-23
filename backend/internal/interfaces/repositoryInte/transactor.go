package repositoryInte

// TxRepositories 聚合了事务中需要使用的所有 Repository 接口。
// 用于在事务范围内传递一致的数据访问对象。
type TxRepositories struct {
	Project       ProjectRepository
	ProjectBudget ProjectBudgetRepository
	Task          TaskRepository
	TaskPayment   TaskBudgetRepository
	Account       AccountRepository
	User          UserRepository
}

// Transactor 定义了在事务中执行操作的能力。
type Transactor interface {
	// WithinTransaction 在事务上下文中执行给定的函数。
	// param ctx: 上下文（可用于传递请求范围的值，如用户信息）。
	// 参数 fn: 接收一个 TxRepositories 对象，包含事务专用的 Repository 实例。
	// 返回 fn 的错误，如果 fn 返回 nil 则提交事务，否则回滚。
	WithinTransaction(fn func(txRepo TxRepositories) error) error
}
