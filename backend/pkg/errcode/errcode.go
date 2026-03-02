package errcode

const (
	// 成功
	Success = 0
)

// 认证相关 (20001-20099)
const (
	// StatusInvalidToken 表示提供的令牌无效或已过期
	StatusInvalidToken = 20001
	// StatusMissedToken 表示请求中缺少必要的令牌
	StatusMissedToken = 20002
	// StatusUserNotAuthenticated 表示用户尚未通过身份认证
	StatusUserNotAuthenticated = 20003
	// StatusInvalidUserData 表示提供的用户数据格式不正确或无效
	StatusInvalidUserData = 20004
	// StatusInsufficientPermissions 表示当前令牌的权限不足以执行该操作
	StatusInsufficientPermissions = 20005
)

// 通用业务错误 (40001-40099)
const (
	// StatusInvalidParams 表示请求参数无效或缺失
	StatusInvalidParams = 40001
	// StatusUnauthorized 表示请求未获得授权
	StatusUnauthorized = 40002
	// StatusNotFound 表示请求的资源不存在
	StatusNotFound = 40004
	// StatusDuplicate 表示尝试创建的资源已存在（重复数据）
	StatusDuplicate = 40005
	// StatusInsufficientPerm 表示当前用户权限不足（与 StatusInsufficientPermissions 类似，但用于不同场景）
	StatusInsufficientPerm = 40006
)

// 用户相关(40100-40199)
const (
	// StatusRegisterNameExist 表示注册时用户名已存在
	StatusRegisterNameExist = 40101
	// StatusLoginNameOrPasswordWrong 表示登录时用户名或密码错误
	StatusLoginNameOrPasswordWrong = 40102
	// StatusUserNotFount 表示根据提供的标识未找到对应用户（注意拼写：NotFount 应为 NotFound，但为保持兼容保留原拼写）
	StatusUserNotFount = 40103
	// StatusInviteCodeNotFound 表示提供的邀请码无效或不存在
	StatusInviteCodeNotFound = 40104
	//StatusInviteCodeUsed 邀请码已被使用
	StatusInviteCodeUsed = 40105
)

// 服务器内部错误 (90000-90099)
const (
	// StatusServerError 表示服务器内部发生未预期的错误
	StatusServerError = 90000
	// StatusDatabaseError 表示数据库操作失败
	StatusDatabaseError = 90001
	// 可根据需要添加更多
)

var codeMsgMap = map[int]string{
	Success:                        "success",
	StatusInviteCodeUsed:           "邀请码已被使用",
	StatusUserNotFount:             "用户不存在",
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
}

// CodeMsg 根据错误码返回对应的描述信息，若未找到则返回“未知错误”
func CodeMsg(code int) string {
	msg, ok := codeMsgMap[code]
	if ok {
		return msg
	}
	return "未知错误"
}
