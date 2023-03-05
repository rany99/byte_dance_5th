package main

import (
	user "byte_dance_5th/kitex_gen/user/userservice"
	"byte_dance_5th/pkg/ymlconfig"
	"fmt"
	"log"
)

var (
	Config      = ymlconfig.ConfigInit("ByteDance_USER", "user_config")
	ServiceName = Config.Viper.GetString("Server.Name")
	ServiceAddr = fmt.Sprintf("%s:%d", Config.Viper.GetString("Server.Address"), Config.Viper.GetInt("Server.Port"))
	EtcdAddress = fmt.Sprintf("%s:%d", Config.Viper.GetString("Etcd.Address"), Config.Viper.GetInt("Etcd.Port"))
)

func Init() {

}

func main() {

	svr := user.NewServer(new(UserServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
