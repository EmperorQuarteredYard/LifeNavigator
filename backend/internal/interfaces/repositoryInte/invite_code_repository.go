package repositoryInte

import (
	"LifeNavigator/internal/models"
)

// InviteCodeRepository 定义了邀请码相关的数据访问接口。
type InviteCodeRepository interface {
	// Create 创建一个新的邀请码记录。
	// 成功返回 nil，失败返回错误（如数据库错误）。
	Create(code *models.InviteCode) error

	// UseCodeByToken 根据令牌使用邀请码（增加使用次数）。
	// 参数 token: 邀请码的令牌字符串。
	// 实现行为：
	//   - 原子性地将 count 字段加1，前提是 count < amount。
	//   - 如果更新影响行数为0，会再次查询确认原因：
	//       若记录不存在返回 ErrNotFound；
	//       若 count >= amount 返回 ErrInviteCodeUsed。
	UseCodeByToken(token string) error

	// Delete 硬删除指定的邀请码记录。
	// 参数 code: 要删除的邀请码对象（其 ID 需有效）。
	// 错误：如果记录不存在返回 ErrNotFound。
	Delete(code *models.InviteCode) error

	// FindByToken 根据令牌查询邀请码。
	// 错误：若不存在返回 ErrNotFound。
	FindByToken(token string) (*models.InviteCode, error)

	// FindByID 根据主键 ID 查询邀请码。
	// 错误：若不存在返回 ErrNotFound。
	FindByID(id int64) (*models.InviteCode, error)

	// GetByUserID 分页查询某个用户创建的邀请码列表。
	// 参数 userID: 创建者 ID（invited_by 字段）。
	// 参数 offset, limit: 分页偏移量和限制条数。
	// 返回按创建时间倒序排列的邀请码列表。
	GetByUserID(uid uint64, offset, limit int) ([]models.InviteCode, error)
}
