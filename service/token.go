package service

import (
	"OperatingRoomSchedulingSystem/config"
	"OperatingRoomSchedulingSystem/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var jwtKey = []byte(config.JwtKey)

// CreateToken 创建token
func CreateToken(user model.User, ExpireTime int64, tokenType string) (string, error) {
	cla := model.MyClaims{
		User: user,
		Type: tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + ExpireTime, //过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)

	return token.SignedString(jwtKey) // 进行签名生成对应的token
}

// ParseToken 解析token
func ParseToken(tokenString string) (*model.MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 获取 token
func GetToken(url string) (*model.Token, error) {
	// 形成请求
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodPost, url, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")

	// 发送请求并获得响应
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		return nil, err
	}

	// 将响应体解析为 token，并返回
	var token model.Token
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		return nil, err
	}
	return &token, nil
}

// 获取第三方用户信息
func GetUserInfo(token *model.Token) (map[string]interface{}, error) {

	// 形成请求
	var userInfoUrl = "https://api.github.com/user" // github用户信息获取接口
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))

	// 发送请求并获取响应
	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return nil, err
	}

	// 将响应的数据写入 userInfo 中，并返回
	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}
