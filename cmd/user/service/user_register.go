package service

import (
	"byte_dance_5th/kitex_gen/user"
	"byte_dance_5th/models"
	"byte_dance_5th/pkg/errortype"
	"byte_dance_5th/pkg/information"
	"byte_dance_5th/pkg/password"
	"byte_dance_5th/util/jwt"
	"context"
	"errors"
)

type UserRegisterService struct {
	ctx context.Context
}

func NewUserRegisterService(ctx context.Context) *UserRegisterService {
	return &UserRegisterService{ctx: ctx}
}

func (u *UserRegisterService) Do(req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	// 判断用户名是否已经被使用
	if models.NewLoginDao().UserAlreadyExist(req.Username) {
		return nil, errors.New(errortype.UserNameExistErr)
	}
	// 采用 MD5 对密码加密
	password := password.Md5(req.Password)
	// 准备好待填入的结构体
	usr := models.User{
		Username: req.Username,
		Password: password,
	}
	userInfo := models.UserInfo{
		Name:            req.Username,
		Avatar:          information.GetAvatarUrl(),
		BackgroundImage: information.GetBackGroundUrl(),
		Signature:       information.GetSignature(),
		User:            &usr,
	}

	err = models.NewUserInfoDAO().AddUserInfo(&userInfo)
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(usr)
	if err != nil {
		return nil, err
	}

	return &user.UserRegisterResponse{
		StatusCode: 0,
		UserId:     usr.Id,
		Token:      token,
	}, nil
}
