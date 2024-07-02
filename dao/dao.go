package dao

import (
	"OperatingRoomSchedulingSystem/config"
	"OperatingRoomSchedulingSystem/model"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *sql.DB
var GormDB *gorm.DB

func InitDB() {
	dsn := config.DbUser + ":" + config.DbPassWord + "@tcp(" + config.DbHost + ":" + config.DbPort + ")/" + config.DbName + "?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	DB = db

	//gormDB连接mysql
	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		fmt.Println("mysql conn err:", err)
		return
	}
	GormDB = gormDB
	gormDB.AutoMigrate(&model.Surgery{}, &model.User{})

}
