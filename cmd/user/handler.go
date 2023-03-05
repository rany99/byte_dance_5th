package main

import (
	"byte_dance_5th/cmd/user/service"
	"byte_dance_5th/kitex_gen/user"
	"context"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// UserRegister implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserRegister(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {

	if len(req.Username) == 0 || len(req.Password) == 0 {
		msg := "用户名和密码不能为空"
		return &user.UserRegisterResponse{StatusCode: 1, StatusMsg: &msg}, nil
	}
	resp, err = service.NewUserRegisterService(ctx).Do(req)

	if err != nil {
		msg := err.Error()
		return &user.UserRegisterResponse{
			StatusCode: 1,
			StatusMsg:  &msg,
		}, nil
	}

	return resp, err
}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {

	if len(req.Username) == 0 || len(req.Password) == 0 {
		msg := "用户名和密码不能为空"
		return &user.UserLoginResponse{
			StatusCode: 1,
			StatusMsg:  &msg,
		}, nil
	}

	resp, err = service.NewUserLoginService(ctx).Do(req)
	if err != nil {
		msg := err.Error()
		return &user.UserLoginResponse{
			StatusCode: 1,
			StatusMsg:  &msg,
		}, nil
	}

	return resp, nil
}
