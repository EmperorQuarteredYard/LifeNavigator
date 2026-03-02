package service

import (
	"LifeNavigator/backend/internal/models"
	"log"
	"strconv"
)

type InviteUserService interface {
	InviteUser(user *models.User, inviteCodeToken string) (err error)
	CreateInviteCode(user *models.User, amount int) (inviteCode *models.InviteCode, err error)
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

func (s *inviteUserService) InviteUser(user *models.User, inviteCodeToken string) (err error) {
	err = s.userServ.Register(user)
	if err != nil {
		return
	}
	err = s.codeServ.UseCode(inviteCodeToken)
	if err != nil {
		log.Printf("failed to invite user by code token : %v", err)
		err = s.userServ.HardDeleteById(user.ID)
		return
	}
	return
}

func (s *inviteUserService) CreateInviteCode(user *models.User, amount int) (inviteCode *models.InviteCode, err error) {
	inviteCode, err = s.codeServ.GenerateInviteCode(amount, strconv.FormatUint(user.ID, 10)+user.Username)
	if err != nil {
		return nil, err
	}
	return
}
