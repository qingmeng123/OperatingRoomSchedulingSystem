/*******
* @Author:qingmeng
* @Description:
* @File:room
* @Date:2024/5/9
 */

package cache

import (
	"OperatingRoomSchedulingSystem/config"
	"OperatingRoomSchedulingSystem/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

// 时间段结构体
type TimeSlot struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// 分数计算
func count(t time.Time) int64 {
	return int64(t.Hour()*60 + t.Minute())
}

// 预约手术房间
func ReserveOperatingRoom(roomID int, startTime time.Time, endTime time.Time) error {
	key := fmt.Sprintf("operating_room_schedule:%d", roomID)
	bookingInfo := map[string]interface{}{
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
	}
	bookingJSON, err := json.Marshal(bookingInfo)
	if err != nil {
		return err
	}
	err = RedisClient.ZAdd(context.Background(), key, &redis.Z{
		Score:  float64(count(endTime)), // 使用预订开始时间作为分数
		Member: string(bookingJSON),     // 存储预订信息的 JSON 字符串
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

// 删除手术房间
func DelReserveRoom(roomID int, startTime time.Time, endTime time.Time) error {
	key := fmt.Sprintf("operating_room_schedule:%d", roomID)
	bookingInfo := map[string]interface{}{
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
	}
	bookingJSON, err := json.Marshal(bookingInfo)
	if err != nil {
		return err
	}
	err = RedisClient.ZRem(context.Background(), key, string(bookingJSON)).Err()
	if err != nil {
		return err
	}
	return nil
}

// 查询指定范围预定信息
func GetReserveRoomByTimeRange(roomID int, startTime time.Time, endTime time.Time) ([]TimeSlot, error) {
	key := fmt.Sprintf("operating_room_schedule:%d", roomID)
	// 从 Redis 中获取手术室的预订信息
	bookings, err := RedisClient.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min: strconv.FormatInt(count(startTime), 10), // 最小分数，即开始时间
		Max: strconv.FormatInt(count(endTime), 10),   // 最大分数，即结束时间
	}).Result()
	if err != nil {
		return nil, err
	}

	//使用时间
	usedTime := []TimeSlot{}
	//解析
	for _, bookingInfoJSON := range bookings {
		var slot TimeSlot
		err := json.Unmarshal([]byte(bookingInfoJSON), &slot)
		if err != nil {
			return usedTime, err
		}
		usedTime = append(usedTime, slot)
	}

	return usedTime, nil
}

// 时间上获取可用房间
func GetAvailableOperatingRoomsByTime(rooms []model.OperatingRoom, startTime time.Time, endTime time.Time) ([]model.OperatingRoom, error) {
	var availableRooms []model.OperatingRoom
	for _, room := range rooms {
		isAvailable, err := IsOperatingRoomAvailable(room.Id, startTime, endTime)
		if err != nil {
			return nil, err
		}
		if isAvailable {
			availableRooms = append(availableRooms, room)
		}
	}
	return availableRooms, nil
}

// 是否是时间允许的房间
func IsOperatingRoomAvailable(roomID int, startTime time.Time, endTime time.Time) (bool, error) {
	key := fmt.Sprintf("operating_room_schedule:%d", roomID)
	bookings, err := RedisClient.ZRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return false, err
	}
	for _, bookingInfoJSON := range bookings {
		var bookingInfo map[string]interface{}
		err := json.Unmarshal([]byte(bookingInfoJSON), &bookingInfo)
		if err != nil {
			return false, err
		}
		bookingStartTime, _ := time.Parse(time.RFC3339, bookingInfo["start_time"].(string))
		bookingEndTime, _ := time.Parse(time.RFC3339, bookingInfo["end_time"].(string))
		if !(endTime.Before(bookingStartTime) || startTime.After(bookingEndTime)) {
			// Room is booked during requested time range
			return false, nil
		}
	}
	// Room is available for the entire requested time range
	return true, nil
}

// RecommendOperatingTime 根据不同手术室的最早时间比较后取最早的作为推荐时间
func RecommendOperatingTime(rooms []model.OperatingRoom, duration time.Duration) (time.Time, error) {
	// 假设我们只考虑接下来的一天
	startTime := time.Now()
	year, month, day := startTime.Date()
	endTime := time.Date(year, month, day, 0, 0, 0, 0, startTime.Location()).Add(time.Duration(config.EndTime) * time.Hour)

	var earliestRecommendation time.Time
	earliestFound := false

	// 遍历每个手术室
	for _, room := range rooms {
		// 查找该手术室在推荐时间范围内的第一个可用时间段
		recommendedTime, err := FindAvailableOperatingTime(room.Id, startTime, endTime, duration)
		if err != nil {
			return time.Time{}, err
		}

		// 如果找到可用时间段
		if !recommendedTime.IsZero() {
			// 如果是第一个找到的时间段或者比当前最早时间更早，则更新最早时间
			if !earliestFound || recommendedTime.Before(earliestRecommendation) {
				earliestRecommendation = recommendedTime
				earliestFound = true
			}
		}
	}

	if earliestFound {
		return earliestRecommendation, nil
	}

	return time.Time{}, errors.New("未找到可用手术室")
}

// FindAvailableOperatingTime 查找可用手术室时间
func FindAvailableOperatingTime(roomID int, startTime time.Time, endTime time.Time, duration time.Duration) (time.Time, error) {
	// 查询手术室在推荐时间范围内的可用时间段
	availableSlots, err := FindAvailableSlots(roomID, startTime, endTime, duration)
	if err != nil {
		return time.Time{}, err
	}

	// 如果找到可用时间段，则返回第一个时间段的开始时间作为推荐时间
	if len(availableSlots) > 0 {
		return availableSlots[0].StartTime, nil
	}

	return time.Time{}, nil
}

// FindAvailableSlots 查询可用时间段
func FindAvailableSlots(roomID int, startTime time.Time, endTime time.Time, duration time.Duration) ([]TimeSlot, error) {
	usedTime, err := GetReserveRoomByTimeRange(roomID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	// 处理预订信息并找出可用时间段
	availableSlots := []TimeSlot{}

	//如果没找到，说明还未使用
	if len(usedTime) == 0 {
		availableSlots = append(availableSlots, TimeSlot{
			StartTime: startTime,
			EndTime:   endTime,
		})
		return availableSlots, nil
	}

	//加上15分钟休息时间
	//第一个区间
	if startTime.Add(duration).Before(usedTime[0].StartTime) {
		availableSlots = append(availableSlots, TimeSlot{
			StartTime: startTime,
			EndTime:   usedTime[0].StartTime,
		})
	}
	//中间区间
	for i := 1; i < len(usedTime); i++ {
		if usedTime[i-1].EndTime.Add(duration).Before(usedTime[i].StartTime) {
			availableSlots = append(availableSlots, TimeSlot{
				StartTime: usedTime[i-1].EndTime,
				EndTime:   usedTime[i].StartTime,
			})
		}
	}
	//末尾
	if usedTime[len(usedTime)-1].EndTime.Add(duration).Before(endTime) {
		availableSlots = append(availableSlots, TimeSlot{
			StartTime: usedTime[len(usedTime)-1].EndTime,
			EndTime:   endTime,
		})
	}

	return availableSlots, nil
}
