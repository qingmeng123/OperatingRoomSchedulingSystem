package api

import (
	"OperatingRoomSchedulingSystem/service"
	"OperatingRoomSchedulingSystem/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		context.Header("Access-Control-Allow-Headers", "Content-Category,AccessToken,X-CSRF-Token, Authorization, Token")
		context.Header("Access-Control-Allow-Methods", "POST,PUT,GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Category")
		context.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept,Authorization")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}

// jwt验证
func jwtAuth(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		tool.RespSuccessfulWithData(ctx, "token为空")
		ctx.Abort()
		return
	}

	mc, err := service.ParseToken(token)
	if err != nil {
		fmt.Println("jwtAuthErr:", err.Error())
		tool.RespSuccessfulWithData(ctx, "token无效")
		ctx.Abort()
		return
	}
	ctx.Set("iUserGroupId", mc.User.GroupId)
	ctx.Set("iUsername", mc.User.Username)
	ctx.Next()
}

// JWT认证管理员
// jwt验证
func jwtAuthAdmin(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		tool.RespSuccessfulWithData(ctx, "token为空")
		ctx.Abort()
		return
	}

	mc, err := service.ParseToken(token)
	if err != nil {
		fmt.Println("jwtAuthErr:", err.Error())
		tool.RespSuccessfulWithData(ctx, "token无效")
		ctx.Abort()
		return
	}
	if mc.User.GroupId != 1 {
		tool.RespSuccessfulWithData(ctx, "用户权限不足")
		ctx.Abort()
		return
	}
	ctx.Set("iUsername", mc.User.Username)
	ctx.Next()
}
