package errcode

import (
	"log"
	"net/http"
)

const (
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

// 用户相关 (40100-40199)
const (
	StatusRegisterNameExist        = 40101 // 表示注册时用户名已存在
	StatusLoginNameOrPasswordWrong = 40102 // 表示登录时用户名或密码错误
	StatusUserNotFound             = 40103 // 表示根据提供的标识未找到对应用户
	StatusInviteCodeNotFound       = 40104 // 表示提供的邀请码无效或不存在
	StatusInviteCodeUsed           = 40105 // 邀请码已被使用
	StatusAvatarImageTooLarge      = 40106 // 头像体积过大
)

// Project 相关 (40200-40299)
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

// codeMsgMap 错误码 -> 错误描述
var codeMsgMap map[int]string

// codeToHttpStatus 错误码 -> HTTP 状态码
var codeToHttpStatus map[int]int

func init() {
	// 初始化映射表
	codeMsgMap = make(map[int]string)
	codeToHttpStatus = make(map[int]int)

	// 注册成功状态
	register(Success, "success", http.StatusOK)

	// 认证相关
	register(StatusInvalidToken, "无效的令牌", http.StatusUnauthorized)
	register(StatusMissedToken, "缺少令牌", http.StatusUnauthorized)
	register(StatusUserNotAuthenticated, "用户未认证", http.StatusUnauthorized)
	register(StatusInvalidUserData, "无效的用户数据", http.StatusBadRequest)
	register(StatusInsufficientPermissions, "Token权限不足", http.StatusForbidden)

	// 通用业务错误
	register(StatusInvalidParams, "请求参数错误", http.StatusBadRequest)
	register(StatusUnauthorized, "未认证", http.StatusUnauthorized)
	register(StatusNotFound, "资源不存在", http.StatusNotFound)
	register(StatusDuplicate, "重复数据", http.StatusConflict)
	register(StatusInsufficientPerm, "权限不足", http.StatusForbidden)

	// 用户相关
	register(StatusRegisterNameExist, "用户名已存在", http.StatusConflict)
	register(StatusLoginNameOrPasswordWrong, "用户名或密码错误", http.StatusUnauthorized)
	register(StatusUserNotFound, "用户不存在", http.StatusNotFound)
	register(StatusInviteCodeNotFound, "邀请码无效", http.StatusNotFound)
	register(StatusInviteCodeUsed, "邀请码已被使用", http.StatusBadRequest)
	register(StatusAvatarImageTooLarge, "头像体积过大", http.StatusBadRequest)

	// Project 相关
	register(StatusProjectNotFound, "项目不存在", http.StatusNotFound)
	register(StatusTaskNotFound, "任务不存在", http.StatusNotFound)
	register(StatusBudgetNotFound, "预算项不存在", http.StatusNotFound)
	register(StatusPrerequisiteNotFound, "依赖关系不存在", http.StatusNotFound)

	// 服务器内部错误
	register(StatusServerError, "服务器内部错误", http.StatusInternalServerError)
	register(StatusDatabaseError, "数据库错误", http.StatusInternalServerError)
}

// register 将错误码、描述、HTTP状态码注册到映射表中
func register(code int, msg string, httpStatus int) {
	codeMsgMap[code] = msg
	codeToHttpStatus[code] = httpStatus
}

// CodeMsg 根据错误码返回对应的描述信息，若未找到则返回“未知错误”
func CodeMsg(code int) string {
	if msg, ok := codeMsgMap[code]; ok {
		return msg
	}
	return "未知错误"
}

// CodeHttpStatus 根据错误码返回对应的 HTTP 状态码，若未找到则返回 500
func CodeHttpStatus(code int) int {
	if status, ok := codeToHttpStatus[code]; ok {
		return status
	}
	log.Printf("错误码未正确注册:%d", code)
	return http.StatusBadRequest
}
