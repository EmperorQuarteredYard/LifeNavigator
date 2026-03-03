package service

import (
	"LifeNavigator/backend/internal/models"
	"log"
)

type InviteUserService interface {
	//InviteUser 邀请注册：需要知道谁在邀请（当前用户），以便在失败时回滚
	InviteUser(user *models.User, inviteCodeToken string, inviterID uint64) error
	//CreateInviteCode 创建邀请码：需要当前用户
	CreateInviteCode(inviter *models.User, amount int) (*models.InviteCode, error)
}

type inviteUserService struct {
	userServ UserService
	codeServ InviteCodeService
}

func NewInviteUserService(userServ UserService, codeServ InviteCodeService) InviteUserService {
	return &inviteUserService{
		userServ: userServ,
		codeServ: codeServ,
	}
}

func (s *inviteUserService) InviteUser(user *models.User, inviteCodeToken string, inviterID uint64) error {
	// 注册用户（无需传入当前用户ID，因为是新用户注册）
	if err := s.userServ.Register(user); err != nil {
		return err
	}
	// 使用邀请码（无需权限）
	if err := s.codeServ.UseCode(inviteCodeToken); err != nil {
		log.Printf("failed to use invite code: %v", err)
		// 回滚：删除刚创建的用户（需传入当前用户ID，即被删除的用户ID）
		if delErr := s.userServ.HardDeleteByID(user.ID, user.ID); delErr != nil {
			log.Printf("failed to rollback user creation: %v", delErr)
		}
		return err
	}
	return nil
}

func (s *inviteUserService) CreateInviteCode(inviter *models.User, amount int) (*models.InviteCode, error) {
	// 生成邀请码需要当前用户ID（inviter.ID）
	return s.codeServ.GenerateInviteCode(amount, inviter.ID)
}
