package dao

import (
	"context"
	"gorm.io/gorm"
	userData "userService/internal/data"
	"userService/internal/database/gorms"
)

func NewUserDao() *UserDao {
	return &UserDao{
		conn: gorms.New(),
	}
}

type UserDao struct {
	conn *gorms.GormConn
}

func (u *UserDao) SetIsSecretEnabled(ctx context.Context, id int64) error {
	err := u.conn.Session(ctx).Model(&userData.User{}).Where("id=?", id).Update("is_mfa_enabled", 1).Error
	return err
}

func (u *UserDao) FindSecretById(ctx context.Context, id int64) (string, error) {
	secret := ""
	err := u.conn.Session(ctx).Model(&userData.User{}).Select("mfa_secret").Where("id=?", id).First(&secret).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	return secret, err
}

func (u *UserDao) SaveMFASecret(ctx context.Context, id int64, qcode string) error {
	err := u.conn.Session(ctx).Model(&userData.User{}).Where("id=?", id).Update("mfa_secret", qcode).Error
	return err
}

func (u *UserDao) UpdateAvatar(ctx context.Context, id int64, url string) error {
	err := u.conn.Session(ctx).Model(&userData.User{}).Where("id=?", id).Update("avatar_url", url).Error
	return err
}

func (u *UserDao) FindUserById(ctx context.Context, uid int64) (user *userData.User, err error) {
	err = u.conn.Session(ctx).Where("id = ? ", uid).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}

func (u *UserDao) Create(ctx context.Context, user *userData.User) error {
	err := u.conn.Session(ctx).Create(user).Error
	return err
}

func (u *UserDao) FindUserByUserName(ctx context.Context, username string) (us *userData.User, err error) {
	err = u.conn.Session(ctx).Where("username = ?", username).First(&us).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return
}
