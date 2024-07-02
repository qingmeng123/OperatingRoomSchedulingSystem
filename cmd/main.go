package main

import (
	"OperatingRoomSchedulingSystem/api"
	"OperatingRoomSchedulingSystem/cache"
	"OperatingRoomSchedulingSystem/config"
	"OperatingRoomSchedulingSystem/dao"
)

func main() {
	config.Init()
	cache.InitCache()
	dao.InitDB()
	api.InitEngine()
}
