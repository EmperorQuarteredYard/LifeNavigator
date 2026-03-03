package controller

import (
	"LifeNavigator/pkg/errcode"
	"LifeNavigator/pkg/jwt"
	"LifeNavigator/pkg/response"

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

// Success 成功响应（带数据）
func (b *BaseController) Success(c *gin.Context, data interface{}) {
	response.Success(c, data)
}

func (b *BaseController) HandleCode(c *gin.Context, code int) {
	response.Code(c, code)
}
