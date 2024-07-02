package api

import (
	"OperatingRoomSchedulingSystem/cache"
	"OperatingRoomSchedulingSystem/config"
	"OperatingRoomSchedulingSystem/dao"
	"OperatingRoomSchedulingSystem/model"
	"OperatingRoomSchedulingSystem/service"
	"OperatingRoomSchedulingSystem/tool"
	_ "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// 注册
func register(ctx *gin.Context) {

	var user model.User
	if err := ctx.Bind(&user); err != nil {
		log.Println("parse user err:", err)
		tool.RespParamError(ctx)
		return
	}
	//检验用户名是否含有敏感词
	flag := tool.CheckIfSensitive(user.Username)
	flag2 := tool.CheckIfSensitive(user.Name)
	if flag || flag2 {
		tool.RespSensitiveError(ctx)
		return
	}
	u := service.UserService{}

	//用户名是否已存在
	flag, err := u.IsExistUsername(user.Username)
	if err != nil {
		fmt.Println("judge exist username err:", err)
		tool.RespInternalError(ctx)
		return
	}

	if flag {
		tool.RespErrorWithData(ctx, "用户名已存在")
		return
	}
	flag = u.IsPasswordReasonable(user.Password)
	if !flag {
		tool.RespErrorWithData(ctx, "密码不合理")
		return
	}

	//密码加盐
	user.Password, err = tool.AddSalt(user.Password)
	if err != nil {
		log.Println("user_register AddSalt err:", err)
		tool.RespParamError(ctx)
		return
	}
	//注册
	err = u.Register(user)
	if err != nil {
		fmt.Println("register err: ", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulWithData(ctx, map[string]interface{}{
		"data":     "注册成功",
		"username": user.Username,
		"name":     user.Name,
		"groupId":  0,
	})

}

// 登陆
func login(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	u := service.UserService{}

	flag, err := u.IsPasswordCorrect(username, password)
	if err != nil {
		fmt.Println("judge password correct err:", err)
		tool.RespInternalError(ctx)
		return
	}
	if !flag {
		tool.RespErrorWithData(ctx, "密码错误")
		return
	}

	//获取用户固定信息
	basicUserinfo, err := u.GetBasicUserinfo(username)
	basicUserinfo.Username = username
	if err != nil {
		fmt.Println("getBasicUserInfo err:", err)
		tool.RespInternalError(ctx)
		return
	}
	//创建token,有效期5天
	tokenString, err := service.CreateToken(basicUserinfo, 5*24*60*60000000, "TOKEN")
	if err != nil {
		fmt.Println("create token err:", err)
		tool.RespInternalError(ctx)
		return
	}
	//创建refreshToken，有效期5天
	refreshTokenString, err := service.CreateToken(basicUserinfo, 5*24*60*60000000, "TOKEN")
	if err != nil {
		fmt.Println("create token err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulWithData(ctx, map[string]interface{}{
		"data":         "登陆成功",
		"id":           basicUserinfo.Id,
		"groupId":      basicUserinfo.GroupId,
		"token":        tokenString,
		"refreshToken": refreshTokenString,
	})

}

// 管理员更改用户信息
func adModifyInfo(ctx *gin.Context) {
	//指定用户
	idStr := ctx.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		tool.RespErrorWithData(ctx, "用户id无效")
		return
	}

	u := service.UserService{}
	username := ctx.PostForm("username")

	//是否存在该用户
	flag, err := u.IsExistUsername(username)
	if username == "" || !flag {
		tool.RespErrorWithData(ctx, "username无效")
		return
	}
	if err != nil {
		fmt.Println("check user err:", err)
		tool.RespInternalError(ctx)
		return
	}

	name := ctx.PostForm("name")
	if name != "" {
		err := u.UpdateName(username, name)
		if err != nil {
			fmt.Println("update name err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	//修改职位
	iPostId := ctx.PostForm("post_id")
	if iPostId != "" {
		postId, err := strconv.Atoi(iPostId)
		if err != nil {
			fmt.Println("postId to int err:", err)
			tool.RespInternalError(ctx)
			return
		}
		err = u.UpdatePostId(username, postId)
		if err != nil {
			fmt.Println("update postId err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	//更新电话
	phone := ctx.PostForm("phone")

	if phone != "" {
		if len(phone) != 11 {
			tool.RespErrorWithData(ctx, "phone无效")
			return
		}
		err := u.UpdatePhone(username, phone)
		if err != nil {
			fmt.Println("UpdatePhone err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	workYear := ctx.PostForm("work_year")
	if workYear != "" {
		wYear, err := strconv.Atoi(workYear)
		if err != nil {
			fmt.Println("workYear string to float err:", err)
			tool.RespErrorWithData(ctx, "工作年龄错误")
			return
		}
		if wYear <= 0 {
			tool.RespErrorWithData(ctx, "工作年龄错误")
			return
		}

		err = u.UpdateYear(username, wYear)
		if err != nil {
			fmt.Println("adModifyInfo err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	//修改密码
	psw := ctx.PostForm("password")
	if psw != "" {
		//密码加盐
		password, err := tool.AddSalt(psw)
		if err != nil {
			log.Println(" AddSalt err:", err)
			tool.RespParamError(ctx)
			return
		}
		err = u.ChangePassword(username, password)
		if err != nil {
			log.Println(" ChangePassword err:", err)
			tool.RespParamError(ctx)
			return
		}
	}

	//修改权限
	ipower := ctx.PostForm("group_id")
	if ipower != "" {
		power, err := strconv.Atoi(ipower)
		if err != nil {
			fmt.Println("group_id string to float err:", err)
			tool.RespErrorWithData(ctx, "操作权限错误")
			return
		}
		err = u.UpdateGroupId(username, power)
		if err != nil {
			log.Println("UpdateGroupId err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	//更新性别
	gender := ctx.PostForm("gender")
	if gender != "" {
		gen := false
		if gender == "true" {
			gen = true
		}
		err = u.UpdateGender(username, gen)
		if err != nil {
			fmt.Println("update gender err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	tool.RespSuccessfulWithData(ctx, "修改成功")
}

// 更新信息
func changeInformation(ctx *gin.Context) {
	iUsername, _ := ctx.Get("iUsername")
	username := iUsername.(string)
	u := service.UserService{}

	//更新电话
	phone := ctx.PostForm("phone")

	if phone != "" {
		if len(phone) != 11 {
			tool.RespErrorWithData(ctx, "phone无效")
			return
		}
		err := u.UpdatePhone(username, phone)
		if err != nil {
			fmt.Println("UpdatePhone err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	//更新昵称
	name := ctx.PostForm("name")
	if name != "" {
		flag := tool.CheckIfSensitive(name)
		if flag {
			tool.RespErrorWithData(ctx, "name格式不符合要求")
			return
		}
		err := u.UpdateName(username, name)
		if err != nil {
			fmt.Println("update name err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	//更新性别
	gender := ctx.PostForm("gender")
	if gender != "" {
		gen := false
		if gender == "true" {
			gen = true
		}
		err := u.UpdateGender(username, gen)
		if err != nil {
			fmt.Println("update gender err:", err)
			tool.RespInternalError(ctx)
			return
		}
	}

	tool.RespSuccessfulWithData(ctx, "修改成功")
}

// 查看个人信息
func viewUserInfo(ctx *gin.Context) {
	iUsername, _ := ctx.Get("iUsername")
	username := iUsername.(string)
	u := service.UserService{}

	userinfo, err := u.GetUserinfoByUserName(username)
	if err != nil {
		fmt.Println("getUserinfo err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessfulWithData(ctx, userinfo)

}

func getUser(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		tool.RespParamError(ctx)
		return
	}
	u := service.UserService{}

	userinfo, err := u.GetUser(id)
	if err != nil {
		fmt.Println("getUserinfo err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulWithData(ctx, userinfo)

}

// 登陆后修改密码
func changePassword(ctx *gin.Context) {
	oldPassword := ctx.PostForm("old_password")
	newPassword := ctx.PostForm("new_password")
	iUsername, _ := ctx.Get("iUsername")
	username := iUsername.(string)
	u := service.UserService{}

	//检验旧密码是否正确
	flag, err := u.IsPasswordCorrect(username, oldPassword)
	if err != nil {
		fmt.Println("judge password correct err:", err)
		tool.RespInternalError(ctx)
		return
	}
	if !flag {
		tool.RespErrorWithData(ctx, "旧密码错误")
		return
	}

	//验证新密码合理性
	flag = u.IsPasswordReasonable(newPassword)
	if !flag {
		tool.RespErrorWithData(ctx, "新密码无效")
		return
	}

	//修改为新密码
	password, err := tool.AddSalt(newPassword)
	if err != nil {
		fmt.Println("change password err:", err)
		tool.RespInternalError(ctx)
		return
	}
	err = u.ChangePassword(username, password)
	if err != nil {
		fmt.Println("change password err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulWithData(ctx, "修改成功,请重新登陆！")
}

func delUsers(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		tool.RespErrorWithData(ctx, "选择用户错误")
		return
	}

	u := dao.UserDao{}
	err = u.DeleteUser(id)
	if err != nil {
		log.Println("delete user err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessful(ctx)
}

func queryUserList(ctx *gin.Context) {
	pageNum := ctx.PostForm("pageNum")
	num, err := strconv.Atoi(pageNum)
	if err != nil || num <= 0 {
		num = 1
	}
	pageSize := config.PageSize
	name := ctx.PostForm("name")
	phoneStr := ctx.PostForm("phone")
	ud := dao.UserDao{}
	userList, totalPages, err := ud.QueryUserList(num, pageSize, phoneStr, name)
	if err != nil {
		log.Println("get user list err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulWithData(ctx, map[string]interface{}{
		"pageList":  userList,
		"totalPage": totalPages,
		"pageNum":   pageNum,
	})
}

func getUserList(ctx *gin.Context) {
	pageNum := ctx.Query("pageNum")
	num, err := strconv.Atoi(pageNum)
	if err != nil || num <= 0 {
		num = 1
	}
	pageSize := config.PageSize
	ud := dao.UserDao{}
	userList, totalPages, err := ud.GetUserList(num, pageSize)
	if err != nil {
		log.Println("get user list err:", err)
		tool.RespInternalError(ctx)
		return
	}
	tool.RespSuccessfulWithData(ctx, map[string]interface{}{
		"pageList":  userList,
		"totalPage": totalPages,
		"pageNum":   pageNum,
	})
}

// 获取空闲医生
func getAvailableUsers(ctx *gin.Context) {
	startTimeStr := ctx.PostForm("start_time")
	endTimeStr := ctx.PostForm("end_time")
	if startTimeStr == "" || endTimeStr == "" {
		tool.RespParamError(ctx)
		return
	}

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

	//获取所有医生
	ud := dao.UserDao{}
	users, err := ud.GetUsers()
	if err != nil {
		log.Println("get users err:", err)
		tool.RespInternalError(ctx)
		return
	}

	users, err = cache.GetAvailableUsersByTime(users, startTime, endTime)
	if err != nil {
		log.Println("GetAvailableOperatingRoomsByTime err:", err)
		tool.RespInternalError(ctx)
		return
	}

	tool.RespSuccessfulWithData(ctx, users)

}
