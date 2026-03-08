package controller

import (
	"LifeNavigator/internal/service"
	"LifeNavigator/pkg/errcode"
	"LifeNavigator/pkg/jwt"
	"LifeNavigator/pkg/response"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

var codeToHttpStatus = map[int]int{
	errcode.Success:                        http.StatusOK,
	errcode.StatusInvalidToken:             http.StatusUnauthorized,
	errcode.StatusMissedToken:              http.StatusUnauthorized,
	errcode.StatusUserNotAuthenticated:     http.StatusUnauthorized,
	errcode.StatusInvalidUserData:          http.StatusBadRequest,
	errcode.StatusInsufficientPermissions:  http.StatusForbidden,
	errcode.StatusInvalidParams:            http.StatusBadRequest,
	errcode.StatusUnauthorized:             http.StatusUnauthorized,
	errcode.StatusNotFound:                 http.StatusNotFound,
	errcode.StatusDuplicate:                http.StatusConflict,
	errcode.StatusInsufficientPerm:         http.StatusForbidden,
	errcode.StatusRegisterNameExist:        http.StatusConflict,
	errcode.StatusLoginNameOrPasswordWrong: http.StatusUnauthorized,
	errcode.StatusUserNotFound:             http.StatusNotFound,
	errcode.StatusInviteCodeNotFound:       http.StatusNotFound,
	errcode.StatusInviteCodeUsed:           http.StatusBadRequest,
	errcode.StatusProjectNotFound:          http.StatusNotFound,
	errcode.StatusTaskNotFound:             http.StatusNotFound,
	errcode.StatusBudgetNotFound:           http.StatusNotFound,
	errcode.StatusPrerequisiteNotFound:     http.StatusNotFound,
	errcode.StatusServerError:              http.StatusInternalServerError,
	errcode.StatusDatabaseError:            http.StatusInternalServerError,
}

func getHttpStatus(code int) int {
	if httpStatus, ok := codeToHttpStatus[code]; ok {
		return httpStatus
	}
	return http.StatusInternalServerError
}

// BindJSON 绑定 JSON 并自动处理错误
func (b *BaseController) BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		if err == io.EOF {
			log.Printf("Empty request body from %s", c.ClientIP())
		} else {
			log.Printf("JSON bind error from %s: %v", c.ClientIP(), err)
		}
		b.HandleCode(c, errcode.StatusInvalidParams)
		return false
	}
	return true
}

// GetAuthUser 获取当前认证用户
func (b *BaseController) GetAuthUser(c *gin.Context) (*jwt.AuthUser, bool) {
	val, exists := c.Get("user")
	if !exists {
		b.HandleCode(c, errcode.StatusUnauthorized)
		return nil, false
	}
	user, ok := val.(*jwt.AuthUser)
	if !ok {
		b.HandleCode(c, errcode.StatusInvalidUserData)
		return nil, false
	}
	return user, true
}

func (b *BaseController) ServerError(c *gin.Context) {
	b.HandleCode(c, errcode.StatusServerError)
}

// 为保证性能，请视情况使用
func (b *BaseController) Error(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		b.HandleCode(c, errcode.StatusUserNotFound)
	case errors.Is(err, service.ErrUserNameExists):
		b.HandleCode(c, errcode.StatusRegisterNameExist)
	case errors.Is(err, service.ErrPasswordWrong):
		b.HandleCode(c, errcode.StatusLoginNameOrPasswordWrong)
	case errors.Is(err, service.ErrForbidden):
		b.HandleCode(c, errcode.StatusInsufficientPerm)
	case errors.Is(err, service.ErrUserInfoIncomplete):
		b.HandleCode(c, errcode.StatusInvalidUserData)
	case errors.Is(err, service.ErrInviteCodeNotFound):
		b.HandleCode(c, errcode.StatusInviteCodeNotFound)
	case errors.Is(err, service.ErrInviteCodeUsed):
		b.HandleCode(c, errcode.StatusInviteCodeUsed)
	case errors.Is(err, service.ErrInvalidToken):
		b.HandleCode(c, errcode.StatusInvalidToken)
	case errors.Is(err, service.ErrProjectNotFound):
		b.HandleCode(c, errcode.StatusProjectNotFound)
	case errors.Is(err, service.ErrTaskNotFound):
		b.HandleCode(c, errcode.StatusTaskNotFound)
	case errors.Is(err, service.ErrTaskDependencyNotFound):
		b.HandleCode(c, errcode.StatusPrerequisiteNotFound)
	case errors.Is(err, service.ErrTaskBudgetNotFound) ||
		errors.Is(err, service.ErrProjectBudgetNotFound) ||
		errors.Is(err, service.ErrBudgetNotFound):
		b.HandleCode(c, errcode.StatusBudgetNotFound)
	case errors.Is(err, service.ErrInternal):
		b.HandleCode(c, errcode.StatusServerError)
	case errors.Is(err, service.ErrInvalidInput):
		b.HandleCode(c, errcode.StatusInvalidParams)
	case errors.Is(err, service.ErrDuplicate):
		b.HandleCode(c, errcode.StatusDuplicate)
	default:
		log.Printf("service 层有未映射的错误：%v", err)
		b.HandleCode(c, errcode.StatusServerError)
	}
}

// Success 成功响应（带数据）
func (b *BaseController) Success(c *gin.Context, data interface{}) {
	response.Success(c, data)
}

func (b *BaseController) HandleCode(c *gin.Context, code int) {
	response.HandleCode(c, getHttpStatus(code), code)
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
