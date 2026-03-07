package controller

import (
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/errcode"
	"LifeNavigator/pkg/jwt"
	"LifeNavigator/pkg/response"
	"errors"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// BindJSON 绑定 JSON 并自动处理错误
func (b *BaseController) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.Code(c, errcode.StatusInvalidParams)
		return false
	}
	return true
}

// GetAuthUser 获取当前认证用户
func (b *BaseController) GetAuthUser(c *gin.Context) (*jwt.AuthUser, bool) {
	val, exists := c.Get("user")
	if !exists {
		response.Code(c, errcode.StatusUnauthorized)
		return nil, false
	}
	user, ok := val.(jwt.AuthUser)
	if !ok {
		response.Code(c, errcode.StatusInvalidUserData)
		return nil, false
	}
	return &user, true
}

func (b *BaseController) ServerError(c *gin.Context) {
	response.Code(c, errcode.StatusServerError)
}

// 为保证性能，请视情况使用
func (b *BaseController) Error(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		response.Code(c, errcode.StatusUserNotFound)
	case errors.Is(err, service.ErrUserNameExists):
		response.Code(c, errcode.StatusRegisterNameExist)
	case errors.Is(err, service.ErrPasswordWrong):
		response.Code(c, errcode.StatusLoginNameOrPasswordWrong)
	case errors.Is(err, service.ErrForbidden):
		response.Code(c, errcode.StatusInsufficientPerm)
	case errors.Is(err, service.ErrUserInfoIncomplete):
		response.Code(c, errcode.StatusInvalidUserData)
	case errors.Is(err, service.ErrInviteCodeNotFound):
		response.Code(c, errcode.StatusInviteCodeNotFound)
	case errors.Is(err, service.ErrInviteCodeUsed):
		response.Code(c, errcode.StatusInviteCodeUsed)
	case errors.Is(err, service.ErrInvalidToken):
		response.Code(c, errcode.StatusInvalidToken)
	case errors.Is(err, service.ErrProjectNotFound):
		response.Code(c, errcode.StatusProjectNotFound)
	case errors.Is(err, service.ErrTaskNotFound):
		response.Code(c, errcode.StatusTaskNotFound)
	case errors.Is(err, service.ErrTaskDependencyNotFound):
		response.Code(c, errcode.StatusPrerequisiteNotFound)
	case errors.Is(err, service.ErrTaskBudgetNotFound) ||
		errors.Is(err, service.ErrProjectBudgetNotFound) ||
		errors.Is(err, service.ErrBudgetNotFound):
		response.Code(c, errcode.StatusBudgetNotFound)
	case errors.Is(err, service.ErrInternal):
		response.Code(c, errcode.StatusServerError)
	case errors.Is(err, service.ErrInvalidInput):
		response.Code(c, errcode.StatusInvalidParams)
	case errors.Is(err, service.ErrDuplicate):
		response.Code(c, errcode.StatusDuplicate)
	default:
		log.Printf("service 层有未映射的错误：%v", err)
		response.Code(c, errcode.StatusServerError)
	}
}

// Success 成功响应（带数据）
func (b *BaseController) Success(c *gin.Context, data interface{}) {
	response.Success(c, data)
}

func (b *BaseController) Code(c *gin.Context, code int) {
	response.Code(c, code)
}

func (b *BaseController) parsePagination(c *gin.Context) (page, pageSize int) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 20
	}
	// 限制最大 pageSize 防止过载
	if pageSize > 100 {
		pageSize = 100
	}
	return
}
