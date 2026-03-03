package service

import (
	"LifeNavigator/backend/internal/models"
	"LifeNavigator/backend/internal/repository"
	"errors"
	"log"
	"strconv"

	"github.com/google/uuid"
)

type InviteCodeService interface {
	//GenerateInviteCode 生成邀请码：需要传入当前用户ID，创建者标记为当前用户
	GenerateInviteCode(amount int, currentUserID uint64) (*models.InviteCode, error)

	//UseCode 使用邀请码：无需权限，任何人只要持有 token 即可使用
	UseCode(token string) error

	//GetCodeInfo 获取邀请码详情：只有创建者或管理员可查看（这里简化：仅创建者可查看）
	GetCodeInfo(token string, currentUserID uint64) (*models.InviteCode, error)

	//ListCodesByUser 列出当前用户创建的邀请码
	ListCodesByUser(currentUserID uint64, offset, limit int) ([]models.InviteCode, error)
}

type inviteCodeService struct {
	repo repository.InviteCodeRepository
}

func NewInviteCodeService(repo repository.InviteCodeRepository) InviteCodeService {
	return &inviteCodeService{repo: repo}
}

func (s *inviteCodeService) GenerateInviteCode(amount int, currentUserID uint64) (*models.InviteCode, error) {
	token := uuid.New().String()
	// 将用户ID转换为字符串作为创建者标识（示例格式，可按需调整）
	invitedBy := strconv.FormatUint(currentUserID, 10)
	code := &models.InviteCode{
		Token:     token,
		Amount:    amount,
		Count:     0,
		InvitedBy: invitedBy,
	}
	if err := s.repo.Create(code); err != nil {
		log.Printf("failed to create invite code: %v", err)
		return nil, ErrInternal
	}
	return code, nil
}

func (s *inviteCodeService) UseCode(token string) error {
	err := s.repo.UseCodeByToken(token)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrInviteCodeNotFound
	}
	if errors.Is(err, repository.ErrInviteCodeUsed) {
		return ErrErrInviteCodeUsed // 直接返回 repository 层的错误（也可包装）
	}
	if err != nil {
		log.Printf("failed to use invite code %s: %v", token, err)
		return ErrInternal
	}
	return nil
}

func (s *inviteCodeService) GetCodeInfo(token string, currentUserID uint64) (*models.InviteCode, error) {
	code, err := s.repo.FindByToken(token)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInviteCodeNotFound
		}
		log.Printf("failed to find invite code %s: %v", token, err)
		return nil, ErrInternal
	}
	// 校验：只有创建者才能查看详情
	if code.InvitedBy != strconv.FormatUint(currentUserID, 10) {
		return nil, ErrForbidden
	}
	return code, nil
}

func (s *inviteCodeService) ListCodesByUser(currentUserID uint64, offset, limit int) ([]models.InviteCode, error) {
	// 直接使用用户ID查询（repo 的 ListByUserID 接收 uint64）
	return s.repo.GetByUserID(currentUserID, offset, limit)
}
