package domain

import (
	"context"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
	"strconv"
	"time"
	"userService/config"
	"userService/internal/dao"
	userData "userService/internal/data"
	"userService/internal/repo"
	"west2-video/common/encrypts"
	"west2-video/common/errs"
	"west2-video/common/jwts"
	"west2-video/common/logs"
	"west2-video/common/mfa"
	"west2-video/common/upload"
	"west2-video/gateway/biz/model"
)

func NewUserDomain() *UserDomain {
	return &UserDomain{
		userRepo: dao.NewUserDao(),
	}
}

type UserDomain struct {
	userRepo repo.UserRepo
}

func (d *UserDomain) Register(ctx context.Context, username, password string) (*userData.User, *errs.BError) {
	userByUn, err := d.userRepo.FindUserByUserName(ctx, username)
	if err != nil {
		logs.LG.Error("UserDomain FindUserByUserName FindUserByUserName error", zap.Error(err))
		return nil, model.DBError
	}
	if userByUn != nil {
		return nil, model.AccountExited
	}
	user := &userData.User{
		Username:  username,
		Password:  encrypts.Md5(password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = d.userRepo.Create(ctx, user)
	if err != nil {
		logs.LG.Error("UserDomain Register Create error", zap.Error(err))
		return nil, model.DBError
	}
	return user, nil
}

func (d *UserDomain) Login(ctx context.Context, username string, password string) (*userData.User, *errs.BError) {
	//判断账号是否存在
	u, err := d.userRepo.FindUserByUserName(ctx, username)
	if err != nil {
		logs.LG.Error("UserDomain FindUserByUserName error", zap.Error(err))
		return nil, model.DBError
	}
	if u == nil {
		return nil, model.UserUnExist
	}
	if u.Password != encrypts.Md5(password) {
		return nil, model.PassWordError
	}
	return u, nil
}

func (d *UserDomain) CreateToken(ctx context.Context, id uint64, username string, ip string) (*jwts.JwtToken, *errs.BError) {
	userIdStr := strconv.FormatInt(int64(id), 10)
	token, err := jwts.CreateToken(userIdStr, username, time.Duration(config.C.Jc.AccessExp), time.Duration(config.C.Jc.RefreshExp), config.C.Jc.AccessSecret, config.C.Jc.RefreshSecret, ip)
	if err != nil {
		logs.LG.Error("UserDomain CreateToken error", zap.Error(err))
		return nil, model.JWTError
	}
	return token, nil

}

func (d *UserDomain) RefreshToken(refreshToken string, ip string) (*jwts.JwtToken, *errs.BError) {
	token, err := jwts.RefreshToken(refreshToken, time.Duration(config.C.Jc.AccessExp), time.Duration(config.C.Jc.RefreshExp), config.C.Jc.AccessSecret, config.C.Jc.RefreshSecret, ip)
	if err != nil {
		logs.LG.Error("UserDomain RefreshToken error", zap.Error(err))
		return nil, model.JWTError
	}
	return token, nil
}

func (d *UserDomain) FindUserById(ctx context.Context, uid int64) (*userData.User, *errs.BError) {
	userById, err := d.userRepo.FindUserById(ctx, uid)
	if err != nil {
		logs.LG.Error("UserDomain FindUserById error", zap.Error(err))
		return nil, model.DBError
	}
	return userById, nil
}

func (d *UserDomain) UpdateAvatar(ctx context.Context, id int64, url string) *errs.BError {
	err := d.userRepo.UpdateAvatar(ctx, id, url)
	if err != nil {
		logs.LG.Error("UserDomain UpdateAvatar error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *UserDomain) CreateMfaQrcode(ctx context.Context, id int64, username string) (string, string, *errs.BError) {
	// 1. 生成 secret（用用户名做 account）
	secret, err := mfa.GenerateSecret(username)
	if err != nil {
		logs.LG.Error("UserDomain CreateMfaQrcode error", zap.Error(err))
		return "", "", model.CreateMfaError
	}

	// 2. 构造 otpauth://totp/... URL
	// issuer: 在 App 里显示的应用名，如 "west2"
	// account: 一般用用户名或邮箱，req 里如果有就用 req.Username/Email
	issuer := "west2"

	otpURL := mfa.BuildOtpAuthURL(secret, username, issuer)

	// 3. 生成二维码 data:image/png;base64,...
	qrcodeDataURL, err := mfa.GenerateMFAQRCodeDataURL(otpURL)
	if err != nil {
		logs.LG.Error("UserDomain CreateMfaQrcodeDataURL error", zap.Error(err))
		return "", "", model.CreateMfaCodeError
	}
	return secret, qrcodeDataURL, nil

}

func (d *UserDomain) SaveMFASecret(ctx context.Context, id int64, qcode string) *errs.BError {
	err := d.userRepo.SaveMFASecret(ctx, id, qcode)
	if err != nil {
		logs.LG.Error("UserDomain SaveMFASecret error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *UserDomain) FindSecretById(ctx context.Context, id int64) (string, *errs.BError) {
	secret, err := d.userRepo.FindSecretById(ctx, id)
	if err != nil {
		logs.LG.Error("UserDomain FindSecretById error", zap.Error(err))
		return "", model.DBError
	}
	if secret == "" {
		return "", model.NotSetSecret
	}
	return secret, nil
}

func (d *UserDomain) VerifySecret(ctx context.Context, secret string, newSecret, code string) *errs.BError {
	if secret != newSecret {
		return model.SecretError
	}
	valid := totp.Validate(code, secret)
	if !valid {
		return model.MFACodeInvalid
	}
	return nil
}

func (d *UserDomain) UpdateIsSecretEnabled(ctx context.Context, id int64) *errs.BError {
	err := d.userRepo.SetIsSecretEnabled(ctx, id)
	if err != nil {
		logs.LG.Error("UserDomain UpdateIsSecretEnabled error", zap.Error(err))
		return model.DBError
	}
	return nil
}

func (d *UserDomain) CreatePngUrl(ctx context.Context, data []byte, idStr string) (string, *errs.BError) {
	image, err := upload.Uploader.UploadImage(data, idStr)
	if err != nil {
		logs.LG.Error("UserDomain CreatePngUrl error", zap.Error(err))
		return "", model.UploadImgError
	}
	return image, nil
}
