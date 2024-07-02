package model

type User struct {
	Id       int    `json:"id" gorm:"column:id"`
	Username string `json:"username" form:"username" binding:"required" gorm:"size:20"`
	Password string `json:"password" form:"password" binding:"required"`
	Gender   bool   `json:"gender" form:"gender" `                      //false男，true女
	Name     string `json:"name" form:"name" binding:"required"`        //姓名
	Phone    string `json:"phone" form:"phone" binding:"required"`      //账号电话
	WorkYear int    `json:"work_year" form:"work_year" `                //工作年龄
	GroupId  int    `json:"group_id"`                                   //成员组id,1为管理员，0为医师护士等
	PostId   int    `json:"post_id" form:"post_id" binding:"required" ` //职位id
}

func (User) TableName() string {
	return "user"
}
