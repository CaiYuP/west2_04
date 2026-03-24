package model

import "west2-video/common/errs"

var (
	BindError                = errs.NewError(9999, "参数有误")
	DBError                  = errs.NewError(9888, "数据库出错")
	RedisError               = errs.NewError(9555, "redis出错")
	JWTError                 = errs.NewError(9777, "jwt出错")
	TokenError               = errs.NewError(9666, "token认证失败")
	RpcError                 = errs.NewError(9555, "获取Rpc失败")
	AccountExited            = errs.NewError(101, "账号已注册")
	AccountNotExist          = errs.NewError(102, "账号不存在")
	PassWordError            = errs.NewError(103, "密码错误")
	UserUnExist              = errs.NewError(104, "账号不存在错误")
	AvatarUploadError        = errs.NewError(105, "头像上传错误")
	CreateMfaError           = errs.NewError(106, "生成 MFA 密钥失败")
	CreateMfaCodeError       = errs.NewError(107, "生成 MFA 二维码失败")
	NotSetSecret             = errs.NewError(108, "mfa密钥未设置")
	SecretError              = errs.NewError(109, "mfa密钥与数据库不符")
	MFACodeInvalid           = errs.NewError(110, "mfa密钥不正确")
	UploadImgError           = errs.NewError(110, "上传图片至图床失败")
	SeedError                = errs.NewError(201, "时间戳超过当前时间")
	CantLikeMp4              = errs.NewError(202, "文件太小,无法判断是否为有效的MP4")
	MP4NoEffect              = errs.NewError(203, "不是有效的MP4文件（未找到ftyp签名）")
	MP4NoSafe                = errs.NewError(204, "MP4文件的主品牌不受支持")
	MinioError               = errs.NewError(205, "minio错误")
	FromTimeMustBigger       = errs.NewError(206, "终止时间比初始时间小")
	IsNotAuthor              = errs.NewError(301, "不是作者无法删除评论")
	ConnotDeleteOtherComment = errs.NewError(302, "不能删除别人的评论")
	IsAlreadyLike            = errs.NewError(303, "已经点赞过了")
	CannotFollowSelf         = errs.NewError(401, "不能关注自己")
	CannotCancelFollowSelf   = errs.NewError(402, "不能取关自己")
	ConnotRepeatFollow       = errs.NewError(403, "不能重复关注")
)

const (
	NoToken       = 10001
	TokenNoBearer = 10002
	TokenIsError  = 10003
)
