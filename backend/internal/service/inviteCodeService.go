package service

import (
	"LifeNavigator/internal/interfaces/Repository"
	"LifeNavigator/internal/interfaces/Service"
	"LifeNavigator/internal/models"
	"errors"
	"log"
	"strconv"

	"github.com/google/uuid"
)

type inviteCodeService struct {
	repo Repository.InviteCodeRepository
}

func NewInviteCodeService(repo Repository.InviteCodeRepository) Service.InviteCodeService {
	return &inviteCodeService{repo: repo}
}

func (s *inviteCodeService) GenerateInviteCode(amount int, currentUserID uint64, role string) (*models.InviteCode, error) {
	token := uuid.New().String()
	// 将用户 ID转换为字符串作为创建者标识
	invitedBy := strconv.FormatUint(currentUserID, 10)
	code := &models.InviteCode{
		Token:        token,
		Amount:       amount,
		Count:        0,
		InvitedBy:    invitedBy,
		InviteAsRole: role,
	}
	if err := s.repo.Create(code); err != nil {
		log.Printf("failed to create invite code: %v", err)
		return nil, Service.ErrInternal
	}
	return code, nil
}

func (s *inviteCodeService) UseCode(token string) error {
	err := s.repo.UseCodeByToken(token)
	if errors.Is(err, Repository.ErrNotFound) {
		return Service.ErrInviteCodeNotFound
	}
	if errors.Is(err, Repository.ErrInviteCodeUsed) {
		return Service.ErrInviteCodeUsed // 直接返回 Repository 层的错误（也可包装）
	}
	if err != nil {
		log.Printf("failed to use invite code %s: %v", token, err)
		return Service.ErrInternal
	}
	return nil
}

func (s *inviteCodeService) GetCodeInfo(token string, currentUserID uint64) (*models.InviteCode, error) {
	code, err := s.repo.FindByToken(token)
	if err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return nil, Service.ErrInviteCodeNotFound
		}
		log.Printf("failed to find invite code %s: %v", token, err)
		return nil, Service.ErrInternal
	}
	// 校验：只有创建者才能查看详情
	if code.InvitedBy != strconv.FormatUint(currentUserID, 10) {
		return nil, Service.ErrForbidden
	}
	return code, nil
}

func (s *inviteCodeService) ListCodesByUser(currentUserID uint64, offset, limit int) ([]models.InviteCode, error) {
	// 直接使用用户ID查询（repo 的 ListByUserID 接收 uint64）
	return s.repo.GetByUserID(currentUserID, offset, limit)
}
