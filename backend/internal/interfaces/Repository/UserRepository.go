package Repository

import (
	"LifeNavigator/internal/models"
)

// UserRepository 定义了用户数据访问接口。
type UserRepository interface {
	// Create 创建新用户。
	// 错误：如果用户名已存在返回 ErrRecordExist。
	Create(user *models.User) error

	// GetByID 根据主键 ID 查询用户。
	// 错误：如果用户不存在返回 ErrNotFound。
	GetByID(id uint64) (*models.User, error)

	// GetByUsername 根据用户名查询用户。
	// 错误：如果用户不存在返回 ErrNotFound。
	GetByUsername(username string) (*models.User, error)

	// GetByEmail 根据邮箱查询用户。
	// 错误：如果用户不存在返回 ErrNotFound。
	GetByEmail(email string) (*models.User, error)

	// Update 更新用户信息（全量更新）。
	// 错误：如果用户不存在返回 ErrNotFound。
	Update(user *models.User) error

	// UpdateProfile 更新用户的个人简介。
	// 实现行为：如果 profile 长度超过 2000 字符，则自动截断至 2000。
	// 错误：如果用户不存在返回 ErrNotFound。
	UpdateProfile(userID uint64, profile string) error

	// Delete 软删除用户（设置 deleted_at）。
	// 错误：如果用户不存在返回 ErrNotFound。
	Delete(id uint64) error

	// List 分页列出所有未软删除的用户。
	List(offset, limit int) ([]models.User, error)

	// HardDeleteById 永久删除用户（物理删除），不可逆。
	// 错误：如果用户不存在返回 ErrNotFound。
	HardDeleteById(id uint64) error

	// SoftDeleteById 软删除用户（设置 deleted_at）。
	// 错误：如果用户不存在返回 ErrNotFound。
	SoftDeleteById(id uint64) error

	// UpdateAvatar 更新用户头像 URL。
	// 错误：如果用户不存在返回 ErrNotFound。
	UpdateAvatar(userID uint64, avatarURL string) error
}
