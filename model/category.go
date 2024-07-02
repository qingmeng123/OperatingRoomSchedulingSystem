/*******
* @Author:qingmeng
* @Description:
* @File:category
* @Date:2024/5/6
 */

package model

type Category struct {
	Id   int    `json:"id" form:"id"`
	Name string `json:"name" form:"name" binding:"required"` //类别名
}

func (Category) TableName() string {
	return "category"
}
