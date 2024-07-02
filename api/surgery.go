/*******
* @Author:qingmeng
* @Description:
* @File:Category
* @Date:2024/4/4
 */

package api

import (
	"OperatingRoomSchedulingSystem/cache"
	"OperatingRoomSchedulingSystem/config"
	"OperatingRoomSchedulingSystem/dao"
	"OperatingRoomSchedulingSystem/model"
	"OperatingRoomSchedulingSystem/tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

func timeRecommend(ctx *gin.Context) {
	durationStr := ctx.PostForm("duration")
	categoryIdStr := ctx.PostForm("category_id")
	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration <= 0 {
		tool.RespParamError(ctx)
		return
	}
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil || categoryId < 0 {
		tool.RespParamError(ctx)
		return
	}
	rd := dao.OperatingRoomDao{}
	rooms, err := rd.GetOperatingRoomsByCategory(categoryId)
	if err != nil {
		log.Println("GetOperatingRoomsByCategory err:", err)
		tool.RespInternalError(ctx)
		return
	}

	recommendOperatingTime, err := cache.RecommendOperatingTime(rooms, time.Minute*time.Duration(duration))
	if err != nil {
		log.Println("RecommendOperatingTime err:", err)
		tool.RespErrorWithData(ctx, "暂无推荐时间")
		return
	}
	tool.RespSuccessfulWithData(ctx, recommendOperatingTime)
	return
}

func querySurgeryList(ctx *gin.Context) {
	pageNum := ctx.PostForm("pageNum")
	num, err := strconv.Atoi(pageNum)
	if err != nil || num == 0 {
		num = 1
	}
	pageSize := config.PageSize

	//查询条件
	name := ctx.PostForm("name")
	categoryIdStr := ctx.PostForm("category_id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil || categoryId < 0 {
		categoryId = 0
	}
	roomIdStr := ctx.PostForm("room_id")
	roomId, err := strconv.Atoi(roomIdStr)
	if err != nil || roomId < 0 {
		roomId = 0
	}
	username := ctx.PostForm("username")

	od := dao.SurgeryDao{}
	surgery, totalPages, err := od.QuerySurgeryList(num, pageSize, categoryId, roomId, name, username)
	if err != nil {
		log.Println("get user list err:", err)
		tool.RespInternalError(ctx)
		return
	}
	//将username放在首位
	if username != "" {
		for i, s := range surgery {
			var users []model.User
			users = append(users, model.User{})
			for _, user := range s.Users {
				if strings.Contains(user.Name, username) {
					users[0] = user
				} else {
					users = append(users, user)
				}
			}
			surgery[i].Users = users
		}
	}

	//按天筛选
	dateStr := ctx.PostForm("date")
	if dateStr != "" {
		cstLocation, _ := time.LoadLocation("Asia/Shanghai")
		date, err := time.ParseInLocation("2006-01-02", dateStr, cstLocation)
		if err != nil {
			tool.RespParamError(ctx)
			return
		}
		var res []model.Surgery
		for i, _ := range surgery {
			surgery[i].StartTime = surgery[i].StartTime.Local()
			surgery[i].EndTime = surgery[i].EndTime.Local()
			if date.Day() == surgery[i].StartTime.Day() && date.Month() == surgery[i].StartTime.Month() && date.Year() == surgery[i].StartTime.Year() {
				res = append(res, surgery[i])
			}
		}
		surgery = res
	} else {
		for i, _ := range surgery {
			surgery[i].StartTime = surgery[i].StartTime.Local()
			surgery[i].EndTime = surgery[i].EndTime.Local()
		}
	}

	tool.RespSuccessfulWithData(ctx, map[string]interface{}{
		"pageList":  surgery,
		"totalPage": totalPages,
		"pageNum":   pageNum,
	})
}

// 完成操作
func finishSurgery(ctx *gin.Context) {
	codeStr := ctx.PostForm("state")
	state, err := strconv.Atoi(codeStr)
	if err != nil || state != 1 {
		tool.RespParamError(ctx)
		return
	}
	idStr := ctx.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		tool.RespParamError(ctx)
		return
	}
	//检查id是否存在
	sd := dao.SurgeryDao{}
	surgery, err := sd.GetSurgery(id)
	if err == gorm.ErrRecordNotFound {
		tool.RespParamError(ctx)
		return
	}
	if err != nil {
		log.Println("get surgery err:", err)
		tool.RespInternalError(ctx)
		return
	}

	if surgery.State != 0 {
		tool.RespParamError(ctx)
		return
	}

	//检查是否已经开始的手术
	if surgery.StartTime.After(time.Now()) {
		tool.RespErrorWithData(ctx, "手术还未开始")
		return
	}

	//更新
	//删除缓存
	err = cache.DelReserveRoom(surgery.RoomId, surgery.StartTime, surgery.EndTime)
	if err != nil {
		log.Println("DelReserveRoom err:", err)
		tool.RespInternalError(ctx)
		return
	}
	for _, user := range surgery.Users {
		err = cache.DelReserveUser(user.Id, surgery.StartTime, surgery.EndTime)
		if err != nil {
			log.Println("DelReserveUser err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	t := time.Now()
	err = sd.UpdateSurgeryDateAndState(surgery.Id, t, state)
	if err != nil {
		log.Println("UpdateSurgeryDateAndState err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessful(ctx)
	return
}

func delSurgery(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		tool.RespParamError(ctx)
		return
	}

	//检查id是否存在
	sd := dao.SurgeryDao{}
	surgery, err := sd.GetSurgery(id)
	if err == gorm.ErrRecordNotFound {
		tool.RespParamError(ctx)
		return
	}
	if err != nil {
		log.Println("get surgery err:", err)
		tool.RespInternalError(ctx)
		return
	}
	//是否已完成
	if surgery.State == 1 {
		tool.RespErrorWithData(ctx, "已完成手术，不能取消")
		return
	}

	//缓存删除
	err = cache.DelReserveRoom(surgery.RoomId, surgery.StartTime, surgery.EndTime)
	if err != nil {
		log.Println("DelReserveRoom err:", err)
		tool.RespInternalError(ctx)
		return
	}
	for _, user := range surgery.Users {
		err = cache.DelReserveUser(user.Id, surgery.StartTime, surgery.EndTime)
		if err != nil {
			log.Println("DelReserveUser err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}
	//数据库删除
	err = sd.DelSurgery(surgery)
	if err != nil {
		log.Println("del surgery from db err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessful(ctx)
	return
}

func addSurgery(ctx *gin.Context) {
	var req struct {
		Id         int          `json:"id" form:"id"`
		Name       string       `json:"name" form:"name"`
		CategoryId int          `json:"category_id" form:"category_id" `
		RoomId     int          `json:"room_id" form:"room_id" binding:"required"`
		StartTime  string       `json:"start_time" form:"start_time" binding:"required"`
		EndTime    string       `json:"end_time" form:"end_time" binding:"required"`
		Users      []model.User `json:"users"`
	}

	err := ctx.BindJSON(&req)
	if err != nil {
		log.Println("parse err:", err)
		tool.RespParamError(ctx)
		return
	}

	sd := dao.SurgeryDao{}
	surgery := model.Surgery{
		Id:         req.Id,
		Name:       req.Name,
		CategoryId: req.CategoryId,
		RoomId:     req.RoomId,
		StartTime:  time.Time{},
		EndTime:    time.Time{},
		State:      0,
		Users:      req.Users,
	}
	//解析时间
	surgery.StartTime, err = tool.ParseTime(req.StartTime)
	if err != nil {
		tool.RespParamError(ctx)
		return
	}
	surgery.EndTime, err = tool.ParseTime(req.EndTime)
	if err != nil {
		tool.RespParamError(ctx)
		return
	}

	//存入缓存
	for _, user := range req.Users {
		err = cache.ReserveUser(user.Id, surgery.StartTime, surgery.EndTime.Add(time.Minute*(-1)))
		if err != nil {
			log.Println("ReserveUser err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}
	err = cache.ReserveOperatingRoom(surgery.RoomId, surgery.StartTime, surgery.EndTime.Add(time.Minute*(-1)))
	if err != nil {
		log.Println("ReserveOperatingRoom err:", err)
		tool.RespInternalError(ctx)
		return
	}

	err = sd.InsertSurgery(surgery)
	if err != nil {
		log.Println("insert Category err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessful(ctx)
	return
}
