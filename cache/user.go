/*******
* @Author:qingmeng
* @Description:
* @File:room
* @Date:2024/5/9
 */

package cache

import (
	"OperatingRoomSchedulingSystem/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

// 预约医生
func ReserveUser(uid int, startTime time.Time, endTime time.Time) error {
	key := fmt.Sprintf("user_schedule:%d", uid)
	bookingInfo := map[string]interface{}{
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
	}
	bookingJSON, err := json.Marshal(bookingInfo)
	if err != nil {
		return err
	}
	err = RedisClient.ZAdd(context.Background(), key, &redis.Z{
		Score:  float64(startTime.Unix()), // 使用预订开始时间作为分数
		Member: string(bookingJSON),       // 存储预订信息的 JSON 字符串
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

// 取消手术医生
func DelReserveUser(uid int, startTime time.Time, endTime time.Time) error {
	key := fmt.Sprintf("user_schedule:%d", uid)
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
func GetReserveUserByTimeRange(uid int, startTime time.Time, endTime time.Time) ([]map[string]interface{}, error) {
	key := fmt.Sprintf("user_schedule:%d", uid)
	reservations, err := RedisClient.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min: strconv.FormatInt(startTime.Unix(), 10), // 最小分数，即开始时间
		Max: strconv.FormatInt(endTime.Unix(), 10),   // 最大分数，即结束时间
	}).Result()
	if err != nil {
		return nil, err
	}

	var bookingInfos []map[string]interface{}
	for _, reservation := range reservations {
		var bookingInfo map[string]interface{}
		err := json.Unmarshal([]byte(reservation), &bookingInfo)
		if err != nil {
			return nil, err
		}
		bookingInfos = append(bookingInfos, bookingInfo)
	}
	return bookingInfos, nil
}

// 时间上获取可用房间
func GetAvailableUsersByTime(users []model.User, startTime time.Time, endTime time.Time) ([]model.User, error) {
	var availableUsers []model.User
	for _, user := range users {
		isAvailable, err := IsUserAvailable(user.Id, startTime, endTime)
		if err != nil {
			return nil, err
		}
		if isAvailable {
			availableUsers = append(availableUsers, user)
		}
	}
	return availableUsers, nil
}

// 是否是时间允许的房间
func IsUserAvailable(uid int, startTime time.Time, endTime time.Time) (bool, error) {
	key := fmt.Sprintf("user_schedule:%d", uid)
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
