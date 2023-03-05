package handlers

import (
	"byte_dance_5th/cmd/api/rpc"
	"byte_dance_5th/kitex_gen/user"
	"byte_dance_5th/pkg/errortype"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
)

func UserRegister(ctx context.Context, c *app.RequestContext) {
	var p UserRegisterParam

	p.UserName = c.Query("username")
	p.PassWord = c.Query("password")

	if len(p.UserName) == 0 || len(p.PassWord) == 0 {
		msg := errortype.UserNameEmptyErr
		SendResponse(c, &user.UserRegisterResponse{
			StatusCode: 1,
			StatusMsg:  &msg,
		})
		return
	}

	resp, err := rpc.UserRegister(ctx, &user.UserRegisterRequest{
		Username: p.UserName,
		Password: p.PassWord,
	})
	if err != nil {
		msg := err.Error()
		SendResponse(c, &user.UserRegisterResponse{
			StatusCode: 1,
			StatusMsg:  &msg,
		})
	}

	SendResponse(c, resp)
}
