/*******
* @Author:qingmeng
* @Description:
* @File:router
* @Date2021/12/10
 */

package api

import (
	"OperatingRoomSchedulingSystem/config"
	"github.com/gin-gonic/gin"
)

func InitEngine() {
	engine := gin.Default()
	engine.Use(Cors())

	apiGroup := engine.Group("/api")
	{
		apiGroup.POST("/user/register", register)              //注册
		apiGroup.POST("/user/login", login)                    //登陆
		apiGroup.GET("/post", getPost)                         //获取职位
		apiGroup.POST("/admin/surgery/list", querySurgeryList) //查询手术
		apiGroup.GET("/admin/operatingRoom", getOperatingRoom)

		userGroup := apiGroup.Group("/user", jwtAuth)
		{
			userGroup.POST("/info", changeInformation) //修改信息
			userGroup.GET("/info", viewUserInfo)       //查看信息
			userGroup.GET("", getUser)
			//密码组
			passwordGroup := userGroup.Group("/password")
			{
				passwordGroup.POST("", changePassword) //登陆后的直接修改密码
			}

		}

		//管理员操作
		adminGroup := apiGroup.Group("/admin", jwtAuthAdmin)
		{
			adminGroup.POST("/modify", adModifyInfo)          //管理员操作用户
			adminGroup.GET("/userList", getUserList)          //获取用户列表
			adminGroup.POST("/query/userList", queryUserList) //查询用户列表
			adminGroup.POST("/delUser", delUsers)             //删除用户
			adminGroup.PUT("/post", addPost)                  //管理员添加职位
			adminGroup.POST("/post", modifyPost)              //管理员修改职位
			adminGroup.DELETE("/post", deletePost)            //删除职位
			//管理员对手术室资源操作
			adminGroup.PUT("/operatingRoom", addOperatingRoom)
			adminGroup.POST("/operatingRoom", modifyOperatingRoom)
			adminGroup.POST("/delOperatingRoom", deleteOperatingRoom)

			adminGroup.POST("/query/operatingRoomList", queryOperatingRoomList) //查询手术室列表
			adminGroup.PUT("/category", addCategory)                            //管理员添加手术类别
			adminGroup.POST("/category", modifyCategory)                        //管理员修改手术类别
			adminGroup.POST("/delCategory", deleteCategory)                     //删除手术类别
			adminGroup.GET("/category", getCategory)                            //获取手术类别

			//执行手术
			surgeryGroup := adminGroup.Group("/surgery")
			{
				surgeryGroup.POST("/reserve", addSurgery)
				surgeryGroup.POST("/users", getAvailableUsers)
				surgeryGroup.POST("/del", delSurgery)
				surgeryGroup.POST("/finish", finishSurgery)
				surgeryGroup.POST("/recommend", timeRecommend)
			}

		}
	}

	engine.Run(config.HttpPort)
}
