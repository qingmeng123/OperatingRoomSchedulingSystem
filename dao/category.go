/*******
* @Author:qingmeng
* @Description:
* @File:Category
* @Date:2024/4/4
 */

package dao

import "OperatingRoomSchedulingSystem/model"

type CategoryDao struct {
}

func (dao *CategoryDao) InsertCategory(name string) error {
	Category := model.Category{Name: name}
	result := GormDB.Create(&Category)
	return result.Error
}

func (dao *CategoryDao) GetCategorys() (Categorys []model.Category, error error) {
	result := GormDB.Find(&Categorys)
	return Categorys, result.Error
}

func (dao *CategoryDao) GetCategory(id int) (Category model.Category, error error) {
	result := GormDB.Where("id=?", id).First(&Category)
	return Category, result.Error
}

func (dao *CategoryDao) GetCategoryByName(name string) (Category model.Category, error error) {
	result := GormDB.Where("name=?", name).First(&Category)
	return Category, result.Error
}

func (dao *CategoryDao) DeleteCategory(id int) error {
	Category := model.Category{Id: id}
	result := GormDB.Delete(&Category)
	return result.Error
}

func (dao *CategoryDao) UpdateCategory(id int, name string, number int) error {
	Category := model.Category{
		Id:   id,
		Name: name,
	}
	result := GormDB.Model(&Category).Updates(Category)
	return result.Error
}
