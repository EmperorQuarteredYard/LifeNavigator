package Service

import (
	"LifeNavigator/internal/models"
)

// InviteUserService 定义邀请注册相关接口
type InviteUserService interface {
	// InviteUser 通过邀请码邀请新用户注册
	// 该方法会先注册用户，然后消耗邀请码，若后续步骤失败则回滚用户创建。
	// 参数:
	//   user: 待注册的用户模型（应包含用户名、密码等）
	//   inviteCodeToken: 邀请码Token
	// 返回值:
	//   accessToken: 生成的访问令牌
	//   refreshToken: 生成的刷新令牌
	//   error: 可能返回 ErrUserInfoIncomplete、ErrUserNameExists、ErrInviteCodeNotFound、ErrInviteCodeUsed、ErrInternal
	InviteUser(user *models.User, inviteCodeToken string) (accessToken string, refreshToken string, err error)

	// CreateInviteCode 创建邀请码（供邀请者使用）
	// 参数:
	//   inviterID: 邀请者用户ID
	//   amount: 可邀请次数
	//   role: 邀请码授予的角色
	// 返回值:
	//   *models.InviteCode: 生成的邀请码对象
	//   error: 可能返回 ErrForbidden（权限不足）、ErrInternal
	CreateInviteCode(inviterID uint64, amount int, role string) (*models.InviteCode, error)
}
