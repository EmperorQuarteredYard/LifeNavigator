package controller

import (
	"LifeNavigator/backend/internal/service"
	"LifeNavigator/backend/pkg/dto"
	"LifeNavigator/backend/pkg/errcode"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InviteController struct {
	BaseController
	inviteCodeService service.InviteCodeService
	inviteUserService service.InviteUserService
	userService       service.UserService
}

func NewInviteController(
	inviteCodeService service.InviteCodeService,
	inviteUserService service.InviteUserService,
	userService service.UserService,
) *InviteController {
	return &InviteController{
		inviteCodeService: inviteCodeService,
		inviteUserService: inviteUserService,
		userService:       userService,
	}
}

// CreateInviteCode 创建邀请码（仅零号用户可操作）
// TODO 后续应当放开
func (ctl *InviteController) CreateInviteCode(c *gin.Context) {
	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.UserID != 0 {
		ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		return
	}

	var req dto.CreateInviteCodeRequest
	if !ctl.BindJSON(c, &req) {
		return
	}

	// 获取完整用户信息（需要 Username）
	user, err := ctl.userService.FindByID(authUser.UserID, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			ctl.HandleCode(c, errcode.StatusUserNotFount)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}

	inviteCode, err := ctl.inviteUserService.CreateInviteCode(user, req.Amount)
	if err != nil {
		// CreateInviteCode 可能返回的错误：service.ErrInternal 等
		ctl.HandleCode(c, errcode.StatusServerError)
		return
	}

	ctl.Success(c, inviteCode)
}

// GetInviteCode 获取邀请码详情
func (ctl *InviteController) GetInviteCode(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}

	code, err := ctl.inviteCodeService.GetCodeInfo(token, authUser.UserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		case errors.Is(err, service.ErrInviteCodeNotFound):
			ctl.HandleCode(c, errcode.StatusInviteCodeNotFound)
		default:
			ctl.HandleCode(c, errcode.StatusServerError)
		}
		return
	}

	ctl.Success(c, code)
}

// ListUserInviteCodes 列出指定用户创建的邀请码（需本人或零号用户）
func (ctl *InviteController) ListUserInviteCodes(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusInvalidParams)
		return
	}

	authUser, ok := ctl.GetAuthUser(c)
	if !ok {
		return
	}
	if authUser.UserID != userID && authUser.UserID != 0 {
		//TODO 这零号用户后期绝对得挨削，不然完蛋(其实也没有那么严重哈哈)
		ctl.HandleCode(c, errcode.StatusInsufficientPermissions)
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	codes, err := ctl.inviteCodeService.ListCodesByUser(userID, offset, limit)
	if err != nil {
		ctl.HandleCode(c, errcode.StatusServerError)
		return
	}

	ctl.Success(c, codes)
}
