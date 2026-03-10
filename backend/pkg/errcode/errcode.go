package errcode

import "net/http"

const (
	// 成功
	Success = 0
)

// 认证相关 (20001-20099)
const (
	StatusInvalidToken            = 20001 // 表示提供的令牌无效或已过期
	StatusMissedToken             = 20002 // 表示请求中缺少必要的令牌
	StatusUserNotAuthenticated    = 20003 // 表示用户尚未通过身份认证
	StatusInvalidUserData         = 20004 // 表示提供的用户数据格式不正确或无效
	StatusInsufficientPermissions = 20005 // 表示当前令牌的权限不足以执行该操作
)

// 通用业务错误 (40001-40099)
const (
	StatusInvalidParams    = 40001 // 表示请求参数无效或缺失
	StatusUnauthorized     = 40002 // 表示请求未获得授权
	StatusNotFound         = 40004 // 表示请求的资源不存在
	StatusDuplicate        = 40005 // 表示尝试创建的资源已存在（重复数据）
	StatusInsufficientPerm = 40006 // 表示当前用户权限不足（与 StatusInsufficientPermissions 类似，但用于不同场景）
)

// 用户相关(40100-40199)
const (
	StatusRegisterNameExist        = 40101 // 表示注册时用户名已存在
	StatusLoginNameOrPasswordWrong = 40102 // 表示登录时用户名或密码错误
	StatusUserNotFound             = 40103 // 表示根据提供的标识未找到对应用户（注意拼写：NotFount 应为 NotFound，但为保持兼容保留原拼写）
	StatusInviteCodeNotFound       = 40104 // 表示提供的邀请码无效或不存在
	StatusInviteCodeUsed           = 40105 // 邀请码已被使用
)

// project 相关(40200-40299)
const (
	StatusProjectNotFound      = 40201 // 项目不存在
	StatusTaskNotFound         = 40202 // 任务不存在
	StatusBudgetNotFound       = 40203 // 预算项不存在
	StatusPrerequisiteNotFound = 40204 // 依赖关系不存在
)

// 服务器内部错误 (90000-90099)
const (
	StatusServerError   = 90000 // 表示服务器内部发生未预期的错误
	StatusDatabaseError = 90001 // 表示数据库操作失败
)

var codeMsgMap = map[int]string{
	Success:                        "success",
	StatusPrerequisiteNotFound:     "依赖关系不存在",
	StatusInviteCodeUsed:           "邀请码已被使用",
	StatusUserNotFound:             "用户不存在",
	StatusLoginNameOrPasswordWrong: "用户名或密码错误",
	StatusRegisterNameExist:        "用户名已存在",
	StatusInviteCodeNotFound:       "邀请码无效",
	StatusInvalidParams:            "请求参数错误",
	StatusUnauthorized:             "未认证",
	StatusNotFound:                 "资源不存在",
	StatusDuplicate:                "重复数据",
	StatusInsufficientPerm:         "权限不足",
	StatusServerError:              "服务器内部错误",
	StatusDatabaseError:            "数据库错误",
	StatusInvalidToken:             "无效的令牌",
	StatusInsufficientPermissions:  "Token权限不足",
	StatusMissedToken:              "缺少令牌",
	StatusUserNotAuthenticated:     "用户未认证",
	StatusInvalidUserData:          "无效的用户数据",
	StatusProjectNotFound:          "项目不存在",
	StatusTaskNotFound:             "任务不存在",
	StatusBudgetNotFound:           "预算项不存在",
}
var codeToHttpStatus = map[int]int{
	Success:                        http.StatusOK,
	StatusInvalidToken:             http.StatusUnauthorized,
	StatusMissedToken:              http.StatusUnauthorized,
	StatusUserNotAuthenticated:     http.StatusUnauthorized,
	StatusInvalidUserData:          http.StatusBadRequest,
	StatusInsufficientPermissions:  http.StatusForbidden,
	StatusInvalidParams:            http.StatusBadRequest,
	StatusUnauthorized:             http.StatusUnauthorized,
	StatusNotFound:                 http.StatusNotFound,
	StatusDuplicate:                http.StatusConflict,
	StatusInsufficientPerm:         http.StatusForbidden,
	StatusRegisterNameExist:        http.StatusConflict,
	StatusLoginNameOrPasswordWrong: http.StatusUnauthorized,
	StatusUserNotFound:             http.StatusNotFound,
	StatusInviteCodeNotFound:       http.StatusNotFound,
	StatusInviteCodeUsed:           http.StatusBadRequest,
	StatusProjectNotFound:          http.StatusNotFound,
	StatusTaskNotFound:             http.StatusNotFound,
	StatusBudgetNotFound:           http.StatusNotFound,
	StatusPrerequisiteNotFound:     http.StatusNotFound,
	StatusServerError:              http.StatusInternalServerError,
	StatusDatabaseError:            http.StatusInternalServerError,
}

// CodeMsg 根据错误码返回对应的描述信息，若未找到则返回“未知错误”
func CodeMsg(code int) string {
	msg, ok := codeMsgMap[code]
	if ok {
		return msg
	}
	return "未知错误"
}

func CodeHttpStatus(code int) int {
	sta, ok := codeToHttpStatus[code]
	if ok {
		return sta
	}
	return http.StatusInternalServerError
}
