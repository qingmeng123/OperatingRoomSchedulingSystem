/*******
* @Author:qingmeng
* @Description:
* @File:resp
* @Date2021/12/10
 */

package tool

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespErrorWithData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": false,
		"data":   data,
	})
}

func RespParamError(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": false,
		"data":   "参数错误",
	})
}

func RespInternalError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"status": false,
		"data":   "服务器错误",
	})
}

func RespSuccessful(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   "成功",
	})
}

func RespSuccessfulWithData(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   data,
	})
}

func RespSensitiveError(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": false,
		"data":   "含有非法词汇",
	})
}
