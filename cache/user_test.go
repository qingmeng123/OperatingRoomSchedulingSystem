/*******
* @Author:qingmeng
* @Description:
* @File:room_test
* @Date:2024/5/9
 */

package cache

import (
	"OperatingRoomSchedulingSystem/model"
	"log"
	"testing"
	"time"
)

func TestUserFunctions(t *testing.T) {
	InitCache()
	// 添加手术室信息
	rooms := []model.OperatingRoom{
		{Id: 1, Name: "User 1"},
		{Id: 2, Name: "User 2"},
		{Id: 3, Name: "User 3"},
	}
	//预约手术
	startTime := time.Date(2024, 5, 9, 8, 0, 0, 0, time.UTC) // 模拟开始时间
	endTime := time.Date(2024, 5, 9, 10, 0, 0, 0, time.UTC)  // 模拟结束时间
	err := ReserveUser(1, startTime, endTime)
	if err != nil {
		log.Println("reserve err:", err)
		return
	}
	err = ReserveOperatingRoom(1, startTime.Add(time.Hour*(3)), endTime.Add(time.Hour*(3)))
	if err != nil {
		log.Println("reserve err:", err)
		return
	}

	//查看预约信息
	res, err := GetReserveRoomByTimeRange(1, startTime.Add(time.Hour*(-1)), endTime.Add(time.Hour*3))
	if err != nil {
		log.Println("get err:", err)
		return
	}
	log.Println("res:", res)

	//删除预约
	err = DelReserveRoom(1, startTime.Add(time.Hour*(3)), endTime.Add(time.Hour*(3)))
	//查看预约信息
	res, err = GetReserveRoomByTimeRange(1, startTime.Add(time.Hour*(-1)), endTime.Add(time.Hour*3))
	if err != nil {
		log.Println("get err:", err)
		return
	}
	log.Println("after del res:", res)
	//获取可用房间
	res1, err := GetAvailableOperatingRoomsByTime(rooms, startTime.Add(time.Hour*(3)), endTime.Add(time.Hour*5))
	if err != nil {
		log.Println("get1 err:", err)
		return
	}
	log.Println("res1:", res1)
}
