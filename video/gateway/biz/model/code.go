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
	CreateMfaCodeError = errs.NewError(107, "生成 MFA 二维码失败")
	NotSetSecret       = errs.NewError(108, "mfa密钥未设置")
	SecretError        = errs.NewError(109, "mfa密钥与数据库不符")
	MFACodeInvalid     = errs.NewError(110, "mfa密钥不正确")
	UploadImgError     = errs.NewError(110, "上传图片至图床失败")
	SeedError          = errs.NewError(201, "时间戳超过当前时间")
)

const (
	NoToken       = 10001
	TokenNoBearer = 10002
	TokenIsError  = 10003
)
