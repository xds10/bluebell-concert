package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorInvalidPassword = errors.New("密码错误")
	ErrorInvalidID       = errors.New("无效id")
)
