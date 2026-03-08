package controller

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/errcode"
	"LifeNavigator/pkg/jwt"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
	userService   service.UserService
	inviteService service.InviteUserService
}

func NewUserController(userService service.UserService, inviteService service.InviteUserService) *UserController {
	return &UserController{
		userService:   userService,
		inviteService: inviteService,
	}
}

// Register 用户注册（需提供有效邀请码）
func (ctl *UserController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	accTok, refTok, err := ctl.inviteService.InviteUser(user, req.InviteCode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserInfoIncomplete):
			ctl.HandleCode(c, errcode.StatusInvalidParams)
		case errors.Is(err, service.ErrUserNameExists):
			ctl.HandleCode(c, errcode.StatusRegisterNameExist)
		case errors.Is(err, service.ErrInviteCodeNotFound):
			ctl.HandleCode(c, errcode.StatusInviteCodeNotFound)
		case errors.Is(err, service.ErrInviteCodeUsed): // 需要处理仓储层返回的错误
			ctl.HandleCode(c, errcode.StatusInviteCodeUsed)
		default:
			ctl.ServerError(c)
		}
		return
	}

	ctl.Success(c, dto.RegisterResponse{
		AccessToken:  accTok,
		RefreshToken: refTok,
		User: dto.UserProfile{
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Login 用户登录
func (ctl *UserController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	user, err := ctl.userService.Login(req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			ctl.HandleCode(c, errcode.StatusUserNotFound)
		case errors.Is(err, service.ErrPasswordWrong):
			ctl.HandleCode(c, errcode.StatusLoginNameOrPasswordWrong)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}

	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusServerError)
		return
	}

	ctl.Success(c, dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserProfile{
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
		},
	})
}

// GetUser 获取指定用户信息（公开）
func (ctl *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	user, err := ctl.userService.GetByID(id, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			ctl.HandleCode(c, errcode.StatusUserNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}

	ctl.Success(c, dto.UserProfile{
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	})
}

// Profile 获取当前登录用户信息
func (ctl *UserController) Profile(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	user, err := ctl.userService.GetByID(authUser.UserID, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			ctl.HandleCode(c, errcode.StatusUserNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}

	ctl.Success(c, dto.UserProfile{
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	})
}

// Refresh 刷新访问令牌
func (ctl *UserController) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	accessToken, refreshToken, err := ctl.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidToken):
			ctl.HandleCode(c, errcode.StatusInvalidToken)
		case errors.Is(err, service.ErrUserNotFound):
			ctl.HandleCode(c, errcode.StatusUserNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}

	ctl.Success(c, dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
