/*******
* @Author:qingmeng
* @Description:
* @File:OperatingRoom
* @Date:2024/4/4
 */

package dao

import (
	"OperatingRoomSchedulingSystem/model"
)

type OperatingRoomDao struct {
}

func (dao *OperatingRoomDao) InsertOperatingRoom(name string, category int) error {
	OperatingRoom := model.OperatingRoom{Name: name, Category: category}
	result := GormDB.Create(&OperatingRoom)
	return result.Error
}

func (dao *OperatingRoomDao) GetOperatingRooms() (OperatingRooms []model.OperatingRoom, error error) {
	result := GormDB.Find(&OperatingRooms)
	return OperatingRooms, result.Error
}

func (dao *OperatingRoomDao) GetOperatingRoom(id int) (OperatingRoom model.OperatingRoom, error error) {
	result := GormDB.Where("id=?", id).First(&OperatingRoom)
	return OperatingRoom, result.Error
}

func (dao *OperatingRoomDao) GetOperatingRoomByName(name string) (OperatingRoom model.OperatingRoom, error error) {
	result := GormDB.Where("name=?", name).First(&OperatingRoom)
	return OperatingRoom, result.Error
}

func (dao *OperatingRoomDao) GetOperatingRoomsByCategory(category int) (OperatingRooms []model.OperatingRoom, error error) {
	result := GormDB.Where("category=?", category).Find(&OperatingRooms)
	return OperatingRooms, result.Error
}

func (dao *OperatingRoomDao) DeleteOperatingRoom(id int) error {
	OperatingRoom := model.OperatingRoom{Id: id}
	result := GormDB.Delete(&OperatingRoom)
	return result.Error
}

func (dao *OperatingRoomDao) UpdateOperatingRoom(id int, name string, category int) error {
	OperatingRoom := model.OperatingRoom{
		Id:       id,
		Name:     name,
		Category: category,
	}
	result := GormDB.Model(&OperatingRoom).Updates(OperatingRoom)
	return result.Error
}

func (dao *OperatingRoomDao) QueryOperatingRoomList(pageNum, pageSize, category int, name string) ([]model.OperatingRoom, int, error) {
	var rooms []model.OperatingRoom
	offset := (pageNum - 1) * pageSize
	// 创建一个 GORM 查询构建器
	query := GormDB.Offset(offset).Limit(pageSize)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if category != 0 {
		query = query.Where("category =?", category)
	}

	// 执行查询
	result := query.Find(&rooms)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	var totalRecords int64
	re := GormDB.Model(&model.OperatingRoom{}).Where(&model.OperatingRoom{Category: category}).Where("name LIKE ?", "%"+name+"%").Count(&totalRecords)
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)
	return rooms, int(totalPages), re.Error
}
