/*******
* @Author:qingmeng
* @Description:
* @File:post
* @Date:2024/4/4
 */

package api

import (
	"OperatingRoomSchedulingSystem/dao"
	"OperatingRoomSchedulingSystem/tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
)

func addPost(ctx *gin.Context) {
	postName := ctx.PostForm("name")
	if postName == "" {
		tool.RespParamError(ctx)
		return
	}
	pd := dao.PostDao{}

	_, err := pd.GetPostByName(postName)
	if err != gorm.ErrRecordNotFound {
		tool.RespErrorWithData(ctx, "该职位已存在")
		return
	}

	err = pd.InsertPost(postName)
	if err != nil {
		log.Println("insert post err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}

func modifyPost(ctx *gin.Context) {
	id := ctx.PostForm("id")
	if id == "" {
		tool.RespParamError(ctx)
		return
	}
	postId, err := strconv.Atoi(id)
	if err != nil {
		tool.RespParamError(ctx)
		return
	}

	number := ctx.PostForm("number")
	num, err := strconv.Atoi(number)
	if err != nil || num < 0 {
		num = 0
	}

	postName := ctx.PostForm("name")
	if postName == "" {
		tool.RespParamError(ctx)
		return
	}
	pd := dao.PostDao{}
	//查找是否存在该名称
	_, err = pd.GetPostByName(postName)
	if err != gorm.ErrRecordNotFound {
		tool.RespErrorWithData(ctx, "该职位已存在")
		return
	}

	err = pd.UpdatePost(postId, postName, num)
	if err != nil {
		log.Println("UpdatePost err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}

func deletePost(ctx *gin.Context) {
	id := ctx.PostForm("id")
	if id == "" {
		tool.RespParamError(ctx)
		return
	}
	postId, err := strconv.Atoi(id)
	if err != nil {
		tool.RespParamError(ctx)
		return
	}

	pd := dao.PostDao{}
	err = pd.DeletePost(postId)

	if err != nil {
		log.Println("deletePost err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}

func getPost(ctx *gin.Context) {
	pd := dao.PostDao{}
	posts, err := pd.GetPosts()
	if err != nil {
		log.Println("getPosts err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessfulWithData(ctx, posts)
	return
}
