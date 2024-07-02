package model

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyClaims struct {
	User User
	Type string //"REFRESH_TOKEN"表示为一个refresh token，"TOKEN"表示为一个token
	Time time.Time
	jwt.StandardClaims
}

// 第三方登陆的token
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
