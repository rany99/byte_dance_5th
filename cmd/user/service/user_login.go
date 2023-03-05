package service

import (
	"byte_dance_5th/kitex_gen/user"
	"byte_dance_5th/models"
	password2 "byte_dance_5th/pkg/password"
	"byte_dance_5th/util/jwt"
	"context"
)

const (
	MaxNameLen     = 32
	MaxPasswordLen = 32
	MinPasswordLen = 5
)

type UserLoginService struct {
	ctx context.Context
}

func NewUserLoginService(ctx context.Context) *UserLoginService {
	return &UserLoginService{ctx: ctx}
}

func (u *UserLoginService) Do(req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	var username = req.Username
	var password = password2.Md5(req.Password)
	var usr models.User
	err = models.NewLoginDao().UserLogin(username, password, &usr)
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(usr)
	if err != nil {
		return nil, err
	}

	return &user.UserLoginResponse{
		StatusCode: 0,
		UserId:     usr.Id,
		Token:      token,
	}, nil
}
