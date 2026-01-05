package model

import "west2-video/common/errs"

var (
	BindError          = errs.NewError(-999, "参数有误")
	DBError            = errs.NewError(-888, "数据库出错")
	RedisError         = errs.NewError(-555, "redis出错")
	JWTError           = errs.NewError(-777, "jwt出错")
	TokenError         = errs.NewError(-666, "token认证失败")
	AccountExited      = errs.NewError(101, "账号已注册")
	AccountNotExist    = errs.NewError(102, "账号不存在")
	PassWordError      = errs.NewError(103, "密码错误")
	UserUnExist        = errs.NewError(104, "账号不存在错误")
	AvatarUploadError  = errs.NewError(105, "头像上传错误")
	CreateMfaError     = errs.NewError(106, "生成 MFA 密钥失败")
	CreateMfaCodeError = errs.NewError(106, "生成 MFA 二维码失败")
)

const (
	NoToken       = 10001
	TokenNoBearer = 10002
	TokenIsError  = 10003
)
