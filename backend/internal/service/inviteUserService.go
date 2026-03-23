package service

import (
	"LifeNavigator/internal/interfaces/Service"
	"LifeNavigator/internal/models"
	"LifeNavigator/middleWare/jwt"
	"LifeNavigator/pkg/roles"
	"log"
)

type inviteUserService struct {
	userServ Service.UserService
	codeServ Service.InviteCodeService
}

func NewInviteUserService(userServ Service.UserService, codeServ Service.InviteCodeService) Service.InviteUserService {
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
		return "", "", Service.ErrInternal
	}
	return
}

func (s *inviteUserService) CreateInviteCode(inviterID uint64, amount int, role string) (*models.InviteCode, error) {
	user, err := s.userServ.GetByID(inviterID, inviterID)
	if err != nil {
		return nil, err
	}
	if roles.GetPrivilegeValue(user.Role) <= roles.GetPrivilegeValue(role) {
		return nil, Service.ErrForbidden
	}
	return s.codeServ.GenerateInviteCode(amount, inviterID, role)
}
