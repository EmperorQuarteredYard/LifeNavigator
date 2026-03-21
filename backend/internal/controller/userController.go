package controller

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/service"
	"LifeNavigator/middleWare/jwt"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/errcode"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

// UploadAvatar 上传头像
func (ctl *UserController) UploadAvatar(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}

	const maxSize = 2 << 20 // 2MB
	if file.Size > maxSize {
		ctl.Code(c, errcode.StatusAvatarImageTooLarge)
		return
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}

	storageRoot := os.Getenv("AVATAR_STORAGE_PATH")
	if storageRoot == "" {
		storageRoot = "uploads/avatars"
	}
	// 确保目录存在
	if err := os.MkdirAll(storageRoot, 0755); err != nil {
		ctl.ServerError(c)
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	uniqueName := fmt.Sprintf("%d_%d%s", authUser.UserID, time.Now().UnixNano(), ext)
	savePath := filepath.Join(storageRoot, uniqueName)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		ctl.ServerError(c)
		return
	}

	// 生成 URL：相对路径，与静态文件服务挂载点一致
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", uniqueName)

	if err := ctl.userService.UpdateAvatar(authUser.UserID, avatarURL); err != nil {
		// 若更新失败，删除已上传的文件
		os.Remove(savePath)
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			ctl.Code(c, errcode.StatusUserNotFound)
		default:
			ctl.ServerError(c)
		}
		return
	}

	ctl.Success(c, gin.H{
		"avatar_url": avatarURL,
	})
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
			ctl.Code(c, errcode.StatusInvalidParams)
		case errors.Is(err, service.ErrUserNameExists):
			ctl.Code(c, errcode.StatusRegisterNameExist)
		case errors.Is(err, service.ErrInviteCodeNotFound):
			ctl.Code(c, errcode.StatusInviteCodeNotFound)
		case errors.Is(err, service.ErrInviteCodeUsed): // 需要处理仓储层返回的错误
			ctl.Code(c, errcode.StatusInviteCodeUsed)
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
			ctl.Code(c, errcode.StatusUserNotFound)
		case errors.Is(err, service.ErrPasswordWrong):
			ctl.Code(c, errcode.StatusLoginNameOrPasswordWrong)
		default:
			ctl.Code(c, errcode.StatusServerError)
		}
		return
	}

	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		ctl.Code(c, errcode.StatusServerError)
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
		ctl.Code(c, errcode.StatusInvalidParams)
		return
	}

	user, err := ctl.userService.GetByID(id, id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			ctl.Code(c, errcode.StatusUserNotFound)
		default:
			ctl.Code(c, errcode.StatusServerError)
		}
		return
	}

	ctl.Success(c, dto.UserProfile{
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		Avatar:    user.Avatar,
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
			ctl.Code(c, errcode.StatusUserNotFound)
		default:
			ctl.Code(c, errcode.StatusServerError)
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
			ctl.Code(c, errcode.StatusInvalidToken)
		case errors.Is(err, service.ErrUserNotFound):
			ctl.Code(c, errcode.StatusUserNotFound)
		default:
			ctl.Code(c, errcode.StatusServerError)
		}
		return
	}

	ctl.Success(c, dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
