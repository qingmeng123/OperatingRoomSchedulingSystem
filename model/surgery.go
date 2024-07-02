/*******
* @Author:qingmeng
* @Description:
* @File:surgery
* @Date:2024/5/8
 */

package model

import (
	"time"
)

type Surgery struct {
	Id         int       `json:"id" form:"id"`
	Name       string    `json:"name" form:"name"`
	CategoryId int       `json:"category_id" form:"category_id" `
	RoomId     int       `json:"room_id" form:"room_id" binding:"required"`
	StartTime  time.Time `json:"start_time" form:"start_time" binding:"required"`
	EndTime    time.Time `json:"end_time" form:"end_time" binding:"required"`
	State      int       `json:"state" form:"state"`
	Users      []User    `json:"users" form:"users" gorm:"many2many:surgery_user;"`
}

func (Surgery) TableName() string {
	return "surgery"
}
