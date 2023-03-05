package rpc

import "byte_dance_5th/pkg/ymlconfig"

func InitRpc(Congfig *ymlconfig.Config) {
	UserConfig := ymlconfig.ConfigInit("ByteDance_USER", "user_config")
	InitUserRpc(&UserConfig)
}
