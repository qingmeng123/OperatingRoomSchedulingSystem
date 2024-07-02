/*******
* @Author:qingmeng
* @Description:
* @File:surgery
* @Date:2024/5/8
 */

package dao

import (
	"OperatingRoomSchedulingSystem/model"
	"gorm.io/gorm/clause"
	"time"
)

type SurgeryDao struct {
}

func (dao *SurgeryDao) InsertSurgery(surgery model.Surgery) error {
	result := GormDB.Create(&surgery)
	return result.Error
}

func (dao *SurgeryDao) DelSurgery(surgery model.Surgery) error {
	result := GormDB.Select(clause.Associations).Delete(&surgery)
	return result.Error
}

func (dao *SurgeryDao) GetSurgery(id int) (model.Surgery, error) {
	var surgery model.Surgery
	surgery.Id = id
	result := GormDB.Preload("Users").First(&surgery)
	return surgery, result.Error
}

func (dao *SurgeryDao) QuerySurgeryList(pageNum, pageSize, categoryId, roomId int, name, username string) ([]model.Surgery, int, error) {
	var surgeries []model.Surgery
	var err error
	offset := (pageNum - 1) * pageSize
	// 创建一个 GORM 查询构建器
	query := GormDB
	if pageNum >= 0 {
		query = GormDB.Offset(offset).Limit(pageSize).Order("start_time desc")
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if categoryId != 0 {
		query = query.Where("category_id=?", categoryId)
	}
	if roomId != 0 {
		query = query.Where("room_id=?", roomId)
	}
	// 执行查询
	if username != "" {
		err = query.Model(&model.Surgery{}).Joins("JOIN surgery_user ON surgery.id = surgery_user.surgery_id").
			Joins("JOIN user ON user.id = surgery_user.user_id").
			Where("user.name LIKE ?", "%"+username+"%").
			Preload("Users").Find(&surgeries).Error
		if err != nil {
			return nil, 0, err
		}
	}
	if username == "" {
		err = query.Model(&model.Surgery{}).Preload("Users").Find(&surgeries).Error
		if err != nil {
			return nil, 0, err
		}
	}
	var totalRecords int64
	re := GormDB.Model(&model.Surgery{}).Count(&totalRecords)
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)
	return surgeries, int(totalPages), re.Error
}

func (dao *SurgeryDao) UpdateSurgeryDateAndState(id int, t time.Time, state int) error {
	result := GormDB.Model(&model.Surgery{}).Where("id=?", id).Updates(model.Surgery{EndTime: t, State: state})
	return result.Error
}
