/*******
* @Author:qingmeng
* @Description:
* @File:time
* @Date:2024/5/10
 */

package tool

import (
	"fmt"
	"log"
	"time"
)

// 解析时间
func ParseTime(timeStr string) (t time.Time, err error) {
	// 获取当前日期
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02") // 今天的日期，格式为年-月-日
	// 将日期和时间字符串合并
	startTimeStrWithDate := fmt.Sprintf("%s %s", currentDate, timeStr)
	layout := "2006-01-02 15:04" // 年-月-日 时:分 的格式
	cstLocation, _ := time.LoadLocation("Asia/Shanghai")
	t, err = time.ParseInLocation(layout, startTimeStrWithDate, cstLocation)
	if err != nil {
		log.Println("parse time err:", err)
		return
	}
	return t, nil
}
