package rpc

import (
	"byte_dance_5th/kitex_gen/user"
	"byte_dance_5th/kitex_gen/user/userservice"
	"byte_dance_5th/pkg/middleware"
	"byte_dance_5th/pkg/ymlconfig"

	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
)

var userClient userservice.Client

func InitUserRpc(Config *ymlconfig.Config) {
	EtcdAddress := fmt.Sprintf("%s:%d", Config.Viper.GetString("Etcd.Address"), Config.Viper.GetInt("Etcd.Port"))
	// 服务发现
	r, err := etcd.NewEtcdResolver([]string{EtcdAddress})
	if err != nil {
		panic(err)
	}
	
	ServiceName := Config.Viper.GetString("Server.Name")

	c, err := userservice.NewClient(
		ServiceName,
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		client.WithMuxConnection(1),                       // mux
		client.WithRPCTimeout(3*time.Second),              // rpc timeout
		client.WithConnectTimeout(50*time.Millisecond),    // conn timeout
		client.WithFailureRetry(retry.NewFailurePolicy()), // retry
		client.WithSuite(trace.NewDefaultClientSuite()),   // tracer
		client.WithResolver(r),                            // resolver
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

func UserRegister(ctx context.Context, req *user.UserRegisterRequest) (resp *user.UserRegisterResponse, err error) {
	resp, err = userClient.UserRegister(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func UserLogin(ctx context.Context, req *user.UserLoginRequest) (resp *user.UserLoginResponse, err error) {
	resp, err = userClient.UserLogin(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
