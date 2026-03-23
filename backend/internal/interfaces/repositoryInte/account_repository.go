package repositoryInte

import (
	"LifeNavigator/internal/models"
)

// AccountRepository 定义了账户相关的数据访问接口。
// 账户支持乐观锁并发更新，部分更新方法会进行重试。
type AccountRepository interface {
	// Create 创建一个新账户并关联指定的用户。
	// 参数 account: 待创建的账户对象（ID 将由实现自动生成）。
	// 参数 userIDs: 至少一个用户 ID，用于将账户与用户关联（多对多关系）。
	// 返回值:
	//   - 成功时返回创建的账户对象（包含生成的 ID）
	//   - 错误：如果 userIDs 为空返回 ErrInvalidInput；如果 type 为空返回 ErrInvalidInput；
	//     数据库错误返回相应错误（如唯一约束冲突、连接错误等）。
	Create(account *models.Account, userIDs []uint64) (*models.Account, error)

	// Delete 删除指定的账户，同时会自动删除关联的 account_users 记录。
	// 参数 account: 需要删除的账户对象，其 ID 字段必须非零。
	// 错误：
	//   - 如果 account.ID == 0 返回 ErrInvalidInput。
	//   - 数据库错误返回相应错误。
	Delete(account *models.Account) error

	// AdjustBalance 原子性地调整账户的可用余额（balance），使用乐观锁。
	// 参数 amount: 调整金额（正数为增加，负数为减少）。
	// 返回值：
	//   - 成功时返回调整后的最新余额。
	//   - 错误：如果账户不存在返回 ErrNotFound；若乐观锁冲突重试耗尽返回 ErrConcurrentUpdate。
	// 实现行为：
	//   - 最多重试 updateMaxTry 次（默认3次），每次重试间隔 10ms。
	//   - 若 amount == 0，则不更新直接返回当前余额。
	AdjustBalance(accountID uint64, amount float64) (float64, error)

	// AdjustNetBalance 原子性地调整账户的净余额（net_balance），使用乐观锁。
	// 行为与 AdjustBalance 相同，仅操作字段不同。
	AdjustNetBalance(accountID uint64, amount float64) (float64, error)

	// GetByID 根据账户 ID 查询账户信息。
	// 错误：如果账户不存在返回 ErrNotFound。
	GetByID(accountID uint64) (*models.Account, error)

	// ListByUserID 查询指定用户所关联的所有账户。
	// 参数 userID: 用户 ID。
	// 返回该用户有权访问的账户列表（可能为空）。
	ListByUserID(userID uint64) ([]models.Account, error)

	// SetUpdateMaxTry 设置乐观锁更新的最大重试次数。
	// 用于调整并发冲突时的重试策略。
	SetUpdateMaxTry(maxTry int)

	// CheckOwnership 检查用户是否对指定账户有所有权（即通过 account_users 关联）。
	// 返回 true 表示存在关联，false 表示无关联或账户不存在。
	CheckOwnership(userID, accountID uint64) (bool, error)
}
