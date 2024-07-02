/*******
* @Author:qingmeng
* @Description:
* @File:post
* @Date:2024/4/4
 */

package model

type Post struct {
	Id     int    `json:"id" form:"id"`
	Name   string `json:"name" form:"name" binding:"required"` //职位名
	Number int    `json:"number" form:"number"`
}

func (Post) TableName() string {
	return "post"
}
