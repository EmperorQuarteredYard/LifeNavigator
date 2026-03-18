package service

import "errors"

// 用户相关
var (
	ErrUserInfoIncomplete = errors.New("用户信息不完整")
	ErrUserNameExists     = errors.New("用户名已存在")
	ErrUserNotFound       = errors.New("用户名不存在")
	ErrPasswordWrong      = errors.New("密码错误")
)

// 邀请码相关
var (
	ErrInviteCodeNotFound = errors.New("邀请码不存在")
	ErrInvalidToken       = errors.New("Token 不正确")
	ErrInviteCodeUsed     = errors.New("邀请码已被使用")
)

// 账户相关
var (
	ErrAccountNotFound = errors.New("账户不存在")
)

// 项目相关
var (
	ErrProjectNotFound        = errors.New("项目不存在")
	ErrTaskNotFound           = errors.New("任务不存在")
	ErrTaskDependencyNotFound = errors.New("任务依赖不存在")
	ErrTaskBudgetNotFound     = errors.New("任务预算不存在")
	ErrProjectBudgetNotFound  = errors.New("项目预算不存在")
	ErrBudgetNotFound         = errors.New("预算项不存在")
	ErrKanbanNotFound         = errors.New("看板不存在")
)

// 通用
var (
	ErrInternal         = errors.New("内部错误")
	ErrInvalidInput     = errors.New("无效输入")
	ErrDuplicate        = errors.New("重复记录")
	ErrConcurrentUpdate = errors.New("并发读写出错")
	ErrForbidden        = errors.New("权限不足")
)
