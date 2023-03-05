package models

import (
	"byte_dance_5th/pkg/errortype"
	"errors"
	"log"
	"sync"
)

type LoginDao struct {
}

// User 登陆信息, UserInfo 获取用户follow等详细信息的外键
type User struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}

var (
	userLoginDao *LoginDao
	// LoginOnce 避免重复注册
	LoginOnce sync.Once
)

// UserLogin 登陆
func (u *LoginDao) UserLogin(username, password string, login *User) error {
	if login == nil {
		return errors.New("UserLogin:" + errortype.PointerIsNilErr)
	}
	DB.Where("username = ?", username).First(login)
	if login.Id == 0 {
		return errors.New(errortype.UserWrongOrNoExistErr)
	}
	DB.Where("username = ? AND password = ?", username, password).First(login)
	if login.Id == 0 {
		return errors.New(errortype.PasswordWrongErr)
	}
	return nil
}

func NewLoginDao() *LoginDao {
	LoginOnce.Do(func() {
		userLoginDao = new(LoginDao)
	})
	return userLoginDao
}

// UserAlreadyExist 查询账户是否已经存在
func (u *LoginDao) UserAlreadyExist(username string) bool {
	var user User
	log.Println(username)
	DB.Where("username = ?", username).First(&user)
	if user.Id == 0 {
		return false
	}
	return true
}
