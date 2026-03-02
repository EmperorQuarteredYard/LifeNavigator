package service

import (
	"LifeNavigator/backend/internal/models"
	"LifeNavigator/backend/internal/repository"
	"errors"

	"github.com/google/uuid"
)

// InviteCodeService 定义了邀请码相关的业务逻辑操作。
type InviteCodeService interface {
	// GenerateInviteCode 生成一个新的邀请码。
	// 参数：
	//   - amount: 邀请码可被使用的次数上限。
	//   - createdBy: 创建者ID（如用户ID）。
	// 返回：
	//   - *models.InviteCode: 生成的邀请码对象，包含 token、amount 等字段。
	//   - error: 如果生成失败（如数据库错误），返回相应错误。
	GenerateInviteCode(amount int, invitedBy string) (*models.InviteCode, error)

	// UseCode 使用指定的邀请码令牌。
	// 该操作是原子的，会检查邀请码是否可用并增加使用次数。
	// 返回的错误包括：
	//   - ErrInviteCodeNotFound: 令牌不存在。
	//   - ErrInviteCodeUsed: 邀请码已用完（使用次数已达上限）。
	//   - 其他数据库错误。
	UseCode(token string) error

	// GetCodeInfo 获取邀请码的详细信息。
	// 如果令牌不存在，返回 ErrInviteCodeNotFound。
	GetCodeInfo(token string) (*models.InviteCode, error)

	// ListCodesByUser 分页查询指定用户创建的邀请码。
	// 参数：
	//   - userID: 创建者ID。
	//   - offset: 偏移量，用于分页。
	//   - limit: 每页数量。
	// 返回邀请码列表和可能的错误。
	ListCodesByUser(userID uint64, offset, limit int) ([]models.InviteCode, error)
}

// inviteCodeService 是 InviteCodeService 接口的实现，封装了邀请码的业务逻辑。
type inviteCodeService struct {
	repo repository.InviteCodeRepository // 依赖仓储层接口
}

// NewInviteCodeService 创建一个新的 InviteCodeService 实例。
func NewInviteCodeService(repo repository.InviteCodeRepository) InviteCodeService {
	return &inviteCodeService{repo: repo}
}

// GenerateInviteCode 实现 InviteCodeService.GenerateInviteCode。
// 它生成一个唯一的 token（使用 UUID），并调用仓储层保存。
func (s *inviteCodeService) GenerateInviteCode(amount int, invitedBy string) (*models.InviteCode, error) {
	// 生成 token
	token := uuid.New().String()
	code := &models.InviteCode{
		Token:     token,
		Amount:    amount,
		Count:     0,
		InvitedBy: invitedBy,
	}
	if err := s.repo.Create(code); err != nil {
		return nil, err
	}
	return code, nil
}

// UseCode 实现 InviteCodeService.UseCode。
// 它直接调用仓储层的乐观锁方法 UseCodeByToken，并转换错误（如果需要）。
func (s *inviteCodeService) UseCode(token string) error {
	err := s.repo.UseCodeByToken(token)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrInviteCodeNotFound
	}
	if errors.Is(err, repository.ErrInviteCodeUsed) {
		return err
	}
	return err
}

// GetCodeInfo 实现 InviteCodeService.GetCodeInfo。
func (s *inviteCodeService) GetCodeInfo(token string) (*models.InviteCode, error) {
	code, err := s.repo.FindByToken(token)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInviteCodeNotFound
		}

		return nil, ErrInternal
	}
	return code, nil
}

// ListCodesByUser 实现 InviteCodeService.ListCodesByUser。
func (s *inviteCodeService) ListCodesByUser(userID uint64, offset, limit int) ([]models.InviteCode, error) {
	return s.repo.GetByUserID(userID, offset, limit)
}
