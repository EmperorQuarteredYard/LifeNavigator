package service

import "errors"

var (
	ErrUserInfoNotFound   = errors.New("用户信息不完整")
	ErrUserNameExists     = errors.New("用户名已存在")
	ErrUserNotFound       = errors.New("用户名不存在")
	ErrPasswordWrong      = errors.New("密码错误")
	ErrInternal           = errors.New("内部错误")
	ErrInviteCodeNotFound = errors.New("邀请码不存在")
	ErrInvalidToken       = errors.New("Token不正确")
)
