// internal/interfaces/serviceInte/inviteCodeService.go
package Service

import (
	"LifeNavigator/internal/models"
)

// InviteCodeService 定义邀请码管理接口
type InviteCodeService interface {
	// GenerateInviteCode 生成邀请码
	// 参数:
	//   amount: 邀请码可邀请的次数（通常为1，表示单次使用）
	//   currentUserID: 当前用户ID，作为创建者标识
	//   role: 邀请码授予的角色（如 "user", "admin"）
	// 返回值:
	//   *models.InviteCode: 生成的邀请码对象，包含唯一的 Token
	//   error: 可能返回 ErrInternal
	GenerateInviteCode(amount int, currentUserID uint64, role string) (*models.InviteCode, error)

	// UseCode 使用邀请码
	// 参数:
	//   token: 邀请码的字符串标识
	// 返回值:
	//   error: 可能返回 ErrInviteCodeNotFound（邀请码不存在）、ErrInviteCodeUsed（已使用）、ErrInternal
	UseCode(token string) error

	// GetCodeInfo 获取邀请码详情
	// 参数:
	//   token: 邀请码Token
	//   currentUserID: 当前用户ID，用于权限校验（只有创建者可以查看）
	// 返回值:
	//   *models.InviteCode: 邀请码详情
	//   error: 可能返回 ErrInviteCodeNotFound、ErrForbidden、ErrInternal
	GetCodeInfo(token string, currentUserID uint64) (*models.InviteCode, error)

	// ListCodesByUser 列出当前用户创建的所有邀请码
	// 参数:
	//   currentUserID: 当前用户ID
	//   offset: 分页偏移量
	//   limit: 每页数量
	// 返回值:
	//   []models.InviteCode: 邀请码列表
	//   error: 可能返回 ErrInternal
	ListCodesByUser(currentUserID uint64, offset, limit int) ([]models.InviteCode, error)
}
