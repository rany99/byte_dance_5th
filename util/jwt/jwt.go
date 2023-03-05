package jwt

import (
	"byte_dance_5th/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var JwtKey = []byte("imwave")

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

// GenerateToken 生成Token
func GenerateToken(user models.User) (string, error) {
	//token有效期设置为一周
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: user.UserInfoId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "ByteDance_5th",
			Subject:   "imwave",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(JwtKey)
	if err != nil {
		//log.Println("token生成失败", tokenStr)
		return "", err
	}
	return tokenStr, nil
}

// DecodeToken 解析token
func DecodeToken(tokenStr string) (*Claims, bool) {
	token, _ := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if token != nil {
		if key, ok := token.Claims.(*Claims); ok {
			if token.Valid {
				return key, true
			} else {
				return key, false
			}
		}
	}
	return nil, false
}
