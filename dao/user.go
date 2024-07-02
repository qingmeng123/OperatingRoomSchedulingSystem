package dao

import (
	"OperatingRoomSchedulingSystem/model"
	"fmt"
	"gorm.io/gorm"
)

type UserDao struct {
}

func (dao *UserDao) InsertUser(user model.User) error {
	result := GormDB.Create(&user)
	return result.Error
}

// 事务处理删除
func (dao *UserDao) DeleteUser(id int) error {
	err := GormDB.Transaction(func(tx *gorm.DB) error {
		user := model.User{}
		result := tx.First(&user, id)
		if result.Error != nil {
			return result.Error
		}
		fmt.Println(user)
		result = tx.Delete(&user, id)

		if result.Error != nil {
			return result.Error
		}

		result = tx.Model(&model.Post{Id: user.PostId}).Update("number", gorm.Expr("number-?", 1))
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	return err
}

func (dao *UserDao) QueryUserList(pageNum, pageSize int, phone, name string) ([]model.User, int, error) {
	var users []model.User
	offset := (pageNum - 1) * pageSize
	// 创建一个 GORM 查询构建器
	query := GormDB.Offset(offset).Limit(pageSize)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}

	// 执行查询
	result := query.Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	var totalRecords int64
	re := GormDB.Model(&model.User{}).Where("name LIKE ?", "%"+name+"%").Where("phone LIKE ?", "%"+phone+"%").Count(&totalRecords)
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)
	return users, int(totalPages), re.Error
}

func (dao *UserDao) GetUserList(pageNum, pageSize int) ([]model.User, int, error) {
	var users []model.User
	offset := (pageNum - 1) * pageSize
	result := GormDB.Offset(offset).Limit(pageSize).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	var totalRecords int64
	re := GormDB.Model(&model.User{}).Count(&totalRecords)
	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)
	return users, int(totalPages), re.Error
}

func (dao *UserDao) GetUsers() ([]model.User, error) {
	var users []model.User
	result := GormDB.Find(&users)
	return users, result.Error
}

func (dao *UserDao) SelectUserByUsername(username string) (model.User, error) {
	user := model.User{}
	result := GormDB.Where("username=?", username).First(&user)
	return user, result.Error
}

func (dao *UserDao) GetUser(id int) (model.User, error) {
	user := model.User{}
	result := GormDB.Where("id=?", id).First(&user)
	return user, result.Error
}

func (dao *UserDao) SelectBasicUserByUsername(username string) (model.User, error) {
	user := model.User{}
	result := GormDB.Select("id", "password", "group_id", "post_id").Where("username=?", username).First(&user)
	user.Username = username
	return user, result.Error
}

func (dao *UserDao) UpdatePassword(username, newPassword string) error {
	result := GormDB.Model(&model.User{}).Where("username=?", username).Update("password", newPassword)
	return result.Error
}

// UpdatePhone 更新用户电话
func (dao *UserDao) UpdatePhone(username string, phone string) error {
	result := GormDB.Model(&model.User{}).Where("username=?", username).Update("phone", phone)
	return result.Error
}

// UpdateName 更新昵称
func (dao *UserDao) UpdateName(username string, name string) error {
	_, err := DB.Exec("update scheduling_system.user set name=? where username=?", name, username)
	return err
}

// UpdateGender 更新性别
func (dao *UserDao) UpdateGender(username string, gender bool) error {

	_, err := DB.Exec("update scheduling_system.user set gender=? where username=?", gender, username)
	return err
}

func (dao *UserDao) UpdatePost(username string, id int) error {
	user, err := dao.SelectUserByUsername(username)
	if err != nil {
		return err
	}
	err = GormDB.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&model.User{}).Where("username=?", username).Update("post_id", id)
		if result.Error != nil {
			return result.Error
		}
		result = tx.Model(&model.Post{Id: user.PostId}).Update("number", gorm.Expr("number-?", 1))
		if result.Error != nil {
			return result.Error
		}

		result = tx.Model(&model.Post{Id: id}).Update("number", gorm.Expr("number+?", 1))
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	return err
}

func (dao *UserDao) UpdateGroupId(username string, i int) error {
	_, err := DB.Exec("update scheduling_system.user set group_id=? where username=?", i, username)
	return err
}

// UpdateYear 更新工龄
func (dao *UserDao) UpdateYear(username string, year int) error {
	_, err := DB.Exec("update scheduling_system.user set work_year=? where username=?", year, username)
	return err
}
