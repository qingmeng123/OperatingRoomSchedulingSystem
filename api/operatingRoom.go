/*******
* @Author:qingmeng
* @Description:
* @File:OperatingRoom
* @Date:2024/4/4
 */

package api

import (
	"OperatingRoomSchedulingSystem/cache"
	"OperatingRoomSchedulingSystem/config"
	"OperatingRoomSchedulingSystem/dao"
	"OperatingRoomSchedulingSystem/tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
)

func addOperatingRoom(ctx *gin.Context) {
	OperatingRoomName := ctx.PostForm("name")
	if OperatingRoomName == "" {
		tool.RespParamError(ctx)
		return
	}

	categoryStr := ctx.PostForm("category")
	category, err := strconv.Atoi(categoryStr)
	if categoryStr == "" || err != nil {
		tool.RespParamError(ctx)
		return
	}
	pd := dao.OperatingRoomDao{}

	_, err = pd.GetOperatingRoomByName(OperatingRoomName)
	if err != gorm.ErrRecordNotFound {
		tool.RespErrorWithData(ctx, "该手术室已存在")
		return
	}

	err = pd.InsertOperatingRoom(OperatingRoomName, category)
	if err != nil {
		log.Println("insert OperatingRoom err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}

func modifyOperatingRoom(ctx *gin.Context) {
	id := ctx.PostForm("id")
	if id == "" {
		tool.RespParamError(ctx)
		return
	}
	OperatingRoomId, err := strconv.Atoi(id)
	if err != nil {
		tool.RespParamError(ctx)
		return
	}

	categoryStr := ctx.PostForm("category")
	category, err := strconv.Atoi(categoryStr)
	if categoryStr == "" || err != nil {
		tool.RespParamError(ctx)
		return
	}

	OperatingRoomName := ctx.PostForm("name")
	if OperatingRoomName == "" {
		tool.RespParamError(ctx)
		return
	}
	pd := dao.OperatingRoomDao{}
	_, err = pd.GetOperatingRoomByName(OperatingRoomName)
	if err != gorm.ErrRecordNotFound {
		tool.RespErrorWithData(ctx, "该手术室已存在")
		return
	}

	err = pd.UpdateOperatingRoom(OperatingRoomId, OperatingRoomName, category)
	if err != nil {
		log.Println("UpdateOperatingRoom err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}

func deleteOperatingRoom(ctx *gin.Context) {
	id := ctx.PostForm("id")
	if id == "" {
		tool.RespParamError(ctx)
		return
	}
	OperatingRoomId, err := strconv.Atoi(id)
	if err != nil {
		tool.RespParamError(ctx)
		return
	}

	pd := dao.OperatingRoomDao{}
	err = pd.DeleteOperatingRoom(OperatingRoomId)

	if err != nil {
		log.Println("deleteOperatingRoom err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}

func getOperatingRoom(ctx *gin.Context) {
	pd := dao.OperatingRoomDao{}
	OperatingRooms, err := pd.GetOperatingRooms()
	if err != nil {
		log.Println("getOperatingRooms err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessfulWithData(ctx, OperatingRooms)
	return
}

func queryOperatingRoomList(ctx *gin.Context) {
	pageNum := ctx.PostForm("pageNum")
	num, err := strconv.Atoi(pageNum)
	if err != nil || num <= 0 {
		num = 1
	}
	pageSize := config.PageSize
	name := ctx.PostForm("name")
	categoryStr := ctx.PostForm("category")
	category, err := strconv.Atoi(categoryStr)
	if categoryStr == "" || err != nil || category < 0 {
		category = 0
	}

	od := dao.OperatingRoomDao{}
	rooms, totalPages, err := od.QueryOperatingRoomList(num, pageSize, category, name)
	if err != nil {
		log.Println("get user list err:", err)
		tool.RespInternalError(ctx)
		return
	}

	//查看可用手术室
	startTimeStr := ctx.PostForm("start_time")
	endTimeStr := ctx.PostForm("end_time")
	if startTimeStr != "" && endTimeStr != "" {
		startTime, err := tool.ParseTime(startTimeStr)
		if err != nil {
			tool.RespParamError(ctx)
			return
		}
		endTime, err := tool.ParseTime(endTimeStr)
		if err != nil {
			tool.RespParamError(ctx)
			return
		}
		rooms, err = cache.GetAvailableOperatingRoomsByTime(rooms, startTime, endTime)
		if err != nil {
			log.Println("GetAvailableOperatingRoomsByTime err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	tool.RespSuccessfulWithData(ctx, map[string]interface{}{
		"pageList":  rooms,
		"totalPage": totalPages,
		"pageNum":   pageNum,
	})
}
