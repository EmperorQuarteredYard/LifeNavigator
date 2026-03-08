package response

import (
	"LifeNavigator/pkg/errcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success,
		Message: "success",
		Data:    data,
	})
}

func Code(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: errcode.CodeMsg(code),
		Data:    nil,
	})
}

func HandleCode(c *gin.Context, httpStatus, code int) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: errcode.CodeMsg(code),
		Data:    nil,
	})
}
