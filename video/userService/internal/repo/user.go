package repo

import (
	"context"
	userData "userService/internal/data"
)

type UserRepo interface {
	FindUserByUserName(ctx context.Context, username string) (u *userData.User, err error)
	Create(ctx context.Context, user *userData.User) error
	FindUserById(ctx context.Context, uid int64) (user *userData.User, err error)
	UpdateAvatar(ctx context.Context, id int64, s string) error
	SaveMFASecret(ctx context.Context, id int64, qcode string) error
}
