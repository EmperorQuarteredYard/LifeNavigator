package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/pkg/jwt"
	"log"
)

type InviteUserService interface {
	//InviteUser 邀请注册：需要知道谁在邀请（当前用户），以便在失败时回滚
	InviteUser(user *models.User, inviteCodeToken string) (accessToken string, refreshToken string, err error)
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

func (s *inviteUserService) InviteUser(user *models.User, inviteCodeToken string) (accessToken string, refreshToken string, err error) {
	if err = s.userServ.Register(user); err != nil {
		return "", "", err
	}
	if err = s.codeServ.UseCode(inviteCodeToken); err != nil {
		log.Printf("failed to use invite code: %v", err)
		// 回滚：删除刚创建的用户（需传入当前用户ID，即被删除的用户ID）
		if delErr := s.userServ.HardDeleteByID(user.ID, user.ID); delErr != nil {
			log.Printf("failed to rollback user creation: %v", delErr)
		}
		return "", "", err
	}

	accessToken, refreshToken, err = jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		if delErr := s.userServ.HardDeleteByID(user.ID, user.ID); delErr != nil {
			log.Printf("failed to rollback user creation: %v", delErr)
		}
		return "", "", ErrInternal
	}
	return
}

func (s *inviteUserService) CreateInviteCode(inviter *models.User, amount int) (*models.InviteCode, error) {
	return s.codeServ.GenerateInviteCode(amount, inviter.ID)
}
