/*******
* @Author:qingmeng
* @Description:
* @File:post
* @Date:2024/4/4
 */

package dao

import "OperatingRoomSchedulingSystem/model"

type PostDao struct {
}

func (dao *PostDao) InsertPost(name string) error {
	post := model.Post{Name: name}
	result := GormDB.Create(&post)
	return result.Error
}

func (dao *PostDao) GetPosts() (posts []model.Post, error error) {
	result := GormDB.Find(&posts)
	return posts, result.Error
}

func (dao *PostDao) GetPost(id int) (post model.Post, error error) {
	result := GormDB.Where("id=?", id).First(&post)
	return post, result.Error
}

func (dao *PostDao) GetPostByName(name string) (post model.Post, error error) {
	result := GormDB.Where("name=?", name).First(&post)
	return post, result.Error
}

func (dao *PostDao) DeletePost(id int) error {
	post := model.Post{Id: id}
	result := GormDB.Delete(&post)
	return result.Error
}

func (dao *PostDao) UpdatePost(id int, name string, number int) error {
	post := model.Post{
		Id:     id,
		Name:   name,
		Number: number,
	}
	result := GormDB.Model(&post).Updates(post)
	return result.Error
}
