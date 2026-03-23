package controller

import (
	"LifeNavigator/internal/interfaces/Service"
	"LifeNavigator/middleWare/jwt"
	"LifeNavigator/pkg/errcode"
	"LifeNavigator/pkg/response"
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

// BindJSON 绑定 JSON 并自动处理错误
func (b *BaseController) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		if err == io.EOF {
			log.Printf("Empty request body from %s", c.ClientIP())
		} else {
			log.Printf("JSON bind error from %s: %v", c.ClientIP(), err)
		}
		b.Code(c, errcode.StatusInvalidParams)
		return false
	}
	return true
}

// GetAuthUser 获取当前认证用户
func (b *BaseController) GetAuthUser(c *gin.Context) (*jwt.AuthUser, bool) {
	val, exists := c.Get("user")
	if !exists {
		b.Code(c, errcode.StatusUnauthorized)
		return nil, false
	}
	user, ok := val.(*jwt.AuthUser)
	if !ok {
		b.Code(c, errcode.StatusInvalidUserData)
		return nil, false
	}
	return user, true
}

func (b *BaseController) ServerError(c *gin.Context) {
	b.Code(c, errcode.StatusServerError)
}

// 为保证性能，请视情况使用
func (b *BaseController) Error(c *gin.Context, err error) {
	switch {
	case errors.Is(err, Service.ErrUserNotFound):
		b.Code(c, errcode.StatusUserNotFound)
	case errors.Is(err, Service.ErrUserNameExists):
		b.Code(c, errcode.StatusRegisterNameExist)
	case errors.Is(err, Service.ErrPasswordWrong):
		b.Code(c, errcode.StatusLoginNameOrPasswordWrong)
	case errors.Is(err, Service.ErrForbidden):
		b.Code(c, errcode.StatusInsufficientPerm)
	case errors.Is(err, Service.ErrUserInfoIncomplete):
		b.Code(c, errcode.StatusInvalidUserData)
	case errors.Is(err, Service.ErrInviteCodeNotFound):
		b.Code(c, errcode.StatusInviteCodeNotFound)
	case errors.Is(err, Service.ErrInviteCodeUsed):
		b.Code(c, errcode.StatusInviteCodeUsed)
	case errors.Is(err, Service.ErrInvalidToken):
		b.Code(c, errcode.StatusInvalidToken)
	case errors.Is(err, Service.ErrProjectNotFound):
		b.Code(c, errcode.StatusProjectNotFound)
	case errors.Is(err, Service.ErrTaskNotFound):
		b.Code(c, errcode.StatusTaskNotFound)
	case errors.Is(err, Service.ErrKanbanNotFound):
		b.Code(c, errcode.StatusNotFound)
	case errors.Is(err, Service.ErrTaskDependencyNotFound):
		b.Code(c, errcode.StatusPrerequisiteNotFound)
	case errors.Is(err, Service.ErrTaskBudgetNotFound) ||
		errors.Is(err, Service.ErrProjectBudgetNotFound) ||
		errors.Is(err, Service.ErrBudgetNotFound):
		b.Code(c, errcode.StatusBudgetNotFound)
	case errors.Is(err, Service.ErrInternal):
		b.Code(c, errcode.StatusServerError)
	case errors.Is(err, Service.ErrInvalidInput):
		b.Code(c, errcode.StatusInvalidParams)
	case errors.Is(err, Service.ErrDuplicate):
		b.Code(c, errcode.StatusDuplicate)
	default:
		log.Printf("service 层有未映射的错误：%v", err)
		b.Code(c, errcode.StatusServerError)
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
