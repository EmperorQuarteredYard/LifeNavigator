package Service

import (
	"LifeNavigator/internal/models"
)

// UserService 定义用户管理接口
type UserService interface {
	// Register 注册新用户
	// 参数:
	//   user: 用户模型，应包含用户名、密码等必要信息
	// 返回值:
	//   error: 可能返回 ErrUserInfoIncomplete、ErrUserNameExists、ErrInternal
	Register(user *models.User) error

	// Login 用户登录
	// 参数:
	//   username: 用户名
	//   password: 密码
	// 返回值:
	//   *models.User: 登录成功的用户信息（不包含密码）
	//   error: 可能返回 ErrUserNotFound、ErrPasswordWrong、ErrInternal
	Login(username, password string) (*models.User, error)

	// GetByID 根据用户ID获取用户信息
	// 参数:
	//   id: 要查询的用户ID
	//   currentUserID: 当前用户ID，用于权限校验（只允许查询自己）
	// 返回值:
	//   *models.User: 用户信息
	//   error: 可能返回 ErrForbidden、ErrUserNotFound、ErrInternal
	GetByID(id uint64, currentUserID uint64) (*models.User, error)

	// GetByUsername 根据用户名获取用户信息
	// 参数:
	//   username: 用户名
	//   currentUserID: 当前用户ID，用于权限校验（只允许查询自己）
	// 返回值:
	//   *models.User: 用户信息
	//   error: 可能返回 ErrForbidden、ErrUserNotFound、ErrInternal
	GetByUsername(username string, currentUserID uint64) (*models.User, error)

	// GetByEmail 根据邮箱获取用户信息
	// 参数:
	//   email: 邮箱
	//   currentUserID: 当前用户ID，用于权限校验（只允许查询自己）
	// 返回值:
	//   *models.User: 用户信息
	//   error: 可能返回 ErrForbidden、ErrUserNotFound、ErrInternal
	GetByEmail(email string, currentUserID uint64) (*models.User, error)

	// Update 更新用户信息
	// 参数:
	//   user: 更新后的用户模型（必须包含ID）
	//   currentUserID: 当前用户ID，用于权限校验（只允许更新自己）
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrUserNotFound、ErrInternal
	Update(user *models.User, currentUserID uint64) error

	// HardDeleteByID 硬删除用户（物理删除）
	// 参数:
	//   id: 要删除的用户ID
	//   currentUserID: 当前用户ID，用于权限校验（只允许删除自己）
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrUserNotFound、ErrInternal
	HardDeleteByID(id uint64, currentUserID uint64) error

	// SoftDeleteByID 软删除用户（逻辑删除）
	// 参数:
	//   id: 要删除的用户ID
	//   currentUserID: 当前用户ID，用于权限校验（只允许删除自己）
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrUserNotFound、ErrInternal
	SoftDeleteByID(id uint64, currentUserID uint64) error

	// RefreshToken 刷新访问令牌
	// 参数:
	//   refreshToken: 刷新令牌
	// 返回值:
	//   accessToken: 新的访问令牌
	//   newRefreshToken: 新的刷新令牌
	//   error: 可能返回 ErrInvalidToken、ErrUserNotFound、ErrInternal
	RefreshToken(refreshToken string) (accessToken, newRefreshToken string, err error)

	// UpdateAvatar 更新用户头像URL
	// 参数:
	//   userID: 用户ID
	//   avatarURL: 头像URL
	// 返回值:
	//   error: 可能返回 ErrUserNotFound、ErrInternal
	UpdateAvatar(userID uint64, avatarURL string) error
}
