/*******
* @Author:qingmeng
* @Description:
* @File:operatingRoom
* @Date:2024/4/6
 */

package model

type OperatingRoom struct {
	Id       int    `json:"id" form:"id"`
	Name     string `json:"name" form:"name" binding:"required"` //手术室名称
	Category int    `json:"category" form:"category" binding:"required"`
}

func (OperatingRoom) TableName() string {
	return "operating_room"
}
