/*******
* @Author:qingmeng
* @Description:
* @File:Category
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

func addCategory(ctx *gin.Context) {
	CategoryName := ctx.PostForm("name")
	if CategoryName == "" {
		tool.RespParamError(ctx)
		return
	}
	pd := dao.CategoryDao{}

	_, err := pd.GetCategoryByName(CategoryName)
	if err != gorm.ErrRecordNotFound {
		tool.RespErrorWithData(ctx, "该职位已存在")
		return
	}

	err = pd.InsertCategory(CategoryName)
	if err != nil {
		log.Println("insert Category err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}

func modifyCategory(ctx *gin.Context) {
	id := ctx.PostForm("id")
	if id == "" {
		tool.RespParamError(ctx)
		return
	}
	CategoryId, err := strconv.Atoi(id)
	if err != nil {
		tool.RespParamError(ctx)
		return
	}

	number := ctx.PostForm("number")
	num, err := strconv.Atoi(number)
	if err != nil || num < 0 {
		num = 0
	}

	CategoryName := ctx.PostForm("name")
	if CategoryName == "" {
		tool.RespParamError(ctx)
		return
	}
	pd := dao.CategoryDao{}
	//查找是否存在该名称
	_, err = pd.GetCategoryByName(CategoryName)
	if err != gorm.ErrRecordNotFound {
		tool.RespErrorWithData(ctx, "该职位已存在")
		return
	}

	err = pd.UpdateCategory(CategoryId, CategoryName, num)
	if err != nil {
		log.Println("UpdateCategory err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}

func deleteCategory(ctx *gin.Context) {
	id := ctx.PostForm("id")
	if id == "" {
		tool.RespParamError(ctx)
		return
	}
	CategoryId, err := strconv.Atoi(id)
	if err != nil {
		tool.RespParamError(ctx)
		return
	}

	pd := dao.CategoryDao{}
	err = pd.DeleteCategory(CategoryId)

	if err != nil {
		log.Println("deleteCategory err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}

func getCategory(ctx *gin.Context) {
	pd := dao.CategoryDao{}
	categorys, err := pd.GetCategorys()
	if err != nil {
		log.Println("getCategorys err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessfulWithData(ctx, categorys)
	return
}
